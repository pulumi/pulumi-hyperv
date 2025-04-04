// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package machine

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-hyperv-provider/provider/go/hyperv/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// # Hyper-V Virtual Machine Management Service (VMMS)
//
// ## Overview
//
// The Virtual Machine Management Service (VMMS) is a core component of Hyper-V that manages virtual machine operations on a Windows Server or Windows Client system. This document provides information about the VMMS as implemented in the Pulumi Hyper-V provider.
//
// ## Features
//
// - Virtual machine lifecycle management (create, start, stop, pause, resume, delete)
// - Resource allocation and monitoring
// - Snapshot management
// - Virtual device configuration
//
// ## Implementation Details in Pulumi
//
// ### Virtual Machine Creation
//
// The `Create` method in the `vmController` is responsible for creating a virtual machine. It performs the following steps:
//
// 1. **Generate a Unique ID**: A unique ID is generated for the virtual machine.
// 2. **Default Values**:
//   - Memory size defaults to `1024 MB` if not specified.
//   - Processor count defaults to `1` if not specified.
//
// 3. **VMMS Client Initialization**: A VMMS client is created to interact with the Hyper-V host.
// 4. **Virtual Machine Settings**:
//   - The virtual machine is configured with `Hyper-V Generation 2`.
//   - Memory and processor settings are applied based on the provided or default values.
//
// 5. **Virtual Machine Creation**: The virtual machine is created using the configured settings.
//
// ### Read Method
//
// The `Read` method is a no-op in the current implementation. It does not perform any operations and always returns an empty state.
//
// ### Update Method
//
// The `Update` method:
//
// - Updates the virtual machine state if an `Update` command is provided.
// - Falls back to the `Create` command if no `Update` command is specified.
// - Does nothing if neither command is provided.
//
// ### Delete Method
//
// The `Delete` method is a no-op unless a `Delete` command is explicitly specified.
//
// ## Default Behavior
//
// - Outputs depend on all inputs by default.
// - No explicit dependency wiring is required.
//
// ## Usage in Pulumi
//
// When using the Pulumi Hyper-V provider, the VMMS is accessed indirectly through the `Vm` resource type. The resource supports the following properties:
//
// - `processorCount`: Number of processors to allocate (default: 1).
// - `memorySize`: Memory size in MB (default: 1024).
//
// ## Authentication and Security
//
// The VMMS requires appropriate permissions to manage Hyper-V objects. When using the Pulumi Hyper-V provider, ensure that:
//
// 1. The user running Pulumi commands has administrative privileges on the Hyper-V host.
// 2. Required firewall rules are configured if managing a remote Hyper-V host.
// 3. Proper credentials are provided when connecting to remote systems.
//
// ## Related Documentation
//
// - [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
// - [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/
type Machine struct {
	pulumi.CustomResourceState

	// The command to run on create.
	Create pulumi.StringPtrOutput `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrOutput `pulumi:"delete"`
	// Generation of the Virtual Machine. Defaults to 2.
	Generation pulumi.IntPtrOutput `pulumi:"generation"`
	// Name of the Virtual Machine
	MachineName pulumi.StringOutput `pulumi:"machineName"`
	// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
	MemorySize pulumi.IntPtrOutput `pulumi:"memorySize"`
	// Number of processors to allocate to the Virtual Machine. Defaults to 1.
	ProcessorCount pulumi.IntPtrOutput `pulumi:"processorCount"`
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
}

// NewMachine registers a new resource with the given unique name, arguments, and options.
func NewMachine(ctx *pulumi.Context,
	name string, args *MachineArgs, opts ...pulumi.ResourceOption) (*Machine, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.MachineName == nil {
		return nil, errors.New("invalid value for required argument 'MachineName'")
	}
	replaceOnChanges := pulumi.ReplaceOnChanges([]string{
		"triggers[*]",
	})
	opts = append(opts, replaceOnChanges)
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource Machine
	err := ctx.RegisterResource("hyperv:machine:Machine", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMachine gets an existing Machine resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMachine(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MachineState, opts ...pulumi.ResourceOption) (*Machine, error) {
	var resource Machine
	err := ctx.ReadResource("hyperv:machine:Machine", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Machine resources.
type machineState struct {
}

type MachineState struct {
}

func (MachineState) ElementType() reflect.Type {
	return reflect.TypeOf((*machineState)(nil)).Elem()
}

type machineArgs struct {
	// The command to run on create.
	Create *string `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete *string `pulumi:"delete"`
	// Generation of the Virtual Machine. Defaults to 2.
	Generation *int `pulumi:"generation"`
	// Name of the Virtual Machine
	MachineName string `pulumi:"machineName"`
	// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
	MemorySize *int `pulumi:"memorySize"`
	// Number of processors to allocate to the Virtual Machine. Defaults to 1.
	ProcessorCount *int `pulumi:"processorCount"`
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
}

// The set of arguments for constructing a Machine resource.
type MachineArgs struct {
	// The command to run on create.
	Create pulumi.StringPtrInput
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrInput
	// Generation of the Virtual Machine. Defaults to 2.
	Generation pulumi.IntPtrInput
	// Name of the Virtual Machine
	MachineName pulumi.StringInput
	// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
	MemorySize pulumi.IntPtrInput
	// Number of processors to allocate to the Virtual Machine. Defaults to 1.
	ProcessorCount pulumi.IntPtrInput
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
}

func (MachineArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*machineArgs)(nil)).Elem()
}

type MachineInput interface {
	pulumi.Input

	ToMachineOutput() MachineOutput
	ToMachineOutputWithContext(ctx context.Context) MachineOutput
}

func (*Machine) ElementType() reflect.Type {
	return reflect.TypeOf((**Machine)(nil)).Elem()
}

func (i *Machine) ToMachineOutput() MachineOutput {
	return i.ToMachineOutputWithContext(context.Background())
}

func (i *Machine) ToMachineOutputWithContext(ctx context.Context) MachineOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MachineOutput)
}

// MachineArrayInput is an input type that accepts MachineArray and MachineArrayOutput values.
// You can construct a concrete instance of `MachineArrayInput` via:
//
//	MachineArray{ MachineArgs{...} }
type MachineArrayInput interface {
	pulumi.Input

	ToMachineArrayOutput() MachineArrayOutput
	ToMachineArrayOutputWithContext(context.Context) MachineArrayOutput
}

type MachineArray []MachineInput

func (MachineArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Machine)(nil)).Elem()
}

func (i MachineArray) ToMachineArrayOutput() MachineArrayOutput {
	return i.ToMachineArrayOutputWithContext(context.Background())
}

func (i MachineArray) ToMachineArrayOutputWithContext(ctx context.Context) MachineArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MachineArrayOutput)
}

// MachineMapInput is an input type that accepts MachineMap and MachineMapOutput values.
// You can construct a concrete instance of `MachineMapInput` via:
//
//	MachineMap{ "key": MachineArgs{...} }
type MachineMapInput interface {
	pulumi.Input

	ToMachineMapOutput() MachineMapOutput
	ToMachineMapOutputWithContext(context.Context) MachineMapOutput
}

type MachineMap map[string]MachineInput

func (MachineMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Machine)(nil)).Elem()
}

func (i MachineMap) ToMachineMapOutput() MachineMapOutput {
	return i.ToMachineMapOutputWithContext(context.Background())
}

func (i MachineMap) ToMachineMapOutputWithContext(ctx context.Context) MachineMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MachineMapOutput)
}

type MachineOutput struct{ *pulumi.OutputState }

func (MachineOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Machine)(nil)).Elem()
}

func (o MachineOutput) ToMachineOutput() MachineOutput {
	return o
}

func (o MachineOutput) ToMachineOutputWithContext(ctx context.Context) MachineOutput {
	return o
}

// The command to run on create.
func (o MachineOutput) Create() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.StringPtrOutput { return v.Create }).(pulumi.StringPtrOutput)
}

// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
// Command resource from previous create or update steps.
func (o MachineOutput) Delete() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.StringPtrOutput { return v.Delete }).(pulumi.StringPtrOutput)
}

// Generation of the Virtual Machine. Defaults to 2.
func (o MachineOutput) Generation() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.IntPtrOutput { return v.Generation }).(pulumi.IntPtrOutput)
}

// Name of the Virtual Machine
func (o MachineOutput) MachineName() pulumi.StringOutput {
	return o.ApplyT(func(v *Machine) pulumi.StringOutput { return v.MachineName }).(pulumi.StringOutput)
}

// Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.
func (o MachineOutput) MemorySize() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.IntPtrOutput { return v.MemorySize }).(pulumi.IntPtrOutput)
}

// Number of processors to allocate to the Virtual Machine. Defaults to 1.
func (o MachineOutput) ProcessorCount() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.IntPtrOutput { return v.ProcessorCount }).(pulumi.IntPtrOutput)
}

// Trigger a resource replacement on changes to any of these values. The
// trigger values can be of any type. If a value is different in the current update compared to the
// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
// Please see the resource documentation for examples.
func (o MachineOutput) Triggers() pulumi.ArrayOutput {
	return o.ApplyT(func(v *Machine) pulumi.ArrayOutput { return v.Triggers }).(pulumi.ArrayOutput)
}

// The command to run on update, if empty, create will
// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
// are set to the stdout and stderr properties of the Command resource from previous
// create or update steps.
func (o MachineOutput) Update() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Machine) pulumi.StringPtrOutput { return v.Update }).(pulumi.StringPtrOutput)
}

type MachineArrayOutput struct{ *pulumi.OutputState }

func (MachineArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Machine)(nil)).Elem()
}

func (o MachineArrayOutput) ToMachineArrayOutput() MachineArrayOutput {
	return o
}

func (o MachineArrayOutput) ToMachineArrayOutputWithContext(ctx context.Context) MachineArrayOutput {
	return o
}

func (o MachineArrayOutput) Index(i pulumi.IntInput) MachineOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *Machine {
		return vs[0].([]*Machine)[vs[1].(int)]
	}).(MachineOutput)
}

type MachineMapOutput struct{ *pulumi.OutputState }

func (MachineMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Machine)(nil)).Elem()
}

func (o MachineMapOutput) ToMachineMapOutput() MachineMapOutput {
	return o
}

func (o MachineMapOutput) ToMachineMapOutputWithContext(ctx context.Context) MachineMapOutput {
	return o
}

func (o MachineMapOutput) MapIndex(k pulumi.StringInput) MachineOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *Machine {
		return vs[0].(map[string]*Machine)[vs[1].(string)]
	}).(MachineOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MachineInput)(nil)).Elem(), &Machine{})
	pulumi.RegisterInputType(reflect.TypeOf((*MachineArrayInput)(nil)).Elem(), MachineArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*MachineMapInput)(nil)).Elem(), MachineMap{})
	pulumi.RegisterOutputType(MachineOutput{})
	pulumi.RegisterOutputType(MachineArrayOutput{})
	pulumi.RegisterOutputType(MachineMapOutput{})
}
