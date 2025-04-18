// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Hyperv.Virtualswitch
{
    /// <summary>
    /// # Virtual Switch Resource Management
    /// 
    /// The `virtualswitch` package provides utilities for managing Hyper-V virtual switches.
    /// 
    /// ## Overview
    /// 
    /// This package enables creating, modifying, and deleting virtual switches through the Pulumi Hyper-V provider. Virtual switches enable network connectivity for virtual machines.
    /// 
    /// ## Key Components
    /// 
    /// ### Types
    /// 
    /// - **VirtualSwitch**: Represents a Hyper-V virtual switch.
    /// 
    /// ### Resource Lifecycle Methods
    /// 
    /// - **Create**: Creates a new virtual switch with specified properties.
    /// - **Read**: Retrieves information about an existing virtual switch.
    /// - **Update**: Modifies properties of an existing virtual switch.
    /// - **Delete**: Removes a virtual switch.
    /// 
    /// ## Available Properties
    /// 
    /// The VirtualSwitch resource supports the following properties:
    /// 
    /// | Property | Type | Description |
    /// |----------|------|-------------|
    /// | `name` | string | Name of the virtual switch |
    /// | `switchType` | string | Type of switch: "External", "Internal", or "Private" |
    /// | `allowManagementOs` | boolean | Allow the management OS to access the switch (External switches) |
    /// | `netAdapterName` | string | Name of the physical network adapter to bind to (External switches) |
    /// 
    /// ## Implementation Details
    /// 
    /// The package uses the WMI interface to interact with Hyper-V's virtual switch management functionality, providing a Go-based interface that integrates with the Pulumi resource model.
    /// 
    /// ## Usage Examples
    /// 
    /// Virtual switches can be defined and managed through the Pulumi Hyper-V provider using the standard resource model.
    /// 
    /// ### Creating an External Switch
    /// 
    /// ### Creating an Internal Switch
    /// </summary>
    [HypervResourceType("hyperv:virtualswitch:VirtualSwitch")]
    public partial class VirtualSwitch : global::Pulumi.CustomResource
    {
        /// <summary>
        /// Allow the management OS to access the switch (External switches)
        /// </summary>
        [Output("allowManagementOs")]
        public Output<bool?> AllowManagementOs { get; private set; } = null!;

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
        /// Name of the virtual switch
        /// </summary>
        [Output("name")]
        public Output<string> Name { get; private set; } = null!;

        /// <summary>
        /// Name of the physical network adapter to bind to (External switches)
        /// </summary>
        [Output("netAdapterName")]
        public Output<string?> NetAdapterName { get; private set; } = null!;

        /// <summary>
        /// Notes or description for the virtual switch
        /// </summary>
        [Output("notes")]
        public Output<string?> Notes { get; private set; } = null!;

        /// <summary>
        /// Type of switch: 'External', 'Internal', or 'Private'
        /// </summary>
        [Output("switchType")]
        public Output<string> SwitchType { get; private set; } = null!;

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
        /// Create a VirtualSwitch resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public VirtualSwitch(string name, VirtualSwitchArgs args, CustomResourceOptions? options = null)
            : base("hyperv:virtualswitch:VirtualSwitch", name, args ?? new VirtualSwitchArgs(), MakeResourceOptions(options, ""))
        {
        }

        private VirtualSwitch(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("hyperv:virtualswitch:VirtualSwitch", name, null, MakeResourceOptions(options, id))
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
        /// Get an existing VirtualSwitch resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static VirtualSwitch Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new VirtualSwitch(name, id, options);
        }
    }

    public sealed class VirtualSwitchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Allow the management OS to access the switch (External switches)
        /// </summary>
        [Input("allowManagementOs")]
        public Input<bool>? AllowManagementOs { get; set; }

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
        /// Name of the virtual switch
        /// </summary>
        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        /// <summary>
        /// Name of the physical network adapter to bind to (External switches)
        /// </summary>
        [Input("netAdapterName")]
        public Input<string>? NetAdapterName { get; set; }

        /// <summary>
        /// Notes or description for the virtual switch
        /// </summary>
        [Input("notes")]
        public Input<string>? Notes { get; set; }

        /// <summary>
        /// Type of switch: 'External', 'Internal', or 'Private'
        /// </summary>
        [Input("switchType", required: true)]
        public Input<string> SwitchType { get; set; } = null!;

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

        public VirtualSwitchArgs()
        {
        }
        public static new VirtualSwitchArgs Empty => new VirtualSwitchArgs();
    }
}
