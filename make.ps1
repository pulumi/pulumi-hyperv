param([switch]$IsWindows=$true) 

#region Configuration
# Stop on first error
$ErrorActionPreference = 'Stop'

$PROJECT_NAME = "Pulumi Hyperv Resource Provider"
$PACK = "hyperv"
$PACKDIR = "sdk"
$PROJECT = "github.com/pulumi/pulumi-hyperv-provider"
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
if ($IsWindows) {
    $VERSION_GENERIC = $PROVIDER_VERSION
} else {
    # Attempt to use pulumictl, but fall back to the raw version if it fails
    $VERSION_GENERIC = try {
        & pulumictl convert-version --language generic --version "$PROVIDER_VERSION" 2>$null
    } catch {
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
        if ($IsWindows) {
            New-Item -ItemType Directory -Force -Path $Path | Out-Null
        } else {
            mkdir -p $Path
        }
    }
    
    Push-Location $Path
    try {
        & $ScriptBlock
    } finally {
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
    if ($IsWindows) {
        Invoke-CommandWithChangeDirectory "sdk" { go mod tidy }
        Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
    } else {
        Invoke-CommandWithChangeDirectory "sdk" { go mod tidy }
        Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
    }
}

function Target-tidy_examples {
    Invoke-CommandWithChangeDirectory "examples" { go mod tidy }
}

function Target-tidy_provider {
    if ($IsWindows) {
        Invoke-CommandWithChangeDirectory "provider" { go mod tidy }
    } else {
        Invoke-CommandWithChangeDirectory "provider" { go mod tidy }
    }
}

function Target-SchemaFile {
    Write-Host "Generating $($SCHEMA_FILE)"
    Target-provider
    if ($IsWindows) {
        $fullPath = "$WORKING_DIR\bin\$PROVIDER$EXE"
        if (Test-Path $fullPath) {
            # Check if Pulumi exists and invoke it with arguments separately
            if (Test-Path $PULUMI) {
                $schemaOutput = & "$PULUMI" "package" "get-schema" "$fullPath"
                $schema = $schemaOutput | ConvertFrom-Json
                $schema.PSObject.Properties.Remove('version')
                $schema | ConvertTo-Json -Depth 10 | Set-Content -Encoding default -Path $SCHEMA_FILE
            } else {
                Write-Error "Pulumi executable not found at $PULUMI"
                exit 1
            }
        } else {
            Write-Error "Provider binary not found at $fullPath"
            exit 1
        }
    } else {
        if (Test-Path $PULUMI) {
            & "$PULUMI" package get-schema "$WORKING_DIR/bin/$PROVIDER$EXE" | jq 'del(.version)' > $SCHEMA_FILE
        } else {
            Write-Error "Pulumi executable not found at $PULUMI"
            exit 1
        }
    }
}

function Target-codegen {
    # Target-Pulumi
    Target-SchemaFile
    Target-sdk_dotnet
	Target-sdk_go
	Target-sdk_nodejs
	Target-sdk_python
	Target-sdk_java
}

function Target-sdk {
    param (
        [string]$language
    )
    
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Write-Error "Pulumi executable not found at $PULUMI"
        exit 1
    }
    
    if ($IsWindows) {
        $sdkPath = "sdk\$language"
        if (Test-Path $sdkPath) {
            Remove-Item -Recurse -Force $sdkPath
        }
        # Call Pulumi executable with separate arguments
        & "$PULUMI" "package" "gen-sdk" "--language" "$language" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    } else {
        $sdkPath = "sdk/$language"
        if (Test-Path $sdkPath) {
            Remove-Item -Recurse -Force $sdkPath
        }
        & "$PULUMI" package gen-sdk --language "$language" "$SCHEMA_FILE" --version "$VERSION_GENERIC"
    }
}

function Target-sdk_java {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Write-Error "Pulumi executable not found at $PULUMI"
        exit 1
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
        Write-Error "Pulumi executable not found at $PULUMI"
        exit 1
    }
    
    $sdkPath = "sdk\python"
    if (Test-Path $sdkPath) {
        Remove-Item -Recurse -Force $sdkPath
    }
    # Call Pulumi executable with separate arguments
    & "$PULUMI" "package" "gen-sdk" "--language" "python" "$SCHEMA_FILE" "--version" "$VERSION_GENERIC"
    Copy-Item README.md "sdk\python\"
}

function Target-sdk_dotnet {
    # Ensure Pulumi exists
    if (-not (Test-Path $PULUMI)) {
        Write-Error "Pulumi executable not found at $PULUMI"
        exit 1
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
        Copy-Item "$WORKING_DIR\assets\logo.png" "logo.png"
    }
}

function Target-sdk_go {
    Target-sdk "go"
}

function Target-sdk_nodejs {
    Target-sdk "nodejs"
}

function Target-provider {
    if ($IsWindows) {
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o "$WORKING_DIR\bin\$PROVIDER$EXE" -ldflags "-X '$PROJECT/$VERSION_PATH=$VERSION_GENERIC'" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    } else {
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o "$WORKING_DIR/bin/$PROVIDER$EXE" -ldflags "-X $PROJECT/$VERSION_PATH=$VERSION_GENERIC" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    }
}

function Target-provider_debug {
    if ($IsWindows) {
        $version = $VERSION_GENERIC
        $outPath = "$WORKING_DIR\bin\$PROVIDER$EXE"
        Invoke-CommandWithChangeDirectory "provider" {
            go build -o $outPath -gcflags "all=-N -l" -ldflags "-X '$PROJECT/$VERSION_PATH=$version'" "$PROJECT/$PROVIDER_PATH/cmd/$PROVIDER"
        }
    } else {
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
    Target-sdk_dotnet
    Invoke-CommandWithChangeDirectory "sdk/dotnet" {
        "$VERSION_GENERIC" | Out-File version.txt
        dotnet build
    }
}

function Target-go_sdk {
    Target-sdk_go
}

function Target-nodejs_sdk {
    Target-sdk_nodejs
    Invoke-CommandWithChangeDirectory "sdk/nodejs" {
        yarn install
        yarn run tsc
    }
    Copy-Item README.md "sdk/nodejs/package.json" "sdk/nodejs/yarn.lock" "sdk/nodejs/bin/" -Destination "sdk/nodejs/"
}

function Target-python_sdk {
    Copy-Item README.md "sdk/python/"
    Invoke-CommandWithChangeDirectory "sdk/python" {
        Remove-Item ./bin/, ../python.bin/ -Recurse -Force
        Copy-Item . ../python.bin -Recurse -Force
        Move-Item ../python.bin ./bin
        python3 -m venv venv
        ./venv/bin/python -m pip install build
        Invoke-CommandWithChangeDirectory "./bin" {
            ../venv/bin/python -m build .
        }
    }
}

function Target-bin_pulumi_java_gen {
    Write-Host "pulumi-java-gen is no longer necessary"
}

function Target-java_sdk {
    $env:PACKAGE_VERSION = $VERSION_GENERIC
    Target-sdk_java
    Invoke-CommandWithChangeDirectory "sdk/java" {
        gradle --console=plain build
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
    Copy-Item "$WORKING_DIR/bin/$PROVIDER$EXE" "$env:GOPATH/bin"
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
    Remove-Item "$WORKING_DIR/nuget/$NUGET_PKG_NAME.*.nupkg" -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Force "$WORKING_DIR/nuget"
    Get-ChildItem . -Filter "*.nupkg" | ForEach-Object { Copy-Item -Path $_.FullName -Destination "$WORKING_DIR/nuget" }
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
    } catch {
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
    $HOME = $WORKING_DIR
    
    # Define the Pulumi location
    if ($IsWindows) {
        $PULUMI_DIR = Join-Path -Path $WORKING_DIR -ChildPath ".pulumi\bin"
        $PULUMI_EXE = Join-Path -Path $PULUMI_DIR -ChildPath "pulumi$EXE"
        $PULUMI_VERSION = Get-Content .pulumi.version
        
        if (Test-Path $PULUMI_EXE) {
            $CURRENT_VERSION = & "$PULUMI_EXE" version
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
            Invoke-WebRequest -Uri https://get.pulumi.com -OutFile pulumi-install.ps1
            & ./pulumi-install.ps1 --version $($PULUMI_VERSION.TrimStart('v'))
        }
        
        # Update the global variable to use the correct path
        $script:PULUMI = $PULUMI_EXE
    } else {
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
        } else {
            $file = "dist/build-provider-sign-windows_windows_$GORELEASER_ARCH/pulumi-resource-hyperv.exe"
            Move-Item $file "$file.unsigned"
            az login --service-principal --username $env:AZURE_SIGNING_CLIENT_ID --password $env:AZURE_SIGNING_CLIENT_SECRET --tenant $env:AZURE_SIGNING_TENANT_ID --output none
            $ACCESS_TOKEN = az account get-access-token --resource "https://vault.azure.net" | ConvertFrom-Json | Select-Object -ExpandProperty accessToken
            java -jar bin/jsign-6.0.jar --storetype AZUREKEYVAULT --keystore "PulumiCodeSigning" --url $env:AZURE_SIGNING_KEY_VAULT_URI --storepass $ACCESS_TOKEN "$file.unsigned"
            Move-Item "$file.unsigned" $file
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
    "ensure"              { Target-ensure }
    "tidy"                { Target-tidy }
    "tidy_examples"       { Target-tidy_examples }
    "tidy_provider"       { Target-tidy_provider }
    "codegen"             { Target-codegen }
    "sdk/java"            { Target-sdk_java }
    "sdk/python"          { Target-sdk_python }
    "sdk/dotnet"          { Target-sdk_dotnet }
	"sdk/go"              { Target-sdk_go }
	"sdk/nodejs"          { Target-sdk_nodejs }
    "provider"            { Target-provider }
    "provider_debug"      { Target-provider_debug }
    "test_provider"       { Target-test_provider }
    "dotnet_sdk"          { Target-dotnet_sdk }
	"go_sdk"              { Target-go_sdk }
	"nodejs_sdk"          { Target-nodejs_sdk }
	"python_sdk"          { Target-python_sdk }
    "bin/pulumi-java-gen" { Target-bin_pulumi_java_gen }
    "java_sdk"            { Target-java_sdk }
    "build"               { Target-build }
    "build_sdks"          { Target-build_sdks }
    "only_build"          { Target-only_build }
    "lint"                { Target-lint }
    "install"             { Target-install }
    "test_all"            { Target-test_all }
    "install_dotnet_sdk"  { Target-install_dotnet_sdk }
    "install_python_sdk"  { Target-install_python_sdk }
    "install_go_sdk"      { Target-install_go_sdk }
    "install_java_sdk"      { Target-install_java_sdk }
    "install_nodejs_sdk"  { Target-install_nodejs_sdk }
    "test"                { Target-test }
    "Pulumi"              { Target-Pulumi }
    "bin/jsign-6.0.jar"    { Target-bin_jsign_6_0_jar }
    "sign-goreleaser-exe-amd64" { Target-sign_goreleaser_exe_amd64 }
    "sign-goreleaser-exe-arm64" { Target-sign_goreleaser_exe_arm64 }
    default {
        Write-Error "Unknown target: $target"
        exit 1
    }
}

#endregion
