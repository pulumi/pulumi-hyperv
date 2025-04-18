// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Hyperv.Networkadapter.Outputs
{

    [OutputType]
    public sealed class NetworkAdapterInputs
    {
        /// <summary>
        /// The command to run on create.
        /// </summary>
        public readonly string? Create;
        /// <summary>
        /// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
        /// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
        /// Command resource from previous create or update steps.
        /// </summary>
        public readonly string? Delete;
        /// <summary>
        /// Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
        /// </summary>
        public readonly bool? DhcpGuard;
        /// <summary>
        /// Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
        /// </summary>
        public readonly bool? IeeePriorityTag;
        /// <summary>
        /// Comma-separated list of IP addresses to assign to the network adapter.
        /// </summary>
        public readonly string? IpAddresses;
        /// <summary>
        /// MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
        /// </summary>
        public readonly string? MacAddress;
        /// <summary>
        /// Name of the network adapter
        /// </summary>
        public readonly string Name;
        /// <summary>
        /// Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
        /// </summary>
        public readonly string? PortMirroring;
        /// <summary>
        /// Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
        /// </summary>
        public readonly bool? RouterGuard;
        /// <summary>
        /// Name of the virtual switch to connect the network adapter to
        /// </summary>
        public readonly string SwitchName;
        /// <summary>
        /// Trigger a resource replacement on changes to any of these values. The
        /// trigger values can be of any type. If a value is different in the current update compared to the
        /// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
        /// Please see the resource documentation for examples.
        /// </summary>
        public readonly ImmutableArray<object> Triggers;
        /// <summary>
        /// The command to run on update, if empty, create will 
        /// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
        /// are set to the stdout and stderr properties of the Command resource from previous 
        /// create or update steps.
        /// </summary>
        public readonly string? Update;
        /// <summary>
        /// VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
        /// </summary>
        public readonly int? VlanId;
        /// <summary>
        /// Name of the virtual machine to attach the network adapter to
        /// </summary>
        public readonly string? VmName;
        /// <summary>
        /// VMQ weight for the network adapter. A value of 0 disables VMQ.
        /// </summary>
        public readonly int? VmqWeight;

        [OutputConstructor]
        private NetworkAdapterInputs(
            string? create,

            string? delete,

            bool? dhcpGuard,

            bool? ieeePriorityTag,

            string? ipAddresses,

            string? macAddress,

            string name,

            string? portMirroring,

            bool? routerGuard,

            string switchName,

            ImmutableArray<object> triggers,

            string? update,

            int? vlanId,

            string? vmName,

            int? vmqWeight)
        {
            Create = create;
            Delete = delete;
            DhcpGuard = dhcpGuard;
            IeeePriorityTag = ieeePriorityTag;
            IpAddresses = ipAddresses;
            MacAddress = macAddress;
            Name = name;
            PortMirroring = portMirroring;
            RouterGuard = routerGuard;
            SwitchName = switchName;
            Triggers = triggers;
            Update = update;
            VlanId = vlanId;
            VmName = vmName;
            VmqWeight = vmqWeight;
        }
    }
}
