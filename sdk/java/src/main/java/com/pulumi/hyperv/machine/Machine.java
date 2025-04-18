// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.hyperv.machine;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.hyperv.Utilities;
import com.pulumi.hyperv.machine.MachineArgs;
import com.pulumi.hyperv.machine.outputs.HardDriveInput;
import com.pulumi.hyperv.networkadapter.outputs.NetworkAdapterInputs;
import java.lang.Boolean;
import java.lang.Integer;
import java.lang.Object;
import java.lang.String;
import java.util.List;
import java.util.Optional;
import javax.annotation.Nullable;

/**
 * # Hyper-V Machine Resource
 * 
 * ## Overview
 * 
 * The Machine resource in the Pulumi Hyper-V provider allows you to create, manage, and delete virtual machines on a Hyper-V host. This resource interacts with the Virtual Machine Management Service (VMMS) to perform virtual machine operations.
 * 
 * ## Features
 * 
 * - Create and delete Hyper-V virtual machines
 * - Configure VM hardware properties including:
 *   - Memory allocation (static or dynamic with min/max)
 *   - Processor count
 *   - VM generation (Gen 1 or Gen 2)
 *   - Auto start/stop actions
 * - Attach hard drives with custom controller configuration
 * - Configure network adapters with virtual switch connections
 * - Unique VM identification with automatic ID generation
 * 
 * ## Implementation Details
 * 
 * ### Resource Structure
 * 
 * The Machine resource implementation consists of multiple files:
 * - `machine.go` - Core resource type definition, input/output models, and annotations
 * - `machineController.go` - Implementation of CRUD operations
 * - `machineOutputs.go` - Output-specific methods
 * 
 * ### Virtual Machine Creation
 * 
 * The `Create` method performs the following steps:
 * 
 * 1. **Initialize Connection**: Establishes a connection to the Hyper-V host using WMI
 * 2. **Configure VM Settings**:
 *    - Sets the virtual machine generation (defaults to Generation 2)
 *    - Configures memory settings (defaults to 1024 MB)
 *    - Sets dynamic memory with min/max values if requested
 *    - Sets processor count (defaults to 1 vCPU)
 *    - Configures auto start/stop actions
 * 3. **Create VM**: Calls the Hyper-V API to create a new virtual machine with the specified settings
 * 4. **Attach Hard Drives**: Attaches any specified hard drives to the VM
 * 5. **Configure Network Adapters**: Adds any specified network adapters to the VM
 * 
 * ### Virtual Machine Read
 * 
 * The `Read` method retrieves the current state of a virtual machine by:
 * 1. Connecting to the Hyper-V host
 * 2. Getting the VM by name
 * 3. Retrieving VM properties including:
 *    - VM ID
 *    - Memory settings (including dynamic memory configuration)
 *    - Processor configuration
 *    - Generation
 *    - Auto start/stop actions
 * 
 * ### Virtual Machine Update
 * 
 * The `Update` method currently provides a minimal implementation that preserves the VM&#39;s state while updating its metadata.
 * 
 * ### Virtual Machine Delete
 * 
 * The `Delete` method:
 * 1. Connects to the Hyper-V host
 * 2. Gets the virtual machine by name
 * 3. Starts the VM (to ensure it&#39;s in a state that can be properly deleted)
 * 4. Gracefully stops the VM
 * 5. Deletes the virtual machine
 * 
 * ## Available Properties
 * 
 * | Property | Type | Description | Default |
 * |----------|------|-------------|---------|
 * | `machineName` | string | Name of the Virtual Machine | (required) |
 * | `generation` | int | Generation of the Virtual Machine (1 or 2) | 2 |
 * | `processorCount` | int | Number of processors to allocate | 1 |
 * | `memorySize` | int | Memory size in MB | 1024 |
 * | `dynamicMemory` | bool | Enable dynamic memory for the VM | false |
 * | `minimumMemory` | int | Minimum memory in MB when using dynamic memory | - |
 * | `maximumMemory` | int | Maximum memory in MB when using dynamic memory | - |
 * | `autoStartAction` | string | Action on host start (Nothing, StartIfRunning, Start) | Nothing |
 * | `autoStopAction` | string | Action on host shutdown (TurnOff, Save, ShutDown) | TurnOff |
 * | `networkAdapters` | array | Network adapters to attach to the VM | [] |
 * | `hardDrives` | array | Hard drives to attach to the VM | [] |
 * | `triggers` | array | Values that trigger resource replacement when changed | (optional) |
 * 
 * ### Network Adapter Properties
 * 
 * | Property | Type | Description | Default |
 * |----------|------|-------------|---------|
 * | `name` | string | Name of the network adapter | &#34;Network Adapter&#34; |
 * | `switchName` | string | Name of the virtual switch to connect to | (required) |
 * 
 * ### Hard Drive Properties
 * 
 * | Property | Type | Description | Default |
 * |----------|------|-------------|---------|
 * | `path` | string | Path to the VHD/VHDX file | (required) |
 * | `controllerType` | string | Type of controller (IDE or SCSI) | SCSI |
 * | `controllerNumber` | int | Controller number | 0 |
 * | `controllerLocation` | int | Controller location | 0 |
 * 
 * ## Usage Examples
 * 
 * ## Related Documentation
 * 
 * - [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
 * - [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/)
 * 
 */
@ResourceType(type="hyperv:machine:Machine")
public class Machine extends com.pulumi.resources.CustomResource {
    /**
     * The action to take when the host starts. Valid values are Nothing, StartIfRunning, and Start. Defaults to Nothing.
     * 
     */
    @Export(name="autoStartAction", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> autoStartAction;

    /**
     * @return The action to take when the host starts. Valid values are Nothing, StartIfRunning, and Start. Defaults to Nothing.
     * 
     */
    public Output<Optional<String>> autoStartAction() {
        return Codegen.optional(this.autoStartAction);
    }
    /**
     * The action to take when the host shuts down. Valid values are TurnOff, Save, and ShutDown. Defaults to TurnOff.
     * 
     */
    @Export(name="autoStopAction", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> autoStopAction;

    /**
     * @return The action to take when the host shuts down. Valid values are TurnOff, Save, and ShutDown. Defaults to TurnOff.
     * 
     */
    public Output<Optional<String>> autoStopAction() {
        return Codegen.optional(this.autoStopAction);
    }
    /**
     * The command to run on create.
     * 
     */
    @Export(name="create", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> create;

    /**
     * @return The command to run on create.
     * 
     */
    public Output<Optional<String>> create() {
        return Codegen.optional(this.create);
    }
    /**
     * The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
     * and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
     * Command resource from previous create or update steps.
     * 
     */
    @Export(name="delete", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> delete;

    /**
     * @return The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
     * and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
     * Command resource from previous create or update steps.
     * 
     */
    public Output<Optional<String>> delete() {
        return Codegen.optional(this.delete);
    }
    /**
     * Whether to enable dynamic memory for the Virtual Machine. Defaults to false.
     * 
     */
    @Export(name="dynamicMemory", refs={Boolean.class}, tree="[0]")
    private Output</* @Nullable */ Boolean> dynamicMemory;

    /**
     * @return Whether to enable dynamic memory for the Virtual Machine. Defaults to false.
     * 
     */
    public Output<Optional<Boolean>> dynamicMemory() {
        return Codegen.optional(this.dynamicMemory);
    }
    /**
     * Generation of the Virtual Machine. Defaults to 2.
     * 
     */
    @Export(name="generation", refs={Integer.class}, tree="[0]")
    private Output</* @Nullable */ Integer> generation;

    /**
     * @return Generation of the Virtual Machine. Defaults to 2.
     * 
     */
    public Output<Optional<Integer>> generation() {
        return Codegen.optional(this.generation);
    }
    /**
     * Hard drives to attach to the Virtual Machine.
     * 
     */
    @Export(name="hardDrives", refs={List.class,HardDriveInput.class}, tree="[0,1]")
    private Output</* @Nullable */ List<HardDriveInput>> hardDrives;

    /**
     * @return Hard drives to attach to the Virtual Machine.
     * 
     */
    public Output<Optional<List<HardDriveInput>>> hardDrives() {
        return Codegen.optional(this.hardDrives);
    }
    /**
     * Name of the Virtual Machine
     * 
     */
    @Export(name="machineName", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> machineName;

    /**
     * @return Name of the Virtual Machine
     * 
     */
    public Output<Optional<String>> machineName() {
        return Codegen.optional(this.machineName);
    }
    /**
     * Maximum amount of memory that can be allocated to the Virtual Machine in MB when using dynamic memory.
     * 
     */
    @Export(name="maximumMemory", refs={Integer.class}, tree="[0]")
    private Output</* @Nullable */ Integer> maximumMemory;

    /**
     * @return Maximum amount of memory that can be allocated to the Virtual Machine in MB when using dynamic memory.
     * 
     */
    public Output<Optional<Integer>> maximumMemory() {
        return Codegen.optional(this.maximumMemory);
    }
    /**
     * Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
     * 
     */
    @Export(name="memorySize", refs={Integer.class}, tree="[0]")
    private Output</* @Nullable */ Integer> memorySize;

    /**
     * @return Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
     * 
     */
    public Output<Optional<Integer>> memorySize() {
        return Codegen.optional(this.memorySize);
    }
    /**
     * Minimum amount of memory to allocate to the Virtual Machine in MB when using dynamic memory.
     * 
     */
    @Export(name="minimumMemory", refs={Integer.class}, tree="[0]")
    private Output</* @Nullable */ Integer> minimumMemory;

    /**
     * @return Minimum amount of memory to allocate to the Virtual Machine in MB when using dynamic memory.
     * 
     */
    public Output<Optional<Integer>> minimumMemory() {
        return Codegen.optional(this.minimumMemory);
    }
    /**
     * Network adapters to attach to the Virtual Machine.
     * 
     */
    @Export(name="networkAdapters", refs={List.class,NetworkAdapterInputs.class}, tree="[0,1]")
    private Output</* @Nullable */ List<NetworkAdapterInputs>> networkAdapters;

    /**
     * @return Network adapters to attach to the Virtual Machine.
     * 
     */
    public Output<Optional<List<NetworkAdapterInputs>>> networkAdapters() {
        return Codegen.optional(this.networkAdapters);
    }
    /**
     * Number of processors to allocate to the Virtual Machine. Defaults to 1.
     * 
     */
    @Export(name="processorCount", refs={Integer.class}, tree="[0]")
    private Output</* @Nullable */ Integer> processorCount;

    /**
     * @return Number of processors to allocate to the Virtual Machine. Defaults to 1.
     * 
     */
    public Output<Optional<Integer>> processorCount() {
        return Codegen.optional(this.processorCount);
    }
    /**
     * Trigger a resource replacement on changes to any of these values. The
     * trigger values can be of any type. If a value is different in the current update compared to the
     * previous update, the resource will be replaced, i.e., the &#34;create&#34; command will be re-run.
     * Please see the resource documentation for examples.
     * 
     */
    @Export(name="triggers", refs={List.class,Object.class}, tree="[0,1]")
    private Output</* @Nullable */ List<Object>> triggers;

    /**
     * @return Trigger a resource replacement on changes to any of these values. The
     * trigger values can be of any type. If a value is different in the current update compared to the
     * previous update, the resource will be replaced, i.e., the &#34;create&#34; command will be re-run.
     * Please see the resource documentation for examples.
     * 
     */
    public Output<Optional<List<Object>>> triggers() {
        return Codegen.optional(this.triggers);
    }
    /**
     * The command to run on update, if empty, create will
     * run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
     * are set to the stdout and stderr properties of the Command resource from previous
     * create or update steps.
     * 
     */
    @Export(name="update", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> update;

    /**
     * @return The command to run on update, if empty, create will
     * run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
     * are set to the stdout and stderr properties of the Command resource from previous
     * create or update steps.
     * 
     */
    public Output<Optional<String>> update() {
        return Codegen.optional(this.update);
    }
    @Export(name="vmId", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> vmId;

    public Output<Optional<String>> vmId() {
        return Codegen.optional(this.vmId);
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public Machine(java.lang.String name) {
        this(name, MachineArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public Machine(java.lang.String name, @Nullable MachineArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public Machine(java.lang.String name, @Nullable MachineArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("hyperv:machine:Machine", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private Machine(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("hyperv:machine:Machine", name, null, makeResourceOptions(options, id), false);
    }

    private static MachineArgs makeArgs(@Nullable MachineArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        return args == null ? MachineArgs.Empty : args;
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .build();
        return com.pulumi.resources.CustomResourceOptions.merge(defaultOptions, options, id);
    }

    /**
     * Get an existing Host resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param options Optional settings to control the behavior of the CustomResource.
     */
    public static Machine get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new Machine(name, id, options);
    }
}
