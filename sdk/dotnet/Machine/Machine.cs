// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Hyperv.Machine
{
    /// <summary>
    /// # Hyper-V Virtual Machine Management Service (VMMS)
    /// 
    /// ## Overview
    /// 
    /// The Virtual Machine Management Service (VMMS) is a core component of Hyper-V that manages virtual machine operations on a Windows Server or Windows Client system. This document provides information about the VMMS as implemented in the Pulumi Hyper-V provider.
    /// 
    /// ## Features
    /// 
    /// - Virtual machine lifecycle management (create, start, stop, pause, resume, delete)
    /// - Resource allocation and monitoring
    /// - Snapshot management
    /// - Virtual device configuration
    /// 
    /// ## Implementation Details in Pulumi
    /// 
    /// ### Virtual Machine Creation
    /// 
    /// The `Create` method in the `vmController` is responsible for creating a virtual machine. It performs the following steps:
    /// 
    /// 1. **Generate a Unique ID**: A unique ID is generated for the virtual machine.
    /// 2. **Default Values**:
    ///    - Memory size defaults to `1024 MB` if not specified.
    ///    - Processor count defaults to `1` if not specified.
    /// 3. **VMMS Client Initialization**: A VMMS client is created to interact with the Hyper-V host.
    /// 4. **Virtual Machine Settings**:
    ///    - The virtual machine is configured with `Hyper-V Generation 2`.
    ///    - Memory and processor settings are applied based on the provided or default values.
    /// 5. **Virtual Machine Creation**: The virtual machine is created using the configured settings.
    /// 
    /// ### Read Method
    /// 
    /// The `Read` method is a no-op in the current implementation. It does not perform any operations and always returns an empty state.
    /// 
    /// ### Update Method
    /// 
    /// The `Update` method:
    /// 
    /// - Updates the virtual machine state if an `Update` command is provided.
    /// - Falls back to the `Create` command if no `Update` command is specified.
    /// - Does nothing if neither command is provided.
    /// 
    /// ### Delete Method
    /// 
    /// The `Delete` method is a no-op unless a `Delete` command is explicitly specified.
    /// 
    /// ## Default Behavior
    /// 
    /// - Outputs depend on all inputs by default.
    /// - No explicit dependency wiring is required.
    /// 
    /// ## Usage in Pulumi
    /// 
    /// When using the Pulumi Hyper-V provider, the VMMS is accessed indirectly through the `Vm` resource type. The resource supports the following properties:
    /// 
    /// - `processorCount`: Number of processors to allocate (default: 1).
    /// - `memorySize`: Memory size in MB (default: 1024).
    /// 
    /// ## Authentication and Security
    /// 
    /// The VMMS requires appropriate permissions to manage Hyper-V objects. When using the Pulumi Hyper-V provider, ensure that:
    /// 
    /// 1. The user running Pulumi commands has administrative privileges on the Hyper-V host.
    /// 2. Required firewall rules are configured if managing a remote Hyper-V host.
    /// 3. Proper credentials are provided when connecting to remote systems.
    /// 
    /// ## Related Documentation
    /// 
    /// - [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
    /// - [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/
    /// </summary>
    [HypervResourceType("hyperv:machine:Machine")]
    public partial class Machine : global::Pulumi.CustomResource
    {
        /// <summary>
        /// The command to run on create.
        /// </summary>
        [Output("create")]
        public Output<string?> Create { get; private set; } = null!;

        /// <summary>
        /// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
        /// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
        /// Command resource from previous create or update steps.
        /// </summary>
        [Output("delete")]
        public Output<string?> Delete { get; private set; } = null!;

        /// <summary>
        /// Name of the Virtual Machine
        /// </summary>
        [Output("machineName")]
        public Output<string> MachineName { get; private set; } = null!;

        /// <summary>
        /// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
        /// </summary>
        [Output("memorySize")]
        public Output<int?> MemorySize { get; private set; } = null!;

        /// <summary>
        /// Number of processors to allocate to the Virtual Machine. Defaults to 1.
        /// </summary>
        [Output("processorCount")]
        public Output<int?> ProcessorCount { get; private set; } = null!;

        /// <summary>
        /// Trigger a resource replacement on changes to any of these values. The
        /// trigger values can be of any type. If a value is different in the current update compared to the
        /// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
        /// Please see the resource documentation for examples.
        /// </summary>
        [Output("triggers")]
        public Output<ImmutableArray<object>> Triggers { get; private set; } = null!;

        /// <summary>
        /// The command to run on update, if empty, create will 
        /// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
        /// are set to the stdout and stderr properties of the Command resource from previous 
        /// create or update steps.
        /// </summary>
        [Output("update")]
        public Output<string?> Update { get; private set; } = null!;


        /// <summary>
        /// Create a Machine resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Machine(string name, MachineArgs args, CustomResourceOptions? options = null)
            : base("hyperv:machine:Machine", name, args ?? new MachineArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Machine(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("hyperv:machine:Machine", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
                ReplaceOnChanges =
                {
                    "triggers[*]",
                },
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing Machine resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Machine Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Machine(name, id, options);
        }
    }

    public sealed class MachineArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// The command to run on create.
        /// </summary>
        [Input("create")]
        public Input<string>? Create { get; set; }

        /// <summary>
        /// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
        /// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
        /// Command resource from previous create or update steps.
        /// </summary>
        [Input("delete")]
        public Input<string>? Delete { get; set; }

        /// <summary>
        /// Name of the Virtual Machine
        /// </summary>
        [Input("machineName", required: true)]
        public Input<string> MachineName { get; set; } = null!;

        /// <summary>
        /// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
        /// </summary>
        [Input("memorySize")]
        public Input<int>? MemorySize { get; set; }

        /// <summary>
        /// Number of processors to allocate to the Virtual Machine. Defaults to 1.
        /// </summary>
        [Input("processorCount")]
        public Input<int>? ProcessorCount { get; set; }

        [Input("triggers")]
        private InputList<object>? _triggers;

        /// <summary>
        /// Trigger a resource replacement on changes to any of these values. The
        /// trigger values can be of any type. If a value is different in the current update compared to the
        /// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
        /// Please see the resource documentation for examples.
        /// </summary>
        public InputList<object> Triggers
        {
            get => _triggers ?? (_triggers = new InputList<object>());
            set => _triggers = value;
        }

        /// <summary>
        /// The command to run on update, if empty, create will 
        /// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
        /// are set to the stdout and stderr properties of the Command resource from previous 
        /// create or update steps.
        /// </summary>
        [Input("update")]
        public Input<string>? Update { get; set; }

        public MachineArgs()
        {
        }
        public static new MachineArgs Empty => new MachineArgs();
    }
}
