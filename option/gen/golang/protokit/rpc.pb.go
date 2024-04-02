// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: protokit/rpc.proto

package protokit

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MethodAskType int32

const (
	MethodAskType_ASK  MethodAskType = 0
	MethodAskType_TELL MethodAskType = 1
)

// Enum value maps for MethodAskType.
var (
	MethodAskType_name = map[int32]string{
		0: "ASK",
		1: "TELL",
	}
	MethodAskType_value = map[string]int32{
		"ASK":  0,
		"TELL": 1,
	}
)

func (x MethodAskType) Enum() *MethodAskType {
	p := new(MethodAskType)
	*p = x
	return p
}

func (x MethodAskType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MethodAskType) Descriptor() protoreflect.EnumDescriptor {
	return file_protokit_rpc_proto_enumTypes[0].Descriptor()
}

func (MethodAskType) Type() protoreflect.EnumType {
	return &file_protokit_rpc_proto_enumTypes[0]
}

// NumberInt32 hacked by protokitgo
func (x MethodAskType) NumberInt32() int32 {
	return int32(x.Number())
}

func (x MethodAskType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MethodAskType.Descriptor instead.
func (MethodAskType) EnumDescriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{0}
}

type QueryPathType int32

const (
	QueryPathType_SNAKE_CASE QueryPathType = 0
	QueryPathType_CAMEL_CASE QueryPathType = 1
)

// Enum value maps for QueryPathType.
var (
	QueryPathType_name = map[int32]string{
		0: "SNAKE_CASE",
		1: "CAMEL_CASE",
	}
	QueryPathType_value = map[string]int32{
		"SNAKE_CASE": 0,
		"CAMEL_CASE": 1,
	}
)

func (x QueryPathType) Enum() *QueryPathType {
	p := new(QueryPathType)
	*p = x
	return p
}

func (x QueryPathType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueryPathType) Descriptor() protoreflect.EnumDescriptor {
	return file_protokit_rpc_proto_enumTypes[1].Descriptor()
}

func (QueryPathType) Type() protoreflect.EnumType {
	return &file_protokit_rpc_proto_enumTypes[1]
}

// NumberInt32 hacked by protokitgo
func (x QueryPathType) NumberInt32() int32 {
	return int32(x.Number())
}

func (x QueryPathType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueryPathType.Descriptor instead.
func (QueryPathType) EnumDescriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{1}
}

type RpcServiceOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rpc                bool          `protobuf:"varint,1,opt,name=rpc,proto3" json:"rpc,omitempty" db:"rpc"` // if (!rpc && !actor && !erpc) rpc = true
	Actor              bool          `protobuf:"varint,2,opt,name=actor,proto3" json:"actor,omitempty" db:"actor"`
	Erpc               bool          `protobuf:"varint,3,opt,name=erpc,proto3" json:"erpc,omitempty" db:"erpc"`
	AskTell            MethodAskType `protobuf:"varint,11,opt,name=ask_tell,json=askTell,proto3,enum=protokit.MethodAskType" json:"ask_tell,omitempty" db:"ask_tell"`
	QueryPath          string        `protobuf:"bytes,12,opt,name=query_path,json=queryPath,proto3" json:"query_path,omitempty" db:"query_path"` // https://dcs-devcenter-site-online.auto.centurygame.io/docs-tools-protokitgo/v2.6/code/service/#querypath
	LangOff            string        `protobuf:"bytes,13,opt,name=lang_off,json=langOff,proto3" json:"lang_off,omitempty" db:"lang_off"`
	QueryPathSnakeCase QueryPathType `protobuf:"varint,14,opt,name=query_path_snake_case,json=queryPathSnakeCase,proto3,enum=protokit.QueryPathType" json:"query_path_snake_case,omitempty" db:"query_path_snake_case"`
}

func (x *RpcServiceOptions) Reset() {
	*x = RpcServiceOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protokit_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcServiceOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcServiceOptions) ProtoMessage() {}

func (x *RpcServiceOptions) ProtoReflect() protoreflect.Message {
	mi := &file_protokit_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcServiceOptions.ProtoReflect.Descriptor instead.
func (*RpcServiceOptions) Descriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *RpcServiceOptions) GetRpc() bool {
	if x != nil {
		return x.Rpc
	}
	return false
}

func (x *RpcServiceOptions) GetActor() bool {
	if x != nil {
		return x.Actor
	}
	return false
}

func (x *RpcServiceOptions) GetErpc() bool {
	if x != nil {
		return x.Erpc
	}
	return false
}

func (x *RpcServiceOptions) GetAskTell() MethodAskType {
	if x != nil {
		return x.AskTell
	}
	return MethodAskType_ASK
}

func (x *RpcServiceOptions) GetQueryPath() string {
	if x != nil {
		return x.QueryPath
	}
	return ""
}

func (x *RpcServiceOptions) GetLangOff() string {
	if x != nil {
		return x.LangOff
	}
	return ""
}

func (x *RpcServiceOptions) GetQueryPathSnakeCase() QueryPathType {
	if x != nil {
		return x.QueryPathSnakeCase
	}
	return QueryPathType_SNAKE_CASE
}

type BackOfficeServiceOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enable  bool   `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty" db:"enable"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" db:"name"`
	Comment string `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty" db:"comment"`
}

func (x *BackOfficeServiceOptions) Reset() {
	*x = BackOfficeServiceOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protokit_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackOfficeServiceOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackOfficeServiceOptions) ProtoMessage() {}

func (x *BackOfficeServiceOptions) ProtoReflect() protoreflect.Message {
	mi := &file_protokit_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackOfficeServiceOptions.ProtoReflect.Descriptor instead.
func (*BackOfficeServiceOptions) Descriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *BackOfficeServiceOptions) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *BackOfficeServiceOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BackOfficeServiceOptions) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

type RpcMethodOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rpc        bool          `protobuf:"varint,1,opt,name=rpc,proto3" json:"rpc,omitempty" db:"rpc"` // if (!rpc && !actor && !erpc) rpc = true
	Actor      bool          `protobuf:"varint,2,opt,name=actor,proto3" json:"actor,omitempty" db:"actor"`
	Erpc       bool          `protobuf:"varint,3,opt,name=erpc,proto3" json:"erpc,omitempty" db:"erpc"`
	AskTell    MethodAskType `protobuf:"varint,11,opt,name=ask_tell,json=askTell,proto3,enum=protokit.MethodAskType" json:"ask_tell,omitempty" db:"ask_tell"`
	Alias      string        `protobuf:"bytes,12,opt,name=alias,proto3" json:"alias,omitempty" db:"alias"`
	ActorAlias string        `protobuf:"bytes,13,opt,name=actor_alias,json=actorAlias,proto3" json:"actor_alias,omitempty" db:"actor_alias"`
	LangOff    string        `protobuf:"bytes,14,opt,name=lang_off,json=langOff,proto3" json:"lang_off,omitempty" db:"lang_off"`
}

func (x *RpcMethodOptions) Reset() {
	*x = RpcMethodOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protokit_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcMethodOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcMethodOptions) ProtoMessage() {}

func (x *RpcMethodOptions) ProtoReflect() protoreflect.Message {
	mi := &file_protokit_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcMethodOptions.ProtoReflect.Descriptor instead.
func (*RpcMethodOptions) Descriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *RpcMethodOptions) GetRpc() bool {
	if x != nil {
		return x.Rpc
	}
	return false
}

func (x *RpcMethodOptions) GetActor() bool {
	if x != nil {
		return x.Actor
	}
	return false
}

func (x *RpcMethodOptions) GetErpc() bool {
	if x != nil {
		return x.Erpc
	}
	return false
}

func (x *RpcMethodOptions) GetAskTell() MethodAskType {
	if x != nil {
		return x.AskTell
	}
	return MethodAskType_ASK
}

func (x *RpcMethodOptions) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *RpcMethodOptions) GetActorAlias() string {
	if x != nil {
		return x.ActorAlias
	}
	return ""
}

func (x *RpcMethodOptions) GetLangOff() string {
	if x != nil {
		return x.LangOff
	}
	return ""
}

type BackOfficeMethodOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enable  bool   `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty" db:"enable"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" db:"name"`
	Comment string `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty" db:"comment"`
}

func (x *BackOfficeMethodOptions) Reset() {
	*x = BackOfficeMethodOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protokit_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackOfficeMethodOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackOfficeMethodOptions) ProtoMessage() {}

func (x *BackOfficeMethodOptions) ProtoReflect() protoreflect.Message {
	mi := &file_protokit_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackOfficeMethodOptions.ProtoReflect.Descriptor instead.
func (*BackOfficeMethodOptions) Descriptor() ([]byte, []int) {
	return file_protokit_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *BackOfficeMethodOptions) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *BackOfficeMethodOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BackOfficeMethodOptions) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

var file_protokit_rpc_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*RpcServiceOptions)(nil),
		Field:         66801,
		Name:          "protokit.rpc_service",
		Tag:           "bytes,66801,opt,name=rpc_service",
		Filename:      "protokit/rpc.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*BackOfficeServiceOptions)(nil),
		Field:         66802,
		Name:          "protokit.back_office_service",
		Tag:           "bytes,66802,opt,name=back_office_service",
		Filename:      "protokit/rpc.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*RpcMethodOptions)(nil),
		Field:         68801,
		Name:          "protokit.rpc_method",
		Tag:           "bytes,68801,opt,name=rpc_method",
		Filename:      "protokit/rpc.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*BackOfficeMethodOptions)(nil),
		Field:         68802,
		Name:          "protokit.back_office_method",
		Tag:           "bytes,68802,opt,name=back_office_method",
		Filename:      "protokit/rpc.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// optional protokit.RpcServiceOptions rpc_service = 66801;
	E_RpcService = &file_protokit_rpc_proto_extTypes[0]
	// optional protokit.BackOfficeServiceOptions back_office_service = 66802;
	E_BackOfficeService = &file_protokit_rpc_proto_extTypes[1]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional protokit.RpcMethodOptions rpc_method = 68801;
	E_RpcMethod = &file_protokit_rpc_proto_extTypes[2]
	// optional protokit.BackOfficeMethodOptions back_office_method = 68802;
	E_BackOfficeMethod = &file_protokit_rpc_proto_extTypes[3]
)

var File_protokit_rpc_proto protoreflect.FileDescriptor

var file_protokit_rpc_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x89, 0x02, 0x0a, 0x11, 0x52, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x70, 0x63, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x03, 0x72, 0x70, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x65, 0x72, 0x70, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x65, 0x72,
	0x70, 0x63, 0x12, 0x32, 0x0a, 0x08, 0x61, 0x73, 0x6b, 0x5f, 0x74, 0x65, 0x6c, 0x6c, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x41, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x61,
	0x73, 0x6b, 0x54, 0x65, 0x6c, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f,
	0x70, 0x61, 0x74, 0x68, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x50, 0x61, 0x74, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x5f, 0x6f, 0x66,
	0x66, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x61, 0x6e, 0x67, 0x4f, 0x66, 0x66,
	0x12, 0x4a, 0x0a, 0x15, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x73,
	0x6e, 0x61, 0x6b, 0x65, 0x5f, 0x63, 0x61, 0x73, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x50, 0x61, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x52, 0x12, 0x71, 0x75, 0x65, 0x72, 0x79, 0x50,
	0x61, 0x74, 0x68, 0x53, 0x6e, 0x61, 0x6b, 0x65, 0x43, 0x61, 0x73, 0x65, 0x22, 0x60, 0x0a, 0x18,
	0x42, 0x61, 0x63, 0x6b, 0x4f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0xd4,
	0x01, 0x0a, 0x10, 0x52, 0x70, 0x63, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x70, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x03, 0x72, 0x70, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x65,
	0x72, 0x70, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x65, 0x72, 0x70, 0x63, 0x12,
	0x32, 0x0a, 0x08, 0x61, 0x73, 0x6b, 0x5f, 0x74, 0x65, 0x6c, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x41, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x61, 0x73, 0x6b, 0x54,
	0x65, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x61,
	0x6e, 0x67, 0x5f, 0x6f, 0x66, 0x66, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x61,
	0x6e, 0x67, 0x4f, 0x66, 0x66, 0x22, 0x5f, 0x0a, 0x17, 0x42, 0x61, 0x63, 0x6b, 0x4f, 0x66, 0x66,
	0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2a, 0x22, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x41, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x53, 0x4b, 0x10, 0x00,
	0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x4c, 0x4c, 0x10, 0x01, 0x2a, 0x2f, 0x0a, 0x0d, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x50, 0x61, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x53,
	0x4e, 0x41, 0x4b, 0x45, 0x5f, 0x43, 0x41, 0x53, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43,
	0x41, 0x4d, 0x45, 0x4c, 0x5f, 0x43, 0x41, 0x53, 0x45, 0x10, 0x01, 0x3a, 0x5f, 0x0a, 0x0b, 0x72,
	0x70, 0x63, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf1, 0x89, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x52,
	0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x0a, 0x72, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3a, 0x75, 0x0a, 0x13,
	0x62, 0x61, 0x63, 0x6b, 0x5f, 0x6f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf2, 0x89, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x4f, 0x66, 0x66, 0x69,
	0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x11, 0x62, 0x61, 0x63, 0x6b, 0x4f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x3a, 0x5b, 0x0a, 0x0a, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xc1, 0x99, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x52, 0x70, 0x63, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x09, 0x72, 0x70, 0x63, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x3a, 0x71, 0x0a, 0x12, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x6f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x5f,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc2, 0x99, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x4f, 0x66,
	0x66, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x10, 0x62, 0x61, 0x63, 0x6b, 0x4f, 0x66, 0x66, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x42, 0x4b, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x61, 0x6e, 0x64, 0x77, 0x69, 0x63, 0x68, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69,
	0x74, 0xaa, 0x02, 0x0c, 0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6b, 0x69, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protokit_rpc_proto_rawDescOnce sync.Once
	file_protokit_rpc_proto_rawDescData = file_protokit_rpc_proto_rawDesc
)

func file_protokit_rpc_proto_rawDescGZIP() []byte {
	file_protokit_rpc_proto_rawDescOnce.Do(func() {
		file_protokit_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_protokit_rpc_proto_rawDescData)
	})
	return file_protokit_rpc_proto_rawDescData
}

var file_protokit_rpc_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protokit_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protokit_rpc_proto_goTypes = []interface{}{
	(MethodAskType)(0),                  // 0: protokit.MethodAskType
	(QueryPathType)(0),                  // 1: protokit.QueryPathType
	(*RpcServiceOptions)(nil),           // 2: protokit.RpcServiceOptions
	(*BackOfficeServiceOptions)(nil),    // 3: protokit.BackOfficeServiceOptions
	(*RpcMethodOptions)(nil),            // 4: protokit.RpcMethodOptions
	(*BackOfficeMethodOptions)(nil),     // 5: protokit.BackOfficeMethodOptions
	(*descriptorpb.ServiceOptions)(nil), // 6: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 7: google.protobuf.MethodOptions
}
var file_protokit_rpc_proto_depIdxs = []int32{
	0,  // 0: protokit.RpcServiceOptions.ask_tell:type_name -> protokit.MethodAskType
	1,  // 1: protokit.RpcServiceOptions.query_path_snake_case:type_name -> protokit.QueryPathType
	0,  // 2: protokit.RpcMethodOptions.ask_tell:type_name -> protokit.MethodAskType
	6,  // 3: protokit.rpc_service:extendee -> google.protobuf.ServiceOptions
	6,  // 4: protokit.back_office_service:extendee -> google.protobuf.ServiceOptions
	7,  // 5: protokit.rpc_method:extendee -> google.protobuf.MethodOptions
	7,  // 6: protokit.back_office_method:extendee -> google.protobuf.MethodOptions
	2,  // 7: protokit.rpc_service:type_name -> protokit.RpcServiceOptions
	3,  // 8: protokit.back_office_service:type_name -> protokit.BackOfficeServiceOptions
	4,  // 9: protokit.rpc_method:type_name -> protokit.RpcMethodOptions
	5,  // 10: protokit.back_office_method:type_name -> protokit.BackOfficeMethodOptions
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	7,  // [7:11] is the sub-list for extension type_name
	3,  // [3:7] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_protokit_rpc_proto_init() }
func file_protokit_rpc_proto_init() {
	if File_protokit_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protokit_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcServiceOptions); i {
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
		file_protokit_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackOfficeServiceOptions); i {
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
		file_protokit_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcMethodOptions); i {
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
		file_protokit_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackOfficeMethodOptions); i {
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
			RawDescriptor: file_protokit_rpc_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_protokit_rpc_proto_goTypes,
		DependencyIndexes: file_protokit_rpc_proto_depIdxs,
		EnumInfos:         file_protokit_rpc_proto_enumTypes,
		MessageInfos:      file_protokit_rpc_proto_msgTypes,
		ExtensionInfos:    file_protokit_rpc_proto_extTypes,
	}.Build()
	File_protokit_rpc_proto = out.File
	file_protokit_rpc_proto_rawDesc = nil
	file_protokit_rpc_proto_goTypes = nil
	file_protokit_rpc_proto_depIdxs = nil
}
