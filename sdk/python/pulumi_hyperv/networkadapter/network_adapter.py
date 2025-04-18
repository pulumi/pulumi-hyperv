# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins
import copy
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from . import _utilities

__all__ = ['NetworkAdapterArgs', 'NetworkAdapter']

@pulumi.input_type
class NetworkAdapterArgs:
    def __init__(__self__, *,
                 name: pulumi.Input[builtins.str],
                 switch_name: pulumi.Input[builtins.str],
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 dhcp_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 ieee_priority_tag: Optional[pulumi.Input[builtins.bool]] = None,
                 ip_addresses: Optional[pulumi.Input[builtins.str]] = None,
                 mac_address: Optional[pulumi.Input[builtins.str]] = None,
                 port_mirroring: Optional[pulumi.Input[builtins.str]] = None,
                 router_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None,
                 vlan_id: Optional[pulumi.Input[builtins.int]] = None,
                 vm_name: Optional[pulumi.Input[builtins.str]] = None,
                 vmq_weight: Optional[pulumi.Input[builtins.int]] = None):
        """
        The set of arguments for constructing a NetworkAdapter resource.
        :param pulumi.Input[builtins.str] name: Name of the network adapter
        :param pulumi.Input[builtins.str] switch_name: Name of the virtual switch to connect the network adapter to
        :param pulumi.Input[builtins.str] create: The command to run on create.
        :param pulumi.Input[builtins.str] delete: The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
               and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
               Command resource from previous create or update steps.
        :param pulumi.Input[builtins.bool] dhcp_guard: Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
        :param pulumi.Input[builtins.bool] ieee_priority_tag: Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
        :param pulumi.Input[builtins.str] ip_addresses: Comma-separated list of IP addresses to assign to the network adapter.
        :param pulumi.Input[builtins.str] mac_address: MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
        :param pulumi.Input[builtins.str] port_mirroring: Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
        :param pulumi.Input[builtins.bool] router_guard: Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
        :param pulumi.Input[Sequence[Any]] triggers: Trigger a resource replacement on changes to any of these values. The
               trigger values can be of any type. If a value is different in the current update compared to the
               previous update, the resource will be replaced, i.e., the "create" command will be re-run.
               Please see the resource documentation for examples.
        :param pulumi.Input[builtins.str] update: The command to run on update, if empty, create will 
               run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
               are set to the stdout and stderr properties of the Command resource from previous 
               create or update steps.
        :param pulumi.Input[builtins.int] vlan_id: VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
        :param pulumi.Input[builtins.str] vm_name: Name of the virtual machine to attach the network adapter to
        :param pulumi.Input[builtins.int] vmq_weight: VMQ weight for the network adapter. A value of 0 disables VMQ.
        """
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "switch_name", switch_name)
        if create is not None:
            pulumi.set(__self__, "create", create)
        if delete is not None:
            pulumi.set(__self__, "delete", delete)
        if dhcp_guard is not None:
            pulumi.set(__self__, "dhcp_guard", dhcp_guard)
        if ieee_priority_tag is not None:
            pulumi.set(__self__, "ieee_priority_tag", ieee_priority_tag)
        if ip_addresses is not None:
            pulumi.set(__self__, "ip_addresses", ip_addresses)
        if mac_address is not None:
            pulumi.set(__self__, "mac_address", mac_address)
        if port_mirroring is not None:
            pulumi.set(__self__, "port_mirroring", port_mirroring)
        if router_guard is not None:
            pulumi.set(__self__, "router_guard", router_guard)
        if triggers is not None:
            pulumi.set(__self__, "triggers", triggers)
        if update is not None:
            pulumi.set(__self__, "update", update)
        if vlan_id is not None:
            pulumi.set(__self__, "vlan_id", vlan_id)
        if vm_name is not None:
            pulumi.set(__self__, "vm_name", vm_name)
        if vmq_weight is not None:
            pulumi.set(__self__, "vmq_weight", vmq_weight)

    @property
    @pulumi.getter
    def name(self) -> pulumi.Input[builtins.str]:
        """
        Name of the network adapter
        """
        return pulumi.get(self, "name")

    @name.setter
    def name(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "name", value)

    @property
    @pulumi.getter(name="switchName")
    def switch_name(self) -> pulumi.Input[builtins.str]:
        """
        Name of the virtual switch to connect the network adapter to
        """
        return pulumi.get(self, "switch_name")

    @switch_name.setter
    def switch_name(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "switch_name", value)

    @property
    @pulumi.getter
    def create(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        The command to run on create.
        """
        return pulumi.get(self, "create")

    @create.setter
    def create(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "create", value)

    @property
    @pulumi.getter
    def delete(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
        and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
        Command resource from previous create or update steps.
        """
        return pulumi.get(self, "delete")

    @delete.setter
    def delete(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "delete", value)

    @property
    @pulumi.getter(name="dhcpGuard")
    def dhcp_guard(self) -> Optional[pulumi.Input[builtins.bool]]:
        """
        Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
        """
        return pulumi.get(self, "dhcp_guard")

    @dhcp_guard.setter
    def dhcp_guard(self, value: Optional[pulumi.Input[builtins.bool]]):
        pulumi.set(self, "dhcp_guard", value)

    @property
    @pulumi.getter(name="ieeePriorityTag")
    def ieee_priority_tag(self) -> Optional[pulumi.Input[builtins.bool]]:
        """
        Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
        """
        return pulumi.get(self, "ieee_priority_tag")

    @ieee_priority_tag.setter
    def ieee_priority_tag(self, value: Optional[pulumi.Input[builtins.bool]]):
        pulumi.set(self, "ieee_priority_tag", value)

    @property
    @pulumi.getter(name="ipAddresses")
    def ip_addresses(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        Comma-separated list of IP addresses to assign to the network adapter.
        """
        return pulumi.get(self, "ip_addresses")

    @ip_addresses.setter
    def ip_addresses(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "ip_addresses", value)

    @property
    @pulumi.getter(name="macAddress")
    def mac_address(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
        """
        return pulumi.get(self, "mac_address")

    @mac_address.setter
    def mac_address(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "mac_address", value)

    @property
    @pulumi.getter(name="portMirroring")
    def port_mirroring(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
        """
        return pulumi.get(self, "port_mirroring")

    @port_mirroring.setter
    def port_mirroring(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "port_mirroring", value)

    @property
    @pulumi.getter(name="routerGuard")
    def router_guard(self) -> Optional[pulumi.Input[builtins.bool]]:
        """
        Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
        """
        return pulumi.get(self, "router_guard")

    @router_guard.setter
    def router_guard(self, value: Optional[pulumi.Input[builtins.bool]]):
        pulumi.set(self, "router_guard", value)

    @property
    @pulumi.getter
    def triggers(self) -> Optional[pulumi.Input[Sequence[Any]]]:
        """
        Trigger a resource replacement on changes to any of these values. The
        trigger values can be of any type. If a value is different in the current update compared to the
        previous update, the resource will be replaced, i.e., the "create" command will be re-run.
        Please see the resource documentation for examples.
        """
        return pulumi.get(self, "triggers")

    @triggers.setter
    def triggers(self, value: Optional[pulumi.Input[Sequence[Any]]]):
        pulumi.set(self, "triggers", value)

    @property
    @pulumi.getter
    def update(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        The command to run on update, if empty, create will 
        run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
        are set to the stdout and stderr properties of the Command resource from previous 
        create or update steps.
        """
        return pulumi.get(self, "update")

    @update.setter
    def update(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "update", value)

    @property
    @pulumi.getter(name="vlanId")
    def vlan_id(self) -> Optional[pulumi.Input[builtins.int]]:
        """
        VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
        """
        return pulumi.get(self, "vlan_id")

    @vlan_id.setter
    def vlan_id(self, value: Optional[pulumi.Input[builtins.int]]):
        pulumi.set(self, "vlan_id", value)

    @property
    @pulumi.getter(name="vmName")
    def vm_name(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        Name of the virtual machine to attach the network adapter to
        """
        return pulumi.get(self, "vm_name")

    @vm_name.setter
    def vm_name(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "vm_name", value)

    @property
    @pulumi.getter(name="vmqWeight")
    def vmq_weight(self) -> Optional[pulumi.Input[builtins.int]]:
        """
        VMQ weight for the network adapter. A value of 0 disables VMQ.
        """
        return pulumi.get(self, "vmq_weight")

    @vmq_weight.setter
    def vmq_weight(self, value: Optional[pulumi.Input[builtins.int]]):
        pulumi.set(self, "vmq_weight", value)


class NetworkAdapter(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 dhcp_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 ieee_priority_tag: Optional[pulumi.Input[builtins.bool]] = None,
                 ip_addresses: Optional[pulumi.Input[builtins.str]] = None,
                 mac_address: Optional[pulumi.Input[builtins.str]] = None,
                 name: Optional[pulumi.Input[builtins.str]] = None,
                 port_mirroring: Optional[pulumi.Input[builtins.str]] = None,
                 router_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 switch_name: Optional[pulumi.Input[builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None,
                 vlan_id: Optional[pulumi.Input[builtins.int]] = None,
                 vm_name: Optional[pulumi.Input[builtins.str]] = None,
                 vmq_weight: Optional[pulumi.Input[builtins.int]] = None,
                 __props__=None):
        """
        # Network Adapter Resource

        The Network Adapter resource allows you to create and manage network adapters for virtual machines in Hyper-V.

        ## Example Usage

        ### Standalone Network Adapter

        ### Using the NetworkAdapters Property in Machine Resource

        You can also define network adapters directly in the Machine resource using the `networkAdapters` property:

        ## Input Properties

        | Property         | Type     | Required | Description |
        |------------------|----------|----------|-------------|
        | name             | string   | Yes      | Name of the network adapter |
        | vmName           | string   | Yes      | Name of the virtual machine to attach the network adapter to |
        | switchName       | string   | Yes      | Name of the virtual switch to connect the network adapter to |
        | macAddress       | string   | No       | MAC address for the network adapter. If not specified, a dynamic MAC address will be generated |
        | vlanId           | number   | No       | VLAN ID for the network adapter. If not specified, no VLAN tagging is used |
        | dhcpGuard        | boolean  | No       | Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages |
        | routerGuard      | boolean  | No       | Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages |
        | portMirroring    | string   | No       | Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None |
        | ieeePriorityTag  | boolean  | No       | Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value |
        | vmqWeight        | number   | No       | VMQ weight for the network adapter. A value of 0 disables VMQ |
        | ipAddresses      | string   | No       | Comma-separated list of IP addresses to assign to the network adapter |

        ## Output Properties

        | Property         | Type     | Description |
        |------------------|----------|-------------|
        | adapterId        | string   | The ID of the network adapter |

        ## Lifecycle Management

        - **Create**: Creates a new network adapter and attaches it to the specified virtual machine.
        - **Read**: Reads the properties of an existing network adapter.
        - **Update**: Updates the properties of an existing network adapter.
        - **Delete**: Removes a network adapter from a virtual machine.

        ## Notes

        - The network adapter creation will fail if the virtual machine or virtual switch does not exist.
        - Dynamic MAC addresses are automatically generated if not specified.
        - IP addresses are specified as a comma-separated string (e.g., "192.168.1.10,192.168.1.11").
        - When updating a network adapter, the virtual machine may need to be powered off depending on the properties being changed.

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[builtins.str] create: The command to run on create.
        :param pulumi.Input[builtins.str] delete: The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
               and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
               Command resource from previous create or update steps.
        :param pulumi.Input[builtins.bool] dhcp_guard: Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
        :param pulumi.Input[builtins.bool] ieee_priority_tag: Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
        :param pulumi.Input[builtins.str] ip_addresses: Comma-separated list of IP addresses to assign to the network adapter.
        :param pulumi.Input[builtins.str] mac_address: MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
        :param pulumi.Input[builtins.str] name: Name of the network adapter
        :param pulumi.Input[builtins.str] port_mirroring: Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
        :param pulumi.Input[builtins.bool] router_guard: Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
        :param pulumi.Input[builtins.str] switch_name: Name of the virtual switch to connect the network adapter to
        :param pulumi.Input[Sequence[Any]] triggers: Trigger a resource replacement on changes to any of these values. The
               trigger values can be of any type. If a value is different in the current update compared to the
               previous update, the resource will be replaced, i.e., the "create" command will be re-run.
               Please see the resource documentation for examples.
        :param pulumi.Input[builtins.str] update: The command to run on update, if empty, create will 
               run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
               are set to the stdout and stderr properties of the Command resource from previous 
               create or update steps.
        :param pulumi.Input[builtins.int] vlan_id: VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
        :param pulumi.Input[builtins.str] vm_name: Name of the virtual machine to attach the network adapter to
        :param pulumi.Input[builtins.int] vmq_weight: VMQ weight for the network adapter. A value of 0 disables VMQ.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: NetworkAdapterArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        # Network Adapter Resource

        The Network Adapter resource allows you to create and manage network adapters for virtual machines in Hyper-V.

        ## Example Usage

        ### Standalone Network Adapter

        ### Using the NetworkAdapters Property in Machine Resource

        You can also define network adapters directly in the Machine resource using the `networkAdapters` property:

        ## Input Properties

        | Property         | Type     | Required | Description |
        |------------------|----------|----------|-------------|
        | name             | string   | Yes      | Name of the network adapter |
        | vmName           | string   | Yes      | Name of the virtual machine to attach the network adapter to |
        | switchName       | string   | Yes      | Name of the virtual switch to connect the network adapter to |
        | macAddress       | string   | No       | MAC address for the network adapter. If not specified, a dynamic MAC address will be generated |
        | vlanId           | number   | No       | VLAN ID for the network adapter. If not specified, no VLAN tagging is used |
        | dhcpGuard        | boolean  | No       | Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages |
        | routerGuard      | boolean  | No       | Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages |
        | portMirroring    | string   | No       | Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None |
        | ieeePriorityTag  | boolean  | No       | Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value |
        | vmqWeight        | number   | No       | VMQ weight for the network adapter. A value of 0 disables VMQ |
        | ipAddresses      | string   | No       | Comma-separated list of IP addresses to assign to the network adapter |

        ## Output Properties

        | Property         | Type     | Description |
        |------------------|----------|-------------|
        | adapterId        | string   | The ID of the network adapter |

        ## Lifecycle Management

        - **Create**: Creates a new network adapter and attaches it to the specified virtual machine.
        - **Read**: Reads the properties of an existing network adapter.
        - **Update**: Updates the properties of an existing network adapter.
        - **Delete**: Removes a network adapter from a virtual machine.

        ## Notes

        - The network adapter creation will fail if the virtual machine or virtual switch does not exist.
        - Dynamic MAC addresses are automatically generated if not specified.
        - IP addresses are specified as a comma-separated string (e.g., "192.168.1.10,192.168.1.11").
        - When updating a network adapter, the virtual machine may need to be powered off depending on the properties being changed.

        :param str resource_name: The name of the resource.
        :param NetworkAdapterArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(NetworkAdapterArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 dhcp_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 ieee_priority_tag: Optional[pulumi.Input[builtins.bool]] = None,
                 ip_addresses: Optional[pulumi.Input[builtins.str]] = None,
                 mac_address: Optional[pulumi.Input[builtins.str]] = None,
                 name: Optional[pulumi.Input[builtins.str]] = None,
                 port_mirroring: Optional[pulumi.Input[builtins.str]] = None,
                 router_guard: Optional[pulumi.Input[builtins.bool]] = None,
                 switch_name: Optional[pulumi.Input[builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None,
                 vlan_id: Optional[pulumi.Input[builtins.int]] = None,
                 vm_name: Optional[pulumi.Input[builtins.str]] = None,
                 vmq_weight: Optional[pulumi.Input[builtins.int]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = NetworkAdapterArgs.__new__(NetworkAdapterArgs)

            __props__.__dict__["create"] = create
            __props__.__dict__["delete"] = delete
            __props__.__dict__["dhcp_guard"] = dhcp_guard
            __props__.__dict__["ieee_priority_tag"] = ieee_priority_tag
            __props__.__dict__["ip_addresses"] = ip_addresses
            __props__.__dict__["mac_address"] = mac_address
            if name is None and not opts.urn:
                raise TypeError("Missing required property 'name'")
            __props__.__dict__["name"] = name
            __props__.__dict__["port_mirroring"] = port_mirroring
            __props__.__dict__["router_guard"] = router_guard
            if switch_name is None and not opts.urn:
                raise TypeError("Missing required property 'switch_name'")
            __props__.__dict__["switch_name"] = switch_name
            __props__.__dict__["triggers"] = triggers
            __props__.__dict__["update"] = update
            __props__.__dict__["vlan_id"] = vlan_id
            __props__.__dict__["vm_name"] = vm_name
            __props__.__dict__["vmq_weight"] = vmq_weight
            __props__.__dict__["adapter_id"] = None
        replace_on_changes = pulumi.ResourceOptions(replace_on_changes=["triggers[*]"])
        opts = pulumi.ResourceOptions.merge(opts, replace_on_changes)
        super(NetworkAdapter, __self__).__init__(
            'hyperv:networkadapter:NetworkAdapter',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'NetworkAdapter':
        """
        Get an existing NetworkAdapter resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = NetworkAdapterArgs.__new__(NetworkAdapterArgs)

        __props__.__dict__["adapter_id"] = None
        __props__.__dict__["create"] = None
        __props__.__dict__["delete"] = None
        __props__.__dict__["dhcp_guard"] = None
        __props__.__dict__["ieee_priority_tag"] = None
        __props__.__dict__["ip_addresses"] = None
        __props__.__dict__["mac_address"] = None
        __props__.__dict__["name"] = None
        __props__.__dict__["port_mirroring"] = None
        __props__.__dict__["router_guard"] = None
        __props__.__dict__["switch_name"] = None
        __props__.__dict__["triggers"] = None
        __props__.__dict__["update"] = None
        __props__.__dict__["vlan_id"] = None
        __props__.__dict__["vm_name"] = None
        __props__.__dict__["vmq_weight"] = None
        return NetworkAdapter(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter(name="adapterId")
    def adapter_id(self) -> pulumi.Output[builtins.str]:
        """
        The ID of the network adapter
        """
        return pulumi.get(self, "adapter_id")

    @property
    @pulumi.getter
    def create(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        The command to run on create.
        """
        return pulumi.get(self, "create")

    @property
    @pulumi.getter
    def delete(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
        and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
        Command resource from previous create or update steps.
        """
        return pulumi.get(self, "delete")

    @property
    @pulumi.getter(name="dhcpGuard")
    def dhcp_guard(self) -> pulumi.Output[Optional[builtins.bool]]:
        """
        Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
        """
        return pulumi.get(self, "dhcp_guard")

    @property
    @pulumi.getter(name="ieeePriorityTag")
    def ieee_priority_tag(self) -> pulumi.Output[Optional[builtins.bool]]:
        """
        Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
        """
        return pulumi.get(self, "ieee_priority_tag")

    @property
    @pulumi.getter(name="ipAddresses")
    def ip_addresses(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        Comma-separated list of IP addresses to assign to the network adapter.
        """
        return pulumi.get(self, "ip_addresses")

    @property
    @pulumi.getter(name="macAddress")
    def mac_address(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
        """
        return pulumi.get(self, "mac_address")

    @property
    @pulumi.getter
    def name(self) -> pulumi.Output[builtins.str]:
        """
        Name of the network adapter
        """
        return pulumi.get(self, "name")

    @property
    @pulumi.getter(name="portMirroring")
    def port_mirroring(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
        """
        return pulumi.get(self, "port_mirroring")

    @property
    @pulumi.getter(name="routerGuard")
    def router_guard(self) -> pulumi.Output[Optional[builtins.bool]]:
        """
        Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
        """
        return pulumi.get(self, "router_guard")

    @property
    @pulumi.getter(name="switchName")
    def switch_name(self) -> pulumi.Output[builtins.str]:
        """
        Name of the virtual switch to connect the network adapter to
        """
        return pulumi.get(self, "switch_name")

    @property
    @pulumi.getter
    def triggers(self) -> pulumi.Output[Optional[Sequence[Any]]]:
        """
        Trigger a resource replacement on changes to any of these values. The
        trigger values can be of any type. If a value is different in the current update compared to the
        previous update, the resource will be replaced, i.e., the "create" command will be re-run.
        Please see the resource documentation for examples.
        """
        return pulumi.get(self, "triggers")

    @property
    @pulumi.getter
    def update(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        The command to run on update, if empty, create will 
        run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
        are set to the stdout and stderr properties of the Command resource from previous 
        create or update steps.
        """
        return pulumi.get(self, "update")

    @property
    @pulumi.getter(name="vlanId")
    def vlan_id(self) -> pulumi.Output[Optional[builtins.int]]:
        """
        VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
        """
        return pulumi.get(self, "vlan_id")

    @property
    @pulumi.getter(name="vmName")
    def vm_name(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        Name of the virtual machine to attach the network adapter to
        """
        return pulumi.get(self, "vm_name")

    @property
    @pulumi.getter(name="vmqWeight")
    def vmq_weight(self) -> pulumi.Output[Optional[builtins.int]]:
        """
        VMQ weight for the network adapter. A value of 0 disables VMQ.
        """
        return pulumi.get(self, "vmq_weight")

