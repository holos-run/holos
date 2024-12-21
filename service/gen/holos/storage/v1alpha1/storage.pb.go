// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: holos/storage/v1alpha1/storage.proto

package storage

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Model represents user-defined and user-supplied form field values stored as a
// marshaled JSON object.  The model is a Struct to ensure any valid JSON object
// defined by the user via the form can be represented and stored.
type Model struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model *structpb.Struct `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
}

func (x *Model) Reset() {
	*x = Model{}
	if protoimpl.UnsafeEnabled {
		mi := &file_holos_storage_v1alpha1_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model) ProtoMessage() {}

func (x *Model) ProtoReflect() protoreflect.Message {
	mi := &file_holos_storage_v1alpha1_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Model.ProtoReflect.Descriptor instead.
func (*Model) Descriptor() ([]byte, []int) {
	return file_holos_storage_v1alpha1_storage_proto_rawDescGZIP(), []int{0}
}

func (x *Model) GetModel() *structpb.Struct {
	if x != nil {
		return x.Model
	}
	return nil
}

// Form represents the Formly input form stored as a marshaled JSON object.
type Form struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// fields represents FormlyFieldConfig[] encoded as an array of JSON objects
	// organized by section.
	FieldConfigs []*structpb.Struct `protobuf:"bytes,1,rep,name=field_configs,json=fieldConfigs,proto3" json:"field_configs,omitempty"`
}

func (x *Form) Reset() {
	*x = Form{}
	if protoimpl.UnsafeEnabled {
		mi := &file_holos_storage_v1alpha1_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Form) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Form) ProtoMessage() {}

func (x *Form) ProtoReflect() protoreflect.Message {
	mi := &file_holos_storage_v1alpha1_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Form.ProtoReflect.Descriptor instead.
func (*Form) Descriptor() ([]byte, []int) {
	return file_holos_storage_v1alpha1_storage_proto_rawDescGZIP(), []int{1}
}

func (x *Form) GetFieldConfigs() []*structpb.Struct {
	if x != nil {
		return x.FieldConfigs
	}
	return nil
}

var File_holos_storage_v1alpha1_storage_proto protoreflect.FileDescriptor

var file_holos_storage_v1alpha1_storage_proto_rawDesc = []byte{
	0x0a, 0x24, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x05,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x2d, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x22, 0x44, 0x0a, 0x04, 0x46, 0x6f, 0x72, 0x6d, 0x12, 0x3c, 0x0a, 0x0d,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0c, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2d, 0x72,
	0x75, 0x6e, 0x2f, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_holos_storage_v1alpha1_storage_proto_rawDescOnce sync.Once
	file_holos_storage_v1alpha1_storage_proto_rawDescData = file_holos_storage_v1alpha1_storage_proto_rawDesc
)

func file_holos_storage_v1alpha1_storage_proto_rawDescGZIP() []byte {
	file_holos_storage_v1alpha1_storage_proto_rawDescOnce.Do(func() {
		file_holos_storage_v1alpha1_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_holos_storage_v1alpha1_storage_proto_rawDescData)
	})
	return file_holos_storage_v1alpha1_storage_proto_rawDescData
}

var file_holos_storage_v1alpha1_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_holos_storage_v1alpha1_storage_proto_goTypes = []any{
	(*Model)(nil),           // 0: holos.storage.v1alpha1.Model
	(*Form)(nil),            // 1: holos.storage.v1alpha1.Form
	(*structpb.Struct)(nil), // 2: google.protobuf.Struct
}
var file_holos_storage_v1alpha1_storage_proto_depIdxs = []int32{
	2, // 0: holos.storage.v1alpha1.Model.model:type_name -> google.protobuf.Struct
	2, // 1: holos.storage.v1alpha1.Form.field_configs:type_name -> google.protobuf.Struct
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_holos_storage_v1alpha1_storage_proto_init() }
func file_holos_storage_v1alpha1_storage_proto_init() {
	if File_holos_storage_v1alpha1_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_holos_storage_v1alpha1_storage_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Model); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_holos_storage_v1alpha1_storage_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Form); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_holos_storage_v1alpha1_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_holos_storage_v1alpha1_storage_proto_goTypes,
		DependencyIndexes: file_holos_storage_v1alpha1_storage_proto_depIdxs,
		MessageInfos:      file_holos_storage_v1alpha1_storage_proto_msgTypes,
	}.Build()
	File_holos_storage_v1alpha1_storage_proto = out.File
	file_holos_storage_v1alpha1_storage_proto_rawDesc = nil
	file_holos_storage_v1alpha1_storage_proto_goTypes = nil
	file_holos_storage_v1alpha1_storage_proto_depIdxs = nil
}