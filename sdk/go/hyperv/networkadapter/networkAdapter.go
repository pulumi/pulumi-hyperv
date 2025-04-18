// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package networkadapter

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-hyperv/provider/go/hyperv/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// # Network Adapter Resource
//
// The Network Adapter resource allows you to create and manage network adapters for virtual machines in Hyper-V.
//
// ## Example Usage
//
// ### Standalone Network Adapter
//
// ### Using the NetworkAdapters Property in Machine Resource
//
// You can also define network adapters directly in the Machine resource using the `networkAdapters` property:
//
// ## Input Properties
//
// | Property         | Type     | Required | Description |
// |------------------|----------|----------|-------------|
// | name             | string   | Yes      | Name of the network adapter |
// | vmName           | string   | Yes      | Name of the virtual machine to attach the network adapter to |
// | switchName       | string   | Yes      | Name of the virtual switch to connect the network adapter to |
// | macAddress       | string   | No       | MAC address for the network adapter. If not specified, a dynamic MAC address will be generated |
// | vlanId           | number   | No       | VLAN ID for the network adapter. If not specified, no VLAN tagging is used |
// | dhcpGuard        | boolean  | No       | Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages |
// | routerGuard      | boolean  | No       | Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages |
// | portMirroring    | string   | No       | Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None |
// | ieeePriorityTag  | boolean  | No       | Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value |
// | vmqWeight        | number   | No       | VMQ weight for the network adapter. A value of 0 disables VMQ |
// | ipAddresses      | string   | No       | Comma-separated list of IP addresses to assign to the network adapter |
//
// ## Output Properties
//
// | Property         | Type     | Description |
// |------------------|----------|-------------|
// | adapterId        | string   | The ID of the network adapter |
//
// ## Lifecycle Management
//
// - **Create**: Creates a new network adapter and attaches it to the specified virtual machine.
// - **Read**: Reads the properties of an existing network adapter.
// - **Update**: Updates the properties of an existing network adapter.
// - **Delete**: Removes a network adapter from a virtual machine.
//
// ## Notes
//
// - The network adapter creation will fail if the virtual machine or virtual switch does not exist.
// - Dynamic MAC addresses are automatically generated if not specified.
// - IP addresses are specified as a comma-separated string (e.g., "192.168.1.10,192.168.1.11").
// - When updating a network adapter, the virtual machine may need to be powered off depending on the properties being changed.
type NetworkAdapter struct {
	pulumi.CustomResourceState

	// The ID of the network adapter
	AdapterId pulumi.StringOutput `pulumi:"adapterId"`
	// The command to run on create.
	Create pulumi.StringPtrOutput `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrOutput `pulumi:"delete"`
	// Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
	DhcpGuard pulumi.BoolPtrOutput `pulumi:"dhcpGuard"`
	// Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
	IeeePriorityTag pulumi.BoolPtrOutput `pulumi:"ieeePriorityTag"`
	// Comma-separated list of IP addresses to assign to the network adapter.
	IpAddresses pulumi.StringPtrOutput `pulumi:"ipAddresses"`
	// MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
	MacAddress pulumi.StringPtrOutput `pulumi:"macAddress"`
	// Name of the network adapter
	Name pulumi.StringOutput `pulumi:"name"`
	// Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
	PortMirroring pulumi.StringPtrOutput `pulumi:"portMirroring"`
	// Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
	RouterGuard pulumi.BoolPtrOutput `pulumi:"routerGuard"`
	// Name of the virtual switch to connect the network adapter to
	SwitchName pulumi.StringOutput `pulumi:"switchName"`
	// Trigger a resource replacement on changes to any of these values. The
	// trigger values can be of any type. If a value is different in the current update compared to the
	// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
	// Please see the resource documentation for examples.
	Triggers pulumi.ArrayOutput `pulumi:"triggers"`
	// The command to run on update, if empty, create will
	// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
	// are set to the stdout and stderr properties of the Command resource from previous
	// create or update steps.
	Update pulumi.StringPtrOutput `pulumi:"update"`
	// VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
	VlanId pulumi.IntPtrOutput `pulumi:"vlanId"`
	// Name of the virtual machine to attach the network adapter to
	VmName pulumi.StringPtrOutput `pulumi:"vmName"`
	// VMQ weight for the network adapter. A value of 0 disables VMQ.
	VmqWeight pulumi.IntPtrOutput `pulumi:"vmqWeight"`
}

// NewNetworkAdapter registers a new resource with the given unique name, arguments, and options.
func NewNetworkAdapter(ctx *pulumi.Context,
	name string, args *NetworkAdapterArgs, opts ...pulumi.ResourceOption) (*NetworkAdapter, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Name == nil {
		return nil, errors.New("invalid value for required argument 'Name'")
	}
	if args.SwitchName == nil {
		return nil, errors.New("invalid value for required argument 'SwitchName'")
	}
	replaceOnChanges := pulumi.ReplaceOnChanges([]string{
		"triggers[*]",
	})
	opts = append(opts, replaceOnChanges)
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource NetworkAdapter
	err := ctx.RegisterResource("hyperv:networkadapter:NetworkAdapter", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetNetworkAdapter gets an existing NetworkAdapter resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetNetworkAdapter(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *NetworkAdapterState, opts ...pulumi.ResourceOption) (*NetworkAdapter, error) {
	var resource NetworkAdapter
	err := ctx.ReadResource("hyperv:networkadapter:NetworkAdapter", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering NetworkAdapter resources.
type networkAdapterState struct {
}

type NetworkAdapterState struct {
}

func (NetworkAdapterState) ElementType() reflect.Type {
	return reflect.TypeOf((*networkAdapterState)(nil)).Elem()
}

type networkAdapterArgs struct {
	// The command to run on create.
	Create *string `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete *string `pulumi:"delete"`
	// Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
	DhcpGuard *bool `pulumi:"dhcpGuard"`
	// Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
	IeeePriorityTag *bool `pulumi:"ieeePriorityTag"`
	// Comma-separated list of IP addresses to assign to the network adapter.
	IpAddresses *string `pulumi:"ipAddresses"`
	// MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
	MacAddress *string `pulumi:"macAddress"`
	// Name of the network adapter
	Name string `pulumi:"name"`
	// Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
	PortMirroring *string `pulumi:"portMirroring"`
	// Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
	RouterGuard *bool `pulumi:"routerGuard"`
	// Name of the virtual switch to connect the network adapter to
	SwitchName string `pulumi:"switchName"`
	// Trigger a resource replacement on changes to any of these values. The
	// trigger values can be of any type. If a value is different in the current update compared to the
	// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
	// Please see the resource documentation for examples.
	Triggers []interface{} `pulumi:"triggers"`
	// The command to run on update, if empty, create will
	// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
	// are set to the stdout and stderr properties of the Command resource from previous
	// create or update steps.
	Update *string `pulumi:"update"`
	// VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
	VlanId *int `pulumi:"vlanId"`
	// Name of the virtual machine to attach the network adapter to
	VmName *string `pulumi:"vmName"`
	// VMQ weight for the network adapter. A value of 0 disables VMQ.
	VmqWeight *int `pulumi:"vmqWeight"`
}

// The set of arguments for constructing a NetworkAdapter resource.
type NetworkAdapterArgs struct {
	// The command to run on create.
	Create pulumi.StringPtrInput
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrInput
	// Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
	DhcpGuard pulumi.BoolPtrInput
	// Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
	IeeePriorityTag pulumi.BoolPtrInput
	// Comma-separated list of IP addresses to assign to the network adapter.
	IpAddresses pulumi.StringPtrInput
	// MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
	MacAddress pulumi.StringPtrInput
	// Name of the network adapter
	Name pulumi.StringInput
	// Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
	PortMirroring pulumi.StringPtrInput
	// Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
	RouterGuard pulumi.BoolPtrInput
	// Name of the virtual switch to connect the network adapter to
	SwitchName pulumi.StringInput
	// Trigger a resource replacement on changes to any of these values. The
	// trigger values can be of any type. If a value is different in the current update compared to the
	// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
	// Please see the resource documentation for examples.
	Triggers pulumi.ArrayInput
	// The command to run on update, if empty, create will
	// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
	// are set to the stdout and stderr properties of the Command resource from previous
	// create or update steps.
	Update pulumi.StringPtrInput
	// VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
	VlanId pulumi.IntPtrInput
	// Name of the virtual machine to attach the network adapter to
	VmName pulumi.StringPtrInput
	// VMQ weight for the network adapter. A value of 0 disables VMQ.
	VmqWeight pulumi.IntPtrInput
}

func (NetworkAdapterArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*networkAdapterArgs)(nil)).Elem()
}

type NetworkAdapterInput interface {
	pulumi.Input

	ToNetworkAdapterOutput() NetworkAdapterOutput
	ToNetworkAdapterOutputWithContext(ctx context.Context) NetworkAdapterOutput
}

func (*NetworkAdapter) ElementType() reflect.Type {
	return reflect.TypeOf((**NetworkAdapter)(nil)).Elem()
}

func (i *NetworkAdapter) ToNetworkAdapterOutput() NetworkAdapterOutput {
	return i.ToNetworkAdapterOutputWithContext(context.Background())
}

func (i *NetworkAdapter) ToNetworkAdapterOutputWithContext(ctx context.Context) NetworkAdapterOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkAdapterOutput)
}

// NetworkAdapterArrayInput is an input type that accepts NetworkAdapterArray and NetworkAdapterArrayOutput values.
// You can construct a concrete instance of `NetworkAdapterArrayInput` via:
//
//	NetworkAdapterArray{ NetworkAdapterArgs{...} }
type NetworkAdapterArrayInput interface {
	pulumi.Input

	ToNetworkAdapterArrayOutput() NetworkAdapterArrayOutput
	ToNetworkAdapterArrayOutputWithContext(context.Context) NetworkAdapterArrayOutput
}

type NetworkAdapterArray []NetworkAdapterInput

func (NetworkAdapterArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*NetworkAdapter)(nil)).Elem()
}

func (i NetworkAdapterArray) ToNetworkAdapterArrayOutput() NetworkAdapterArrayOutput {
	return i.ToNetworkAdapterArrayOutputWithContext(context.Background())
}

func (i NetworkAdapterArray) ToNetworkAdapterArrayOutputWithContext(ctx context.Context) NetworkAdapterArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkAdapterArrayOutput)
}

// NetworkAdapterMapInput is an input type that accepts NetworkAdapterMap and NetworkAdapterMapOutput values.
// You can construct a concrete instance of `NetworkAdapterMapInput` via:
//
//	NetworkAdapterMap{ "key": NetworkAdapterArgs{...} }
type NetworkAdapterMapInput interface {
	pulumi.Input

	ToNetworkAdapterMapOutput() NetworkAdapterMapOutput
	ToNetworkAdapterMapOutputWithContext(context.Context) NetworkAdapterMapOutput
}

type NetworkAdapterMap map[string]NetworkAdapterInput

func (NetworkAdapterMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*NetworkAdapter)(nil)).Elem()
}

func (i NetworkAdapterMap) ToNetworkAdapterMapOutput() NetworkAdapterMapOutput {
	return i.ToNetworkAdapterMapOutputWithContext(context.Background())
}

func (i NetworkAdapterMap) ToNetworkAdapterMapOutputWithContext(ctx context.Context) NetworkAdapterMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkAdapterMapOutput)
}

type NetworkAdapterOutput struct{ *pulumi.OutputState }

func (NetworkAdapterOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**NetworkAdapter)(nil)).Elem()
}

func (o NetworkAdapterOutput) ToNetworkAdapterOutput() NetworkAdapterOutput {
	return o
}

func (o NetworkAdapterOutput) ToNetworkAdapterOutputWithContext(ctx context.Context) NetworkAdapterOutput {
	return o
}

// The ID of the network adapter
func (o NetworkAdapterOutput) AdapterId() pulumi.StringOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringOutput { return v.AdapterId }).(pulumi.StringOutput)
}

// The command to run on create.
func (o NetworkAdapterOutput) Create() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.Create }).(pulumi.StringPtrOutput)
}

// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
// Command resource from previous create or update steps.
func (o NetworkAdapterOutput) Delete() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.Delete }).(pulumi.StringPtrOutput)
}

// Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.
func (o NetworkAdapterOutput) DhcpGuard() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.BoolPtrOutput { return v.DhcpGuard }).(pulumi.BoolPtrOutput)
}

// Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.
func (o NetworkAdapterOutput) IeeePriorityTag() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.BoolPtrOutput { return v.IeeePriorityTag }).(pulumi.BoolPtrOutput)
}

// Comma-separated list of IP addresses to assign to the network adapter.
func (o NetworkAdapterOutput) IpAddresses() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.IpAddresses }).(pulumi.StringPtrOutput)
}

// MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.
func (o NetworkAdapterOutput) MacAddress() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.MacAddress }).(pulumi.StringPtrOutput)
}

// Name of the network adapter
func (o NetworkAdapterOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.
func (o NetworkAdapterOutput) PortMirroring() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.PortMirroring }).(pulumi.StringPtrOutput)
}

// Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.
func (o NetworkAdapterOutput) RouterGuard() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.BoolPtrOutput { return v.RouterGuard }).(pulumi.BoolPtrOutput)
}

// Name of the virtual switch to connect the network adapter to
func (o NetworkAdapterOutput) SwitchName() pulumi.StringOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringOutput { return v.SwitchName }).(pulumi.StringOutput)
}

// Trigger a resource replacement on changes to any of these values. The
// trigger values can be of any type. If a value is different in the current update compared to the
// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
// Please see the resource documentation for examples.
func (o NetworkAdapterOutput) Triggers() pulumi.ArrayOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.ArrayOutput { return v.Triggers }).(pulumi.ArrayOutput)
}

// The command to run on update, if empty, create will
// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
// are set to the stdout and stderr properties of the Command resource from previous
// create or update steps.
func (o NetworkAdapterOutput) Update() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.Update }).(pulumi.StringPtrOutput)
}

// VLAN ID for the network adapter. If not specified, no VLAN tagging is used.
func (o NetworkAdapterOutput) VlanId() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.IntPtrOutput { return v.VlanId }).(pulumi.IntPtrOutput)
}

// Name of the virtual machine to attach the network adapter to
func (o NetworkAdapterOutput) VmName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.StringPtrOutput { return v.VmName }).(pulumi.StringPtrOutput)
}

// VMQ weight for the network adapter. A value of 0 disables VMQ.
func (o NetworkAdapterOutput) VmqWeight() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *NetworkAdapter) pulumi.IntPtrOutput { return v.VmqWeight }).(pulumi.IntPtrOutput)
}

type NetworkAdapterArrayOutput struct{ *pulumi.OutputState }

func (NetworkAdapterArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*NetworkAdapter)(nil)).Elem()
}

func (o NetworkAdapterArrayOutput) ToNetworkAdapterArrayOutput() NetworkAdapterArrayOutput {
	return o
}

func (o NetworkAdapterArrayOutput) ToNetworkAdapterArrayOutputWithContext(ctx context.Context) NetworkAdapterArrayOutput {
	return o
}

func (o NetworkAdapterArrayOutput) Index(i pulumi.IntInput) NetworkAdapterOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *NetworkAdapter {
		return vs[0].([]*NetworkAdapter)[vs[1].(int)]
	}).(NetworkAdapterOutput)
}

type NetworkAdapterMapOutput struct{ *pulumi.OutputState }

func (NetworkAdapterMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*NetworkAdapter)(nil)).Elem()
}

func (o NetworkAdapterMapOutput) ToNetworkAdapterMapOutput() NetworkAdapterMapOutput {
	return o
}

func (o NetworkAdapterMapOutput) ToNetworkAdapterMapOutputWithContext(ctx context.Context) NetworkAdapterMapOutput {
	return o
}

func (o NetworkAdapterMapOutput) MapIndex(k pulumi.StringInput) NetworkAdapterOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *NetworkAdapter {
		return vs[0].(map[string]*NetworkAdapter)[vs[1].(string)]
	}).(NetworkAdapterOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkAdapterInput)(nil)).Elem(), &NetworkAdapter{})
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkAdapterArrayInput)(nil)).Elem(), NetworkAdapterArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkAdapterMapInput)(nil)).Elem(), NetworkAdapterMap{})
	pulumi.RegisterOutputType(NetworkAdapterOutput{})
	pulumi.RegisterOutputType(NetworkAdapterArrayOutput{})
	pulumi.RegisterOutputType(NetworkAdapterMapOutput{})
}
