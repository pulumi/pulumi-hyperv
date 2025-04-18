param([switch]$ForceWindowsMode = $true) # Force Windows mode for testing

# Check PowerShell version
if ($PSVersionTable.PSVersion.Major -lt 7) {
    Write-Error "PowerShell 7 or greater is required to run this script."
    exit 1
}

$IsWindowsEnvironment = if ($PSVersionTable.PSVersion.Major -ge 6) { $IsWindows } else { $ForceWindowsMode }
#region Configuration
# Stop on first error
$ErrorActionPreference = 'Stop'

$PROJECT_NAME = "Pulumi Hyperv Resource Provider"
$PACK = "hyperv"
$PACKDIR = "sdk"
$PROJECT = "github.com/pulumi/pulumi-hyperv"
$NODE_MODULE_NAME = "@pulumi/hyperv"
$NUGET_PKG_NAME = "Pulumi.Hyperv"
$PROVIDER = "pulumi-resource-$PACK"
$PROVIDER_PATH = "provider"
$VERSION_PATH = "$PROVIDER_PATH/pkg/version.Version"
$SCHEMA_FILE = "provider/cmd/pulumi-resource-hyperv/schema.json"
$TESTPARALLELISM = 4

$SHELL = "powershell.exe"
$PATHSEP = ";"
$EXE = ".exe"

# Adjust paths and commands for Windows
$PULUMI = Join-Path -Path (Get-Location) -ChildPath ".pulumi\bin\pulumi$EXE"
$WORKING_DIR = (Get-Location).Path
$COPY = "Copy-Item"
$RM = "Remove-Item -Recurse -Force"
$CP = "Copy-Item -Recurse -Force"
$MKDIR = "New-Item -ItemType Directory -Force"

# Override during CI using `make [TARGET] PROVIDER_VERSION=""` or by setting a PROVIDER_VERSION environment variable
# Local & branch builds will just used this fixed default version unless specified
$PROVIDER_VERSION = "1.0.0-alpha.0+dev"
# Use this normalised version everywhere rather than the raw input to ensure consistency.
if ($IsWindowsEnvironment) {
    $VERSION_GENERIC = $PROVIDER_VERSION
}
else {
    # Attempt to use pulumictl, but fall back to the raw version if it fails
    $VERSION_GENERIC = try {
        & pulumictl convert-version --language generic --version "$PROVIDER_VERSION" 2>$null
    }
    catch {
        $PROVIDER_VERSION
    }
}

# Need to pick up locally pinned pulumi-langage-* plugins.
$env:PULUMI_IGNORE_AMBIENT_PLUGINS = $true

#endregion

#region Helper Functions

function Invoke-CommandWithChangeDirectory {
    param (
        [string]$Path,
        [scriptblock]$ScriptBlock
    )
    Write-Host "Entering directory: $Path"
    
    # Create directory if it doesn't exist
    if (-not (Test-Path $Path)) {
        Write-Host "Directory does not exist, creating: $Path"
        if ($IsWindowsEnvironment) {
            New-Item -ItemType Directory -Force -Path $Path | Out-Null
        }
        else {
            mkdir -p $Path
        }
    }
    
    Push-Location $Path
    try {
        & $ScriptBlock
    }
    finally {
        Pop-Location
        Write-Host "Leaving directory: $Path"
    }
}

function Execute-Command {
    param (
        [string]$Command
    )
    Write-Host "Executing: $Command"
    Invoke-Expression $Command
}

#endregion

#region Targets

function Target-ensure {
    Target-tidy
}

function Target-tidy {
    Target-tidy_provider
    Target-tidy_examples
    if ($IsWindowsEnvironment) {
        Invoke-CommandWithChangeDirectory "sdk" { go mod tidy }
        Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
    }
    else {
        Invoke-CommandWithChangeDirectory "sdk" { go mod tidy }
        Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
    }
}

function Target-tidy_examples {
    Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
}

function Target-tidy_provider {
    if ($IsWindowsEnvironment) {
        Invoke-CommandWithChangeDirectory "provider" { go mod tidy }
    }
    else {
        Invoke-CommandWithChangeDirectory "provider" { go mod tidy }
    }
}

function Target-SchemaFile {
    Write-Host "Generating $($SCHEMA_FILE)"
    # Ensure bin directory exists before continuing
    if (-not (Test-Path "$WORKING_DIR\bin")) {
        New-Item -ItemType Directory -Path "$WORKING_DIR\bin" -Force | Out-Null
    }
    
    # Check if provider exists, build it if not
    $fullPath = "$WORKING_DIR\bin\$PROVIDER$EXE"
    if (-not (Test-Path $fullPath)) {
        Target-provider
    }
    
    # Check if the provider binary exists now
    if (Test-Path $fullPath) {
        # Check if Pulumi exists and invoke it with arguments separately
        if (Test-Path $PULUMI) {
            $schemaOutput = & "$PULUMI" "package" "get-schema" "$fullPath"
            $schema = $schemaOutput | ConvertFrom-Json
            $schema.PSObject.Properties.Remove('version')
            # Convert to JSON and explicitly ensure UTF-8 without BOM and with LF line endings
            $jsonContent = $schema | ConvertTo-Json -Depth 100
            # Replace all Windows line endings with Unix line endings
            $jsonContent = $jsonContent.Replace("`r`n", "`n")
            # Remove any trailing newline to prevent git warnings about EOL at EOF
            $jsonContent = $jsonContent.TrimEnd("`n")
            # Construct full path from working directory
            $fullSchemaPath = Join-Path -Path $WORKING_DIR -ChildPath $SCHEMA_FILE
            Write-Host "Writing schema to: $fullSchemaPath"
            # Use .NET methods to write the file with explicit UTF-8 without BOM encoding
            [System.IO.File]::WriteAllText($fullSchemaPath, $jsonContent, [System.Text.UTF8Encoding]::new($false))
            Write-Host "Schema file generated at $fullSchemaPath"
            # Count and display the number of lines in the schema file
            $lineCount = (Get-Content $fullSchemaPath | Measure-Object -Line).Lines
            Write-Host "Schema file contains $lineCount lines"
        }
        else {
            Write-Error "Pulumi executable not found at $PULUMI"
            exit 1
        }
    }
    else {
        Write-Error "Provider binary not found at $fullPath after attempting to build it"
        exit 1
    }
}

function Target-codegen {
    Target-Pulumi
    Target-provider
    Target-SchemaFile
    Target-sdk_dotnet
    Target-sdk_go
    Target-nodejs_sdk 
    Target-sdk_python
    Target-sdk_java
}

function Target-sdk {
    param (
        [string]$language
    )
    
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Target-Pulumi
    }
    
    if ($IsWindowsEnvironment) {
        $sdkPath = "sdk\$language"
        if (Test-Path $sdkPath) {
            Remove-Item -Recurse -Force $sdkPath
        }
        # Call Pulumi executable with separate arguments
        & "$PULUMI" "package" "gen-sdk" "--language" "$language" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    }
    else {
        $sdkPath = "sdk/$language"
        if (Test-Path $sdkPath) {
            Remove-Item -Recurse -Force $sdkPath
        }
        & "$PULUMI" package gen-sdk --language "$language" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    }
}

function Target-sdk_java {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Target-Pulumi
    }
    
    $sdkPath = "sdk\java"
    if (Test-Path $sdkPath) {
        Remove-Item -Recurse -Force $sdkPath
    }
    # Call Pulumi executable with separate arguments
    & "$PULUMI" "package" "gen-sdk" "--language" "java" "$SCHEMA_FILE"
}

function Target-sdk_python {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Target-Pulumi
    }
    
    $sdkPath = "sdk\python"
    if (Test-Path $sdkPath) {
        Remove-Item -Recurse -Force $sdkPath
    }
    # Call Pulumi executable with separate arguments
    & "$PULUMI" "package" "gen-sdk" "--language" "python" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    Copy-Item README.md "sdk\python\" -Force
}

function Target-sdk_dotnet {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Target-Pulumi
    }
    
    # Ensure bin directory exists before continuing
    if (-not (Test-Path "$WORKING_DIR\bin")) {
        New-Item -ItemType Directory -Path "$WORKING_DIR\bin" -Force
    }
    
    $sdkPath = "sdk\dotnet"
    if (Test-Path $sdkPath) {
        Remove-Item -Recurse -Force $sdkPath
    }
    # Call Pulumi executable with separate arguments
    & "$PULUMI" "package" "gen-sdk" "--language" "dotnet" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    # Copy the logo to the dotnet directory before building so it can be included in the nuget package archive.
    # https://github.com/pulumi/pulumi-hyperv-provider/issues/243
    Invoke-CommandWithChangeDirectory "sdk\dotnet" {
        Copy-Item "$WORKING_DIR\assets\logo.png" "logo.png" -Force
    }
}

function Target-sdk_go {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Target-Pulumi
    }
    
    $sdkPath = "sdk\go"
    if (Test-Path $sdkPath) {
        # Instead of deleting entire directory, preserve doc.go file first
        if (Test-Path "$sdkPath\hyperv\doc.go") {
            $docGoContent = Get-Content "$sdkPath\hyperv\doc.go" -Raw
            Remove-Item -Recurse -Force $sdkPath
            
            # Recreate hyperv directory structure
            New-Item -ItemType Directory -Force -Path "$sdkPath\hyperv" | Out-Null
            
            # Restore doc.go file
            Set-Content -Path "$sdkPath\hyperv\doc.go" -Value $docGoContent
        }
        else {
            Remove-Item -Recurse -Force $sdkPath
        }
    }
    
    # Call Pulumi executable with separate arguments
    & "$PULUMI" "package" "gen-sdk" "--language" "go" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
}

function Target-provider {
    if ($IsWindowsEnvironment) {
        Write-Host "Building provider for Windows: $WORKING_DIR\bin\$PROVIDER$EXE with version $VERSION_GENERIC"
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o "$WORKING_DIR\bin\$PROVIDER$EXE" -ldflags "-X $PROJECT/$VERSION_PATH=$VERSION_GENERIC" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    }
    else {
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o "$WORKING_DIR/bin/$PROVIDER$EXE" -ldflags "-X $PROJECT/$VERSION_PATH=$VERSION_GENERIC" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    }
}

function Target-provider_debug {
    if ($IsWindowsEnvironment) {
        $version = $VERSION_GENERIC
        $outPath = "$WORKING_DIR\bin\$PROVIDER$EXE"
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o $outPath -gcflags "all=-N -l" -ldflags "-X '$PROJECT/$VERSION_PATH=$version'" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    }
    else {
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o "$WORKING_DIR/bin/$PROVIDER$EXE" -gcflags="all=-N -l" -ldflags "-X $PROJECT/$VERSION_PATH=$VERSION_GENERIC" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    }
}

function Target-test_provider {
    Target-tidy_provider
    Invoke-CommandWithChangeDirectory "provider" {
        go test -short -v -count=1 -cover -timeout 2h -parallel $TESTPARALLELISM -coverprofile="coverage.txt" ./...
    }
}

function Target-dotnet_sdk {
    # Use Target-SchemaFile only if schema.json doesn't exist
    if (-not (Test-Path $SCHEMA_FILE)) {
        Target-SchemaFile
    }
    else {
        # Ensure Pulumi exists
        if (-not (Test-Path $PULUMI)) {
            Target-Pulumi
        }
    }
    
    Target-sdk_dotnet
    Invoke-CommandWithChangeDirectory "sdk/dotnet" {
        "$VERSION_GENERIC" | Out-File version.txt
        dotnet build
    }
}

function Target-go_sdk {
    # Use Target-SchemaFile only if schema.json doesn't exist
    if (-not (Test-Path $SCHEMA_FILE)) {
        Target-SchemaFile
    }
    else {
        # Ensure Pulumi exists
        if (-not (Test-Path $PULUMI)) {
            Target-Pulumi
        }
    }
    
    Target-sdk_go
}

function Target-nodejs_sdk {
    # Use Target-SchemaFile only if schema.json doesn't exist
    if (-not (Test-Path $SCHEMA_FILE)) {
        Target-SchemaFile
    }
    else {
        # Ensure Pulumi exists
        if (-not (Test-Path $PULUMI)) {
            Target-Pulumi
        }
    }
    
    Target-sdk "nodejs"
    Invoke-CommandWithChangeDirectory "sdk/nodejs" {
        yarn install
        yarn run tsc --version
        yarn run tsc
    }
    Write-Host "Copying nodejs SDK files to bin directory"
    Copy-Item README.md -Destination "sdk/nodejs/" -Force
    Copy-Item "sdk/nodejs/package.json" -Destination "sdk/nodejs/bin/" -Force
    Copy-Item "sdk/nodejs/yarn.lock" -Destination "sdk/nodejs/bin/" -Force
}

function Target-python_sdk {
    # Use Target-SchemaFile only if schema.json doesn't exist
    if (-not (Test-Path $SCHEMA_FILE)) {
        Target-SchemaFile
    }
    else {
        # Ensure Pulumi exists
        if (-not (Test-Path $PULUMI)) {
            Target-Pulumi
        }
    }
    
    Copy-Item README.md "sdk/python/" -Force
    Invoke-CommandWithChangeDirectory "sdk/python" {
        # Check if directories exist before removing them
        if (Test-Path ./bin/) { Remove-Item ./bin/ -Recurse -Force }
        if (Test-Path ../python.bin/) { Remove-Item ../python.bin/ -Recurse -Force }
        Copy-Item . ../python.bin -Recurse -Force
        Move-Item ../python.bin ./bin -Force
        
        # Use Windows-compatible Python paths
        if ($IsWindowsEnvironment) {
            python -m venv venv
            & ".\venv\Scripts\python" -m pip install build
            Invoke-CommandWithChangeDirectory "./bin" {
                & "..\venv\Scripts\python" -m build .
            }
        } 
        else {
            python3 -m venv venv
            ./venv/bin/python -m pip install build
            Invoke-CommandWithChangeDirectory "./bin" {
                ../venv/bin/python -m build .
            }
        }
    }
}

function Find-JavaExe {
    # First check if java is in the PATH
    $javaPath = Get-Command -Name "java.exe" -ErrorAction SilentlyContinue
    if ($javaPath) {
        return $javaPath.Source
    }
    
    # Check common Java installation locations
    $commonPaths = @(
        "${env:ProgramFiles}\Java\*\bin\java.exe",
        "${env:ProgramFiles(x86)}\Java\*\bin\java.exe",
        "${env:JAVA_HOME}\bin\java.exe",
        "$env:LOCALAPPDATA\Programs\Eclipse Adoptium\*\bin\java.exe",
        "$env:LOCALAPPDATA\Programs\Eclipse Foundation\*\bin\java.exe",
        "$env:LOCALAPPDATA\Programs\Temurin\*\bin\java.exe",
        "$env:LOCALAPPDATA\Programs\Microsoft\jdk-*\bin\java.exe",
        "${env:ProgramFiles}\Eclipse Adoptium\*\bin\java.exe",
        "${env:ProgramFiles}\Eclipse Adoptium\jdk-*\bin\java.exe",
        "${env:ProgramFiles}\Eclipse Foundation\*\bin\java.exe",
        "${env:ProgramFiles}\Temurin\*\bin\java.exe",
        "${env:ProgramFiles}\OpenJDK\*\bin\java.exe",
        "${env:ProgramFiles(x86)}\Eclipse Adoptium\*\bin\java.exe",
        "${env:ProgramFiles(x86)}\Eclipse Foundation\*\bin\java.exe",
        "${env:ProgramFiles(x86)}\Temurin\*\bin\java.exe",
        "${env:ProgramFiles(x86)}\OpenJDK\*\bin\java.exe"
    )
    
    foreach ($path in $commonPaths) {
        $found = Get-ChildItem -Path $path -ErrorAction SilentlyContinue | Sort-Object -Property LastWriteTime -Descending | Select-Object -First 1
        if ($found) {
            Write-Host "Found Java at: $($found.FullName)"
            return $found.FullName
        }
    }
    
    # If Java not found, return null
    Write-Warning "Java executable (java.exe) not found. Please install Java or ensure it's in your PATH."
    return $null
}

function Find-GradleExe {
    # First check if gradle is in the PATH
    $gradlePath = Get-Command -Name "gradle.bat" -ErrorAction SilentlyContinue
    if ($gradlePath) {
        return $gradlePath.Source
    }
    
    # Check common Gradle installation locations
    $commonPaths = @(
        "${env:ProgramFiles}\Gradle\*\bin\gradle.bat",
        "${env:ProgramFiles(x86)}\Gradle\*\bin\gradle.bat",
        "${env:GRADLE_HOME}\bin\gradle.bat",
        "$env:USERPROFILE\.gradle\wrapper\dists\*\*\*\bin\gradle.bat"
    )
    
    foreach ($path in $commonPaths) {
        $found = Get-ChildItem -Path $path -ErrorAction SilentlyContinue | Sort-Object -Property LastWriteTime -Descending | Select-Object -First 1
        if ($found) {
            Write-Host "Found Gradle at: $($found.FullName)"
            return $found.FullName
        }
    }
    
    # If Gradle not found, return null
    Write-Warning "Gradle executable (gradle.bat) not found. Will try to use gradlew wrapper if available."
    return $null
}

function Target-java_sdk {
    # Use Target-SchemaFile only if schema.json doesn't exist
    if (-not (Test-Path $SCHEMA_FILE)) {
        Target-SchemaFile
    }
    else {
        # Ensure Pulumi exists
        if (-not (Test-Path $PULUMI)) {
            Target-Pulumi
        }
    }
    
    $env:PACKAGE_VERSION = $VERSION_GENERIC
    Target-sdk_java
    
    if ($IsWindowsEnvironment) {
        # Find java.exe
        $javaExe = Find-JavaExe
        if ($javaExe) {
            $env:PATH = "$(Split-Path -Parent $javaExe);$env:PATH"
            Write-Host "Using Java from: $javaExe"
        }
        
        # Find gradle.bat
        $gradleExe = Find-GradleExe
        if ($gradleExe) {
            $env:PATH = "$(Split-Path -Parent $gradleExe);$env:PATH"
            Write-Host "Using Gradle from: $gradleExe"
        }
        
        # Check if gradle wrapper exists, use it if available
        $gradleWrapperWin = ".\gradlew.bat"
        $gradleWrapperUnix = "./gradlew"
        
        Invoke-CommandWithChangeDirectory "sdk/java" {
            if (Test-Path $gradleWrapperWin) {
                Write-Host "Using Gradle wrapper: $gradleWrapperWin"
                & $gradleWrapperWin --console=plain build
            }
            elseif (Test-Path $gradleWrapperUnix) {
                Write-Host "Using Gradle wrapper: $gradleWrapperUnix"
                & $gradleWrapperUnix --console=plain build
            }
            elseif ($gradleExe) {
                Write-Host "Using Gradle from PATH"
                & $gradleExe --console=plain build
            }
            else {
                Write-Host "Attempting to use 'gradle' command from PATH"
                gradle --console=plain build
            }
        }
    }
    else {
        # Non-Windows systems
        Invoke-CommandWithChangeDirectory "sdk/java" {
            if (Test-Path "./gradlew") {
                Write-Host "Using Gradle wrapper: ./gradlew"
                chmod +x ./gradlew
                ./gradlew --console=plain build
            }
            else {
                Write-Host "Using gradle from PATH"
                gradle --console=plain build
            }
        }
    }
}

function Target-build {
    Target-provider
    Target-build_sdks
}

function Target-build_sdks {
    Target-dotnet_sdk
    Target-go_sdk
    Target-nodejs_sdk
    Target-python_sdk
    Target-java_sdk
}

function Target-only_build {
    Target-build
}

function Target-lint {
    Invoke-CommandWithChangeDirectory "provider" {
        golangci-lint --path-prefix provider --config ../.golangci.yml run
    }
}

function Target-install {
    Target-install_nodejs_sdk
    Target-install_dotnet_sdk
    if (Test-Path "$WORKING_DIR/bin/$PROVIDER$EXE") {
        New-Item -ItemType Directory -Force -Path "$env:GOPATH/bin" | Out-Null
        Copy-Item "$WORKING_DIR/bin/$PROVIDER$EXE" "$env:GOPATH/bin/" -Force
    }
}

$GO_TEST = "go test -v -count=1 -cover -timeout 2h -parallel $TESTPARALLELISM"

function Target-test_all {
    Target-test
    Invoke-CommandWithChangeDirectory "provider/pkg" { & "$GO_TEST ./..." }
    Invoke-CommandWithChangeDirectory "tests/sdk/nodejs" { & "$GO_TEST ./..." }
    Invoke-CommandWithChangeDirectory "tests/sdk/python" { & "$GO_TEST ./..." }
    Invoke-CommandWithChangeDirectory "tests/sdk/dotnet" { & "$GO_TEST ./..." }
    Invoke-CommandWithChangeDirectory "tests/sdk/go" { & "$GO_TEST ./..." }
}

function Target-install_dotnet_sdk {
    # Create nuget directory if it doesn't exist
    New-Item -ItemType Directory -Force "$WORKING_DIR/nuget"
    
    # Remove any existing packages (with error action silently continue)
    Remove-Item "$WORKING_DIR/nuget/$NUGET_PKG_NAME.*.nupkg" -ErrorAction SilentlyContinue
    
    # Find SDK nupkg files and copy them to nuget directory
    $nupkgFiles = Get-ChildItem -Path "$WORKING_DIR/sdk/dotnet/bin" -Recurse -Filter "*.nupkg" -ErrorAction SilentlyContinue
    if ($nupkgFiles) {
        $nupkgFiles | ForEach-Object { Copy-Item -Path $_.FullName -Destination "$WORKING_DIR/nuget" -Force }
    }
    else {
        Write-Warning "No .nupkg files found in $WORKING_DIR/sdk/dotnet/bin"
    }
}

function Target-install_python_sdk {
    # target intentionally blank
}

function Target-install_go_sdk {
    # target intentionally blank
}

function Target-install_java_sdk {
    # target intentionally blank
}

function Target-install_nodejs_sdk {
    try {
        yarn unlink --cwd "$WORKING_DIR/sdk/nodejs/bin"
    }
    catch {
        Write-Warning "Failed to unlink nodejs sdk: $($_.Exception.Message)"
    }
    yarn link --cwd "$WORKING_DIR/sdk/nodejs/bin"
}

function Target-test {
    Target-tidy_examples
    Target-test_provider
    Invoke-CommandWithChangeDirectory "examples" { go test -v -tags=all -timeout 2h }
}

function Target-Pulumi {
    Write-Host "Ensuring Pulumi CLI is installed and up-to-date"
    Set-Location -Path $WORKING_DIR
    
    # Define the Pulumi location
    if ($IsWindowsEnvironment) {
        $PULUMI_INSTALL_DIR = Join-Path -Path $WORKING_DIR -ChildPath ".pulumi"
        $PULUMI_DIR = Join-Path -Path $WORKING_DIR -ChildPath ".pulumi\bin"
        $PULUMI_EXE = Join-Path -Path $PULUMI_DIR -ChildPath "pulumi$EXE"
        $PULUMI_VERSION = Get-Content .pulumi.version

        if (Test-Path $PULUMI_EXE) {
            $CURRENT_VERSION = $(& "$PULUMI_EXE" version).TrimStart('v')
            if ($CURRENT_VERSION -ne $PULUMI_VERSION) {
                Write-Host "Upgrading $PULUMI_EXE from $CURRENT_VERSION to $PULUMI_VERSION"
                Remove-Item $PULUMI_EXE -Force
            }
        }
        
        if (!(Test-Path $PULUMI_EXE)) {
            Write-Host "Installing Pulumi CLI version $PULUMI_VERSION"
            if (!(Test-Path $PULUMI_DIR)) {
                New-Item -ItemType Directory -Path $PULUMI_DIR -Force | Out-Null
            }
            Invoke-WebRequest -Uri https://get.pulumi.com/install.ps1 -OutFile pulumi-install.ps1
            & "./pulumi-install.ps1" -Version $($PULUMI_VERSION.TrimStart('v')) -InstallRoot $PULUMI_INSTALL_DIR -NoEditPath
            Remove-Item pulumi-install.ps1
        }
        
        # Update the global variable to use the correct path
        $script:PULUMI = $PULUMI_EXE
    }
    else {
        # Similar logic for non-Windows
        $PULUMI_VERSION = Get-Content .pulumi.version
        if (Test-Path $PULUMI) {
            $CURRENT_VERSION = & $PULUMI version
            if ($CURRENT_VERSION -ne $PULUMI_VERSION) {
                Write-Host "Upgrading $PULUMI from $CURRENT_VERSION to $PULUMI_VERSION"
                Remove-Item $PULUMI -Force
            }
        }
        if (!(Test-Path $PULUMI)) {
            Write-Host "Installing Pulumi CLI version $PULUMI_VERSION"
            curl -fsSL https://get.pulumi.com | sh -s -- --version $($PULUMI_VERSION.TrimStart('v'))
        }
    }
}

# Signing targets (conditional)
function Target-bin_jsign_6_0_jar {
    Write-Host "Downloading jsign-6.0.jar"
    Invoke-WebRequest -Uri https://github.com/ebourg/jsign/releases/download/6.0/jsign-6.0.jar -OutFile bin/jsign-6.0.jar
}

function Target-sign_goreleaser_exe {
    param (
        [string]$GORELEASER_ARCH
    )
    Write-Host "Signing goreleaser exe for architecture: $GORELEASER_ARCH"
    Target-bin_jsign_6_0_jar

    # Only sign windows binary if fully configured.
    # Test variables set by joining with | between and looking for || showing at least one variable is empty.
    # Move the binary to a temporary location and sign it there to avoid the target being up-to-date if signing fails.
    if ($env:SKIP_SIGNING -ne "true") {
        if ("|$($env:AZURE_SIGNING_CLIENT_ID)|$($env:AZURE_SIGNING_CLIENT_SECRET)|$($env:AZURE_SIGNING_TENANT_ID)|$($env:AZURE_SIGNING_KEY_VAULT_URI)|" -like "*||*") {
            Write-Host "Can't sign windows binaries as required configuration not set: AZURE_SIGNING_CLIENT_ID, AZURE_SIGNING_CLIENT_SECRET, AZURE_SIGNING_TENANT_ID, AZURE_SIGNING_KEY_VAULT_URI"
            Write-Host "To rebuild with signing delete the unsigned windows exe file and rebuild with the fixed configuration"
            if ($env:CI -eq "true") {
                exit 1
            }
        }
        else {
            $file = "dist/build-provider-sign-windows_windows_$GORELEASER_ARCH/pulumi-resource-hyperv.exe"
            Move-Item $file "$file.unsigned" -Force
            az login --service-principal --username $env:AZURE_SIGNING_CLIENT_ID --password $env:AZURE_SIGNING_CLIENT_SECRET --tenant $env:AZURE_SIGNING_TENANT_ID --output none
            $ACCESS_TOKEN = az account get-access-token --resource "https://vault.azure.net" | ConvertFrom-Json | Select-Object -ExpandProperty accessToken
            java -jar bin/jsign-6.0.jar --storetype AZUREKEYVAULT --keystore "PulumiCodeSigning" --url $env:AZURE_SIGNING_KEY_VAULT_URI --storepass $ACCESS_TOKEN "$file.unsigned"
            Move-Item "$file.unsigned" $file -Force
            az logout
        }
    }
}

function Target-sign_goreleaser_exe_amd64 {
    Target-sign_goreleaser_exe -GORELEASER_ARCH amd64_v1
}

function Target-sign_goreleaser_exe_arm64 {
    Target-sign_goreleaser_exe -GORELEASER_ARCH arm64
}

#endregion

#region Argument Parsing and Execution

# Default target
$target = "build"

# Override target if an argument is provided
if ($args.Length -gt 0) {
    $target = $args[0]
}

switch ($target) {
    "ensure" { Target-ensure }
    "tidy" { Target-tidy }
    "tidy_examples" { Target-tidy_examples }
    "tidy_provider" { Target-tidy_provider }
    "codegen" { Target-codegen }
    "sdk/java" { Target-sdk_java }
    "sdk/python" { Target-sdk_python }
    "sdk/dotnet" { Target-sdk_dotnet }
    "sdk/go" { Target-sdk_go }
    "sdk/nodejs" { Target-nodejs_sdk }
    "provider" { Target-provider }
    "provider_debug" { Target-provider_debug }
    "test_provider" { Target-test_provider }
    "dotnet_sdk" { Target-dotnet_sdk }
    "go_sdk" { Target-go_sdk }
    "nodejs_sdk" { Target-nodejs_sdk }
    "python_sdk" { Target-python_sdk }
    "bin/pulumi-java-gen" { Target-bin_pulumi_java_gen }
    "java_sdk" { Target-java_sdk }
    "build" { Target-build }
    "build_sdks" { Target-build_sdks }
    "only_build" { Target-only_build }
    "lint" { Target-lint }
    "install" { Target-install }
    "test_all" { Target-test_all }
    "install_dotnet_sdk" { Target-install_dotnet_sdk }
    "install_python_sdk" { Target-install_python_sdk }
    "install_go_sdk" { Target-install_go_sdk }
    "install_java_sdk" { Target-install_java_sdk }
    "install_nodejs_sdk" { Target-install_nodejs_sdk }
    "test" { Target-test }
    "Pulumi" { Target-Pulumi }
    "bin/jsign-6.0.jar" { Target-bin_jsign_6_0_jar }
    "sign-goreleaser-exe-amd64" { Target-sign_goreleaser_exe_amd64 }
    "sign-goreleaser-exe-arm64" { Target-sign_goreleaser_exe_arm64 }
    default {
        Write-Error "Unknown target: $target"
        exit 1
    }
}

#endregion
