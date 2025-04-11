// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package vhdfile

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-hyperv-provider/provider/go/hyperv/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// # VHD File Resource Management
//
// The `vhdfile` package provides utilities for managing VHD (Virtual Hard Disk) files for Hyper-V virtual machines.
//
// ## Overview
//
// This package enables creating, modifying, and deleting VHD and VHDX files through the Pulumi Hyper-V provider. It provides a clean abstraction for working with virtual disk files independent of virtual machines.
//
// ## Key Components
//
// ### Types
//
// - **VhdFile**: Represents a VHD or VHDX file for use with Hyper-V virtual machines.
//
// ### Resource Lifecycle Methods
//
// - **Create**: Creates a new VHD/VHDX file with specified properties.
// - **Read**: Retrieves information about an existing VHD/VHDX file.
// - **Update**: Modifies properties of an existing VHD/VHDX file (currently a no-op in the implementation).
// - **Delete**: Removes a VHD/VHDX file.
//
// ## Available Properties
//
// The VhdFile resource supports the following properties:
//
// | Property | Type | Description |
// |----------|------|-------------|
// | `path` | string | Path where the VHD file should be created |
// | `parentPath` | string | Path to parent VHD when creating differencing disks |
// | `diskType` | string | Type of disk (Fixed, Dynamic, Differencing) |
// | `sizeBytes` | number | Size of the disk in bytes (for Fixed and Dynamic disks) |
// | `blockSize` | number | Block size of the disk in bytes (recommended: 1048576 for 1MB) |
//
// ## Implementation Details
//
// The package uses PowerShell commands under the hood to interact with Hyper-V's VHD management functionality, providing a Go-based interface that integrates with the Pulumi resource model.
//
// ### Update Behavior
//
// The current implementation of the `Update` method is a no-op. Any changes to VHD properties that require modification of the underlying file structure will typically require replacing the resource rather than updating it in place.
//
// ## Usage Examples
//
// VHD files can be defined and managed through the Pulumi Hyper-V provider using the standard resource model. These virtual disks can then be attached to virtual machines or managed independently.
//
// ### Creating a Base VHD
//
// ### Creating a Differencing Disk
//
// ### Using with Machine Resource
//
// The VhdFile resource can be used in conjunction with the Machine resource by attaching the VHD files to a virtual machine using the `hardDrives` array:
type VhdFile struct {
	pulumi.CustomResourceState

	// Block size of the VHD file in bytes. Recommended value is 1MB (1048576 bytes) for better compatibility.
	BlockSize pulumi.IntPtrOutput `pulumi:"blockSize"`
	// The command to run on create.
	Create pulumi.StringPtrOutput `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrOutput `pulumi:"delete"`
	// Type of the VHD file (Fixed, Dynamic, or Differencing)
	DiskType pulumi.StringPtrOutput `pulumi:"diskType"`
	// Path to the parent VHD file when creating a differencing disk
	ParentPath pulumi.StringPtrOutput `pulumi:"parentPath"`
	// Path to the VHD file
	Path pulumi.StringOutput `pulumi:"path"`
	// Size of the VHD file in bytes
	SizeBytes pulumi.IntPtrOutput `pulumi:"sizeBytes"`
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

// NewVhdFile registers a new resource with the given unique name, arguments, and options.
func NewVhdFile(ctx *pulumi.Context,
	name string, args *VhdFileArgs, opts ...pulumi.ResourceOption) (*VhdFile, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Path == nil {
		return nil, errors.New("invalid value for required argument 'Path'")
	}
	replaceOnChanges := pulumi.ReplaceOnChanges([]string{
		"triggers[*]",
	})
	opts = append(opts, replaceOnChanges)
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource VhdFile
	err := ctx.RegisterResource("hyperv:vhdfile:VhdFile", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetVhdFile gets an existing VhdFile resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetVhdFile(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *VhdFileState, opts ...pulumi.ResourceOption) (*VhdFile, error) {
	var resource VhdFile
	err := ctx.ReadResource("hyperv:vhdfile:VhdFile", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering VhdFile resources.
type vhdFileState struct {
}

type VhdFileState struct {
}

func (VhdFileState) ElementType() reflect.Type {
	return reflect.TypeOf((*vhdFileState)(nil)).Elem()
}

type vhdFileArgs struct {
	// Block size of the VHD file in bytes. Recommended value is 1MB (1048576 bytes) for better compatibility.
	BlockSize *int `pulumi:"blockSize"`
	// The command to run on create.
	Create *string `pulumi:"create"`
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete *string `pulumi:"delete"`
	// Type of the VHD file (Fixed, Dynamic, or Differencing)
	DiskType *string `pulumi:"diskType"`
	// Path to the parent VHD file when creating a differencing disk
	ParentPath *string `pulumi:"parentPath"`
	// Path to the VHD file
	Path string `pulumi:"path"`
	// Size of the VHD file in bytes
	SizeBytes *int `pulumi:"sizeBytes"`
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

// The set of arguments for constructing a VhdFile resource.
type VhdFileArgs struct {
	// Block size of the VHD file in bytes. Recommended value is 1MB (1048576 bytes) for better compatibility.
	BlockSize pulumi.IntPtrInput
	// The command to run on create.
	Create pulumi.StringPtrInput
	// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
	// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
	// Command resource from previous create or update steps.
	Delete pulumi.StringPtrInput
	// Type of the VHD file (Fixed, Dynamic, or Differencing)
	DiskType pulumi.StringPtrInput
	// Path to the parent VHD file when creating a differencing disk
	ParentPath pulumi.StringPtrInput
	// Path to the VHD file
	Path pulumi.StringInput
	// Size of the VHD file in bytes
	SizeBytes pulumi.IntPtrInput
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

func (VhdFileArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*vhdFileArgs)(nil)).Elem()
}

type VhdFileInput interface {
	pulumi.Input

	ToVhdFileOutput() VhdFileOutput
	ToVhdFileOutputWithContext(ctx context.Context) VhdFileOutput
}

func (*VhdFile) ElementType() reflect.Type {
	return reflect.TypeOf((**VhdFile)(nil)).Elem()
}

func (i *VhdFile) ToVhdFileOutput() VhdFileOutput {
	return i.ToVhdFileOutputWithContext(context.Background())
}

func (i *VhdFile) ToVhdFileOutputWithContext(ctx context.Context) VhdFileOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VhdFileOutput)
}

// VhdFileArrayInput is an input type that accepts VhdFileArray and VhdFileArrayOutput values.
// You can construct a concrete instance of `VhdFileArrayInput` via:
//
//	VhdFileArray{ VhdFileArgs{...} }
type VhdFileArrayInput interface {
	pulumi.Input

	ToVhdFileArrayOutput() VhdFileArrayOutput
	ToVhdFileArrayOutputWithContext(context.Context) VhdFileArrayOutput
}

type VhdFileArray []VhdFileInput

func (VhdFileArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VhdFile)(nil)).Elem()
}

func (i VhdFileArray) ToVhdFileArrayOutput() VhdFileArrayOutput {
	return i.ToVhdFileArrayOutputWithContext(context.Background())
}

func (i VhdFileArray) ToVhdFileArrayOutputWithContext(ctx context.Context) VhdFileArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VhdFileArrayOutput)
}

// VhdFileMapInput is an input type that accepts VhdFileMap and VhdFileMapOutput values.
// You can construct a concrete instance of `VhdFileMapInput` via:
//
//	VhdFileMap{ "key": VhdFileArgs{...} }
type VhdFileMapInput interface {
	pulumi.Input

	ToVhdFileMapOutput() VhdFileMapOutput
	ToVhdFileMapOutputWithContext(context.Context) VhdFileMapOutput
}

type VhdFileMap map[string]VhdFileInput

func (VhdFileMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VhdFile)(nil)).Elem()
}

func (i VhdFileMap) ToVhdFileMapOutput() VhdFileMapOutput {
	return i.ToVhdFileMapOutputWithContext(context.Background())
}

func (i VhdFileMap) ToVhdFileMapOutputWithContext(ctx context.Context) VhdFileMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VhdFileMapOutput)
}

type VhdFileOutput struct{ *pulumi.OutputState }

func (VhdFileOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**VhdFile)(nil)).Elem()
}

func (o VhdFileOutput) ToVhdFileOutput() VhdFileOutput {
	return o
}

func (o VhdFileOutput) ToVhdFileOutputWithContext(ctx context.Context) VhdFileOutput {
	return o
}

// Block size of the VHD file in bytes. Recommended value is 1MB (1048576 bytes) for better compatibility.
func (o VhdFileOutput) BlockSize() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.IntPtrOutput { return v.BlockSize }).(pulumi.IntPtrOutput)
}

// The command to run on create.
func (o VhdFileOutput) Create() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringPtrOutput { return v.Create }).(pulumi.StringPtrOutput)
}

// The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
// and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
// Command resource from previous create or update steps.
func (o VhdFileOutput) Delete() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringPtrOutput { return v.Delete }).(pulumi.StringPtrOutput)
}

// Type of the VHD file (Fixed, Dynamic, or Differencing)
func (o VhdFileOutput) DiskType() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringPtrOutput { return v.DiskType }).(pulumi.StringPtrOutput)
}

// Path to the parent VHD file when creating a differencing disk
func (o VhdFileOutput) ParentPath() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringPtrOutput { return v.ParentPath }).(pulumi.StringPtrOutput)
}

// Path to the VHD file
func (o VhdFileOutput) Path() pulumi.StringOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringOutput { return v.Path }).(pulumi.StringOutput)
}

// Size of the VHD file in bytes
func (o VhdFileOutput) SizeBytes() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.IntPtrOutput { return v.SizeBytes }).(pulumi.IntPtrOutput)
}

// Trigger a resource replacement on changes to any of these values. The
// trigger values can be of any type. If a value is different in the current update compared to the
// previous update, the resource will be replaced, i.e., the "create" command will be re-run.
// Please see the resource documentation for examples.
func (o VhdFileOutput) Triggers() pulumi.ArrayOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.ArrayOutput { return v.Triggers }).(pulumi.ArrayOutput)
}

// The command to run on update, if empty, create will
// run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR
// are set to the stdout and stderr properties of the Command resource from previous
// create or update steps.
func (o VhdFileOutput) Update() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VhdFile) pulumi.StringPtrOutput { return v.Update }).(pulumi.StringPtrOutput)
}

type VhdFileArrayOutput struct{ *pulumi.OutputState }

func (VhdFileArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VhdFile)(nil)).Elem()
}

func (o VhdFileArrayOutput) ToVhdFileArrayOutput() VhdFileArrayOutput {
	return o
}

func (o VhdFileArrayOutput) ToVhdFileArrayOutputWithContext(ctx context.Context) VhdFileArrayOutput {
	return o
}

func (o VhdFileArrayOutput) Index(i pulumi.IntInput) VhdFileOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *VhdFile {
		return vs[0].([]*VhdFile)[vs[1].(int)]
	}).(VhdFileOutput)
}

type VhdFileMapOutput struct{ *pulumi.OutputState }

func (VhdFileMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VhdFile)(nil)).Elem()
}

func (o VhdFileMapOutput) ToVhdFileMapOutput() VhdFileMapOutput {
	return o
}

func (o VhdFileMapOutput) ToVhdFileMapOutputWithContext(ctx context.Context) VhdFileMapOutput {
	return o
}

func (o VhdFileMapOutput) MapIndex(k pulumi.StringInput) VhdFileOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *VhdFile {
		return vs[0].(map[string]*VhdFile)[vs[1].(string)]
	}).(VhdFileOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*VhdFileInput)(nil)).Elem(), &VhdFile{})
	pulumi.RegisterInputType(reflect.TypeOf((*VhdFileArrayInput)(nil)).Elem(), VhdFileArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*VhdFileMapInput)(nil)).Elem(), VhdFileMap{})
	pulumi.RegisterOutputType(VhdFileOutput{})
	pulumi.RegisterOutputType(VhdFileArrayOutput{})
	pulumi.RegisterOutputType(VhdFileMapOutput{})
}
