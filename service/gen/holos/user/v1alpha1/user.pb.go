// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0-devel
// 	protoc        (unknown)
// source: holos/user/v1alpha1/user.proto

package user

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	v1alpha1 "github.com/holos-run/holos/service/gen/holos/object/v1alpha1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// User represents a human user of the system.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique uuid assigned by the server.
	Id *string `protobuf:"bytes,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	// Subject represents the oidc iss and sub claims of the user.
	Subject *v1alpha1.Subject `protobuf:"bytes,2,opt,name=subject,proto3,oneof" json:"subject,omitempty"`
	// Email address of the user.
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	// True if the user email has been verified.
	EmailVerified *bool `protobuf:"varint,4,opt,name=email_verified,json=emailVerified,proto3,oneof" json:"email_verified,omitempty"`
	// Full name provided by the user.
	Name *string `protobuf:"bytes,5,opt,name=name,proto3,oneof" json:"name,omitempty"`
	// Given or first name of the user.
	GivenName *string `protobuf:"bytes,6,opt,name=given_name,json=givenName,proto3,oneof" json:"given_name,omitempty"`
	// Family or last name of the user.
	FamilyName *string `protobuf:"bytes,7,opt,name=family_name,json=familyName,proto3,oneof" json:"family_name,omitempty"`
	// Groups the user is a member of.  This field represents the oidc groups
	// claim.
	Groups []string `protobuf:"bytes,8,rep,name=groups,proto3" json:"groups,omitempty"`
	// https url to an user avatar profile picture.  Should be at least a 200x200 px square image.
	Picture *string `protobuf:"bytes,9,opt,name=picture,proto3,oneof" json:"picture,omitempty"`
	// Detail applicable to all resource objects in the system such as created and
	// updated metadata.
	Detail *v1alpha1.Detail `protobuf:"bytes,10,opt,name=detail,proto3,oneof" json:"detail,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_holos_user_v1alpha1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_holos_user_v1alpha1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_holos_user_v1alpha1_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *User) GetSubject() *v1alpha1.Subject {
	if x != nil {
		return x.Subject
	}
	return nil
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetEmailVerified() bool {
	if x != nil && x.EmailVerified != nil {
		return *x.EmailVerified
	}
	return false
}

func (x *User) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *User) GetGivenName() string {
	if x != nil && x.GivenName != nil {
		return *x.GivenName
	}
	return ""
}

func (x *User) GetFamilyName() string {
	if x != nil && x.FamilyName != nil {
		return *x.FamilyName
	}
	return ""
}

func (x *User) GetGroups() []string {
	if x != nil {
		return x.Groups
	}
	return nil
}

func (x *User) GetPicture() string {
	if x != nil && x.Picture != nil {
		return *x.Picture
	}
	return ""
}

func (x *User) GetDetail() *v1alpha1.Detail {
	if x != nil {
		return x.Detail
	}
	return nil
}

var File_holos_user_v1alpha1_user_proto protoreflect.FileDescriptor

var file_holos_user_v1alpha1_user_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x13, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x22, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x08, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x1d, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x3d,
	0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x48,
	0x01, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48,
	0x04, 0x72, 0x02, 0x60, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x2a, 0x0a, 0x0e,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x0d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0xde, 0x01, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0xc4, 0x01, 0xba, 0x48, 0xc0, 0x01, 0xba, 0x01,
	0x5a, 0x0a, 0x1a, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e, 0x6f, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1d, 0x43,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x20, 0x73, 0x74, 0x61, 0x72, 0x74, 0x20, 0x77, 0x69, 0x74, 0x68,
	0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x1a, 0x1d, 0x21, 0x74,
	0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5e, 0x5b, 0x5b,
	0x3a, 0x73, 0x70, 0x61, 0x63, 0x65, 0x3a, 0x5d, 0x5d, 0x27, 0x29, 0xba, 0x01, 0x59, 0x0a, 0x1b,
	0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e, 0x6f, 0x5f, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67,
	0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1b, 0x43, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x20, 0x65, 0x6e, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x77, 0x68, 0x69,
	0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x1a, 0x1d, 0x21, 0x74, 0x68, 0x69, 0x73, 0x2e,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5b, 0x5b, 0x3a, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x3a, 0x5d, 0x5d, 0x24, 0x27, 0x29, 0x72, 0x05, 0x10, 0x01, 0x18, 0xff, 0x01, 0x48, 0x03,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0xe9, 0x01, 0x0a, 0x0a, 0x67, 0x69,
	0x76, 0x65, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0xc4,
	0x01, 0xba, 0x48, 0xc0, 0x01, 0xba, 0x01, 0x5a, 0x0a, 0x1a, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e,
	0x6f, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x1d, 0x43, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x20, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x2e, 0x1a, 0x1d, 0x21, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x73, 0x28, 0x27, 0x5e, 0x5b, 0x5b, 0x3a, 0x73, 0x70, 0x61, 0x63, 0x65, 0x3a, 0x5d, 0x5d,
	0x27, 0x29, 0xba, 0x01, 0x59, 0x0a, 0x1b, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e, 0x6f, 0x5f, 0x74,
	0x72, 0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x1b, 0x43, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x6e, 0x64, 0x20, 0x77,
	0x69, 0x74, 0x68, 0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x1a,
	0x1d, 0x21, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27,
	0x5b, 0x5b, 0x3a, 0x73, 0x70, 0x61, 0x63, 0x65, 0x3a, 0x5d, 0x5d, 0x24, 0x27, 0x29, 0x72, 0x05,
	0x10, 0x01, 0x18, 0xff, 0x01, 0x48, 0x04, 0x52, 0x09, 0x67, 0x69, 0x76, 0x65, 0x6e, 0x4e, 0x61,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0xeb, 0x01, 0x0a, 0x0b, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0xc4, 0x01, 0xba, 0x48,
	0xc0, 0x01, 0xba, 0x01, 0x5a, 0x0a, 0x1a, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e, 0x6f, 0x5f, 0x6c,
	0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x12, 0x1d, 0x43, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x20, 0x73, 0x74, 0x61, 0x72, 0x74, 0x20,
	0x77, 0x69, 0x74, 0x68, 0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e,
	0x1a, 0x1d, 0x21, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28,
	0x27, 0x5e, 0x5b, 0x5b, 0x3a, 0x73, 0x70, 0x61, 0x63, 0x65, 0x3a, 0x5d, 0x5d, 0x27, 0x29, 0xba,
	0x01, 0x59, 0x0a, 0x1b, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x6e, 0x6f, 0x5f, 0x74, 0x72, 0x61, 0x69,
	0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x1b, 0x43, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x6e, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68,
	0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x1a, 0x1d, 0x21, 0x74,
	0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5b, 0x5b, 0x3a,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x3a, 0x5d, 0x5d, 0x24, 0x27, 0x29, 0x72, 0x05, 0x10, 0x01, 0x18,
	0xff, 0x01, 0x48, 0x05, 0x52, 0x0a, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x08, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x36, 0x0a, 0x07, 0x70,
	0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x48,
	0x14, 0x72, 0x12, 0x10, 0x01, 0x18, 0xff, 0x0f, 0x3a, 0x08, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a,
	0x2f, 0x2f, 0x88, 0x01, 0x01, 0x48, 0x06, 0x52, 0x07, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x3a, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2e, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x48, 0x07, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x42,
	0x05, 0x0a, 0x03, 0x5f, 0x69, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x67, 0x69, 0x76, 0x65, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0e, 0x0a,
	0x0c, 0x5f, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x70, 0x69, 0x63, 0x74, 0x75, 0x72, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x6c, 0x6f, 0x73, 0x2d, 0x72, 0x75, 0x6e, 0x2f, 0x68, 0x6f, 0x6c,
	0x6f, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x68,
	0x6f, 0x6c, 0x6f, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_holos_user_v1alpha1_user_proto_rawDescOnce sync.Once
	file_holos_user_v1alpha1_user_proto_rawDescData = file_holos_user_v1alpha1_user_proto_rawDesc
)

func file_holos_user_v1alpha1_user_proto_rawDescGZIP() []byte {
	file_holos_user_v1alpha1_user_proto_rawDescOnce.Do(func() {
		file_holos_user_v1alpha1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_holos_user_v1alpha1_user_proto_rawDescData)
	})
	return file_holos_user_v1alpha1_user_proto_rawDescData
}

var file_holos_user_v1alpha1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_holos_user_v1alpha1_user_proto_goTypes = []interface{}{
	(*User)(nil),             // 0: holos.user.v1alpha1.User
	(*v1alpha1.Subject)(nil), // 1: holos.object.v1alpha1.Subject
	(*v1alpha1.Detail)(nil),  // 2: holos.object.v1alpha1.Detail
}
var file_holos_user_v1alpha1_user_proto_depIdxs = []int32{
	1, // 0: holos.user.v1alpha1.User.subject:type_name -> holos.object.v1alpha1.Subject
	2, // 1: holos.user.v1alpha1.User.detail:type_name -> holos.object.v1alpha1.Detail
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_holos_user_v1alpha1_user_proto_init() }
func file_holos_user_v1alpha1_user_proto_init() {
	if File_holos_user_v1alpha1_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_holos_user_v1alpha1_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
	file_holos_user_v1alpha1_user_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_holos_user_v1alpha1_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_holos_user_v1alpha1_user_proto_goTypes,
		DependencyIndexes: file_holos_user_v1alpha1_user_proto_depIdxs,
		MessageInfos:      file_holos_user_v1alpha1_user_proto_msgTypes,
	}.Build()
	File_holos_user_v1alpha1_user_proto = out.File
	file_holos_user_v1alpha1_user_proto_rawDesc = nil
	file_holos_user_v1alpha1_user_proto_goTypes = nil
	file_holos_user_v1alpha1_user_proto_depIdxs = nil
}
