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

__all__ = ['VirtualSwitchArgs', 'VirtualSwitch']

@pulumi.input_type
class VirtualSwitchArgs:
    def __init__(__self__, *,
                 name: pulumi.Input[builtins.str],
                 switch_type: pulumi.Input[builtins.str],
                 allow_management_os: Optional[pulumi.Input[builtins.bool]] = None,
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 net_adapter_name: Optional[pulumi.Input[builtins.str]] = None,
                 notes: Optional[pulumi.Input[builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None):
        """
        The set of arguments for constructing a VirtualSwitch resource.
        :param pulumi.Input[builtins.str] name: Name of the virtual switch
        :param pulumi.Input[builtins.str] switch_type: Type of switch: 'External', 'Internal', or 'Private'
        :param pulumi.Input[builtins.bool] allow_management_os: Allow the management OS to access the switch (External switches)
        :param pulumi.Input[builtins.str] create: The command to run on create.
        :param pulumi.Input[builtins.str] delete: The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
               and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
               Command resource from previous create or update steps.
        :param pulumi.Input[builtins.str] net_adapter_name: Name of the physical network adapter to bind to (External switches)
        :param pulumi.Input[builtins.str] notes: Notes or description for the virtual switch
        :param pulumi.Input[Sequence[Any]] triggers: Trigger a resource replacement on changes to any of these values. The
               trigger values can be of any type. If a value is different in the current update compared to the
               previous update, the resource will be replaced, i.e., the "create" command will be re-run.
               Please see the resource documentation for examples.
        :param pulumi.Input[builtins.str] update: The command to run on update, if empty, create will 
               run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
               are set to the stdout and stderr properties of the Command resource from previous 
               create or update steps.
        """
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "switch_type", switch_type)
        if allow_management_os is not None:
            pulumi.set(__self__, "allow_management_os", allow_management_os)
        if create is not None:
            pulumi.set(__self__, "create", create)
        if delete is not None:
            pulumi.set(__self__, "delete", delete)
        if net_adapter_name is not None:
            pulumi.set(__self__, "net_adapter_name", net_adapter_name)
        if notes is not None:
            pulumi.set(__self__, "notes", notes)
        if triggers is not None:
            pulumi.set(__self__, "triggers", triggers)
        if update is not None:
            pulumi.set(__self__, "update", update)

    @property
    @pulumi.getter
    def name(self) -> pulumi.Input[builtins.str]:
        """
        Name of the virtual switch
        """
        return pulumi.get(self, "name")

    @name.setter
    def name(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "name", value)

    @property
    @pulumi.getter(name="switchType")
    def switch_type(self) -> pulumi.Input[builtins.str]:
        """
        Type of switch: 'External', 'Internal', or 'Private'
        """
        return pulumi.get(self, "switch_type")

    @switch_type.setter
    def switch_type(self, value: pulumi.Input[builtins.str]):
        pulumi.set(self, "switch_type", value)

    @property
    @pulumi.getter(name="allowManagementOs")
    def allow_management_os(self) -> Optional[pulumi.Input[builtins.bool]]:
        """
        Allow the management OS to access the switch (External switches)
        """
        return pulumi.get(self, "allow_management_os")

    @allow_management_os.setter
    def allow_management_os(self, value: Optional[pulumi.Input[builtins.bool]]):
        pulumi.set(self, "allow_management_os", value)

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
    @pulumi.getter(name="netAdapterName")
    def net_adapter_name(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        Name of the physical network adapter to bind to (External switches)
        """
        return pulumi.get(self, "net_adapter_name")

    @net_adapter_name.setter
    def net_adapter_name(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "net_adapter_name", value)

    @property
    @pulumi.getter
    def notes(self) -> Optional[pulumi.Input[builtins.str]]:
        """
        Notes or description for the virtual switch
        """
        return pulumi.get(self, "notes")

    @notes.setter
    def notes(self, value: Optional[pulumi.Input[builtins.str]]):
        pulumi.set(self, "notes", value)

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


class VirtualSwitch(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 allow_management_os: Optional[pulumi.Input[builtins.bool]] = None,
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 name: Optional[pulumi.Input[builtins.str]] = None,
                 net_adapter_name: Optional[pulumi.Input[builtins.str]] = None,
                 notes: Optional[pulumi.Input[builtins.str]] = None,
                 switch_type: Optional[pulumi.Input[builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None,
                 __props__=None):
        """
        # Virtual Switch Resource Management

        The `virtualswitch` package provides utilities for managing Hyper-V virtual switches.

        ## Overview

        This package enables creating, modifying, and deleting virtual switches through the Pulumi Hyper-V provider. Virtual switches enable network connectivity for virtual machines.

        ## Key Components

        ### Types

        - **VirtualSwitch**: Represents a Hyper-V virtual switch.

        ### Resource Lifecycle Methods

        - **Create**: Creates a new virtual switch with specified properties.
        - **Read**: Retrieves information about an existing virtual switch.
        - **Update**: Modifies properties of an existing virtual switch.
        - **Delete**: Removes a virtual switch.

        ## Available Properties

        The VirtualSwitch resource supports the following properties:

        | Property | Type | Description |
        |----------|------|-------------|
        | `name` | string | Name of the virtual switch |
        | `switchType` | string | Type of switch: "External", "Internal", or "Private" |
        | `allowManagementOs` | boolean | Allow the management OS to access the switch (External switches) |
        | `netAdapterName` | string | Name of the physical network adapter to bind to (External switches) |

        ## Implementation Details

        The package uses the WMI interface to interact with Hyper-V's virtual switch management functionality, providing a Go-based interface that integrates with the Pulumi resource model.

        ## Usage Examples

        Virtual switches can be defined and managed through the Pulumi Hyper-V provider using the standard resource model.

        ### Creating an External Switch

        ### Creating an Internal Switch

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[builtins.bool] allow_management_os: Allow the management OS to access the switch (External switches)
        :param pulumi.Input[builtins.str] create: The command to run on create.
        :param pulumi.Input[builtins.str] delete: The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
               and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
               Command resource from previous create or update steps.
        :param pulumi.Input[builtins.str] name: Name of the virtual switch
        :param pulumi.Input[builtins.str] net_adapter_name: Name of the physical network adapter to bind to (External switches)
        :param pulumi.Input[builtins.str] notes: Notes or description for the virtual switch
        :param pulumi.Input[builtins.str] switch_type: Type of switch: 'External', 'Internal', or 'Private'
        :param pulumi.Input[Sequence[Any]] triggers: Trigger a resource replacement on changes to any of these values. The
               trigger values can be of any type. If a value is different in the current update compared to the
               previous update, the resource will be replaced, i.e., the "create" command will be re-run.
               Please see the resource documentation for examples.
        :param pulumi.Input[builtins.str] update: The command to run on update, if empty, create will 
               run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
               are set to the stdout and stderr properties of the Command resource from previous 
               create or update steps.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: VirtualSwitchArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        # Virtual Switch Resource Management

        The `virtualswitch` package provides utilities for managing Hyper-V virtual switches.

        ## Overview

        This package enables creating, modifying, and deleting virtual switches through the Pulumi Hyper-V provider. Virtual switches enable network connectivity for virtual machines.

        ## Key Components

        ### Types

        - **VirtualSwitch**: Represents a Hyper-V virtual switch.

        ### Resource Lifecycle Methods

        - **Create**: Creates a new virtual switch with specified properties.
        - **Read**: Retrieves information about an existing virtual switch.
        - **Update**: Modifies properties of an existing virtual switch.
        - **Delete**: Removes a virtual switch.

        ## Available Properties

        The VirtualSwitch resource supports the following properties:

        | Property | Type | Description |
        |----------|------|-------------|
        | `name` | string | Name of the virtual switch |
        | `switchType` | string | Type of switch: "External", "Internal", or "Private" |
        | `allowManagementOs` | boolean | Allow the management OS to access the switch (External switches) |
        | `netAdapterName` | string | Name of the physical network adapter to bind to (External switches) |

        ## Implementation Details

        The package uses the WMI interface to interact with Hyper-V's virtual switch management functionality, providing a Go-based interface that integrates with the Pulumi resource model.

        ## Usage Examples

        Virtual switches can be defined and managed through the Pulumi Hyper-V provider using the standard resource model.

        ### Creating an External Switch

        ### Creating an Internal Switch

        :param str resource_name: The name of the resource.
        :param VirtualSwitchArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(VirtualSwitchArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 allow_management_os: Optional[pulumi.Input[builtins.bool]] = None,
                 create: Optional[pulumi.Input[builtins.str]] = None,
                 delete: Optional[pulumi.Input[builtins.str]] = None,
                 name: Optional[pulumi.Input[builtins.str]] = None,
                 net_adapter_name: Optional[pulumi.Input[builtins.str]] = None,
                 notes: Optional[pulumi.Input[builtins.str]] = None,
                 switch_type: Optional[pulumi.Input[builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[builtins.str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = VirtualSwitchArgs.__new__(VirtualSwitchArgs)

            __props__.__dict__["allow_management_os"] = allow_management_os
            __props__.__dict__["create"] = create
            __props__.__dict__["delete"] = delete
            if name is None and not opts.urn:
                raise TypeError("Missing required property 'name'")
            __props__.__dict__["name"] = name
            __props__.__dict__["net_adapter_name"] = net_adapter_name
            __props__.__dict__["notes"] = notes
            if switch_type is None and not opts.urn:
                raise TypeError("Missing required property 'switch_type'")
            __props__.__dict__["switch_type"] = switch_type
            __props__.__dict__["triggers"] = triggers
            __props__.__dict__["update"] = update
        replace_on_changes = pulumi.ResourceOptions(replace_on_changes=["triggers[*]"])
        opts = pulumi.ResourceOptions.merge(opts, replace_on_changes)
        super(VirtualSwitch, __self__).__init__(
            'hyperv:virtualswitch:VirtualSwitch',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'VirtualSwitch':
        """
        Get an existing VirtualSwitch resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = VirtualSwitchArgs.__new__(VirtualSwitchArgs)

        __props__.__dict__["allow_management_os"] = None
        __props__.__dict__["create"] = None
        __props__.__dict__["delete"] = None
        __props__.__dict__["name"] = None
        __props__.__dict__["net_adapter_name"] = None
        __props__.__dict__["notes"] = None
        __props__.__dict__["switch_type"] = None
        __props__.__dict__["triggers"] = None
        __props__.__dict__["update"] = None
        return VirtualSwitch(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter(name="allowManagementOs")
    def allow_management_os(self) -> pulumi.Output[Optional[builtins.bool]]:
        """
        Allow the management OS to access the switch (External switches)
        """
        return pulumi.get(self, "allow_management_os")

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
    @pulumi.getter
    def name(self) -> pulumi.Output[builtins.str]:
        """
        Name of the virtual switch
        """
        return pulumi.get(self, "name")

    @property
    @pulumi.getter(name="netAdapterName")
    def net_adapter_name(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        Name of the physical network adapter to bind to (External switches)
        """
        return pulumi.get(self, "net_adapter_name")

    @property
    @pulumi.getter
    def notes(self) -> pulumi.Output[Optional[builtins.str]]:
        """
        Notes or description for the virtual switch
        """
        return pulumi.get(self, "notes")

    @property
    @pulumi.getter(name="switchType")
    def switch_type(self) -> pulumi.Output[builtins.str]:
        """
        Type of switch: 'External', 'Internal', or 'Private'
        """
        return pulumi.get(self, "switch_type")

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

