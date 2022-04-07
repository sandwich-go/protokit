package protokit

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Label = string

const (
	LabelOptional Label = "optional"
	LabelSlice    Label = "repeated"
	LabelMap      Label = "map"
)

type ProtoFieldTypeName = string

const (
	ProtoFieldTypeNameDouble   ProtoFieldTypeName = "double"
	ProtoFieldTypeNameFloat                       = "float"
	ProtoFieldTypeNameInt64                       = "int64"
	ProtoFieldTypeNameUInt64                      = "uint64"
	ProtoFieldTypeNameInt32                       = "int32"
	ProtoFieldTypeNameFixed64                     = "fixed64"
	ProtoFieldTypeNameFixed32                     = "fixed32"
	ProtoFieldTypeNameBool                        = "bool"
	ProtoFieldTypeNameString                      = "string"
	ProtoFieldTypeNameBytes                       = "bytes"
	ProtoFieldTypeNameUInt32                      = "uint32"
	ProtoFieldTypeNameSFixed32                    = "sfixed32"
	ProtoFieldTypeNameSFixed64                    = "sfixed64"
	ProtoFieldTypeNameSInt32                      = "sint32"
	ProtoFieldTypeNameSInt64                      = "sint64"
)

type GolangFieldTypeName = string

const (
	GolangFieldTypeNameFloat32 GolangFieldTypeName = "float32"
	GolangFieldTypeNameFloat64                     = "float64"
	GolangFieldTypeNameInt32                       = "int32"
	GolangFieldTypeNameInt64                       = "int64"
	GolangFieldTypeNameUInt32                      = "uint32"
	GolangFieldTypeNameUInt64                      = "uint64"
	GolangFieldTypeNameBool                        = "bool"
	GolangFieldTypeNameString                      = "string"
	GolangFieldTypeNameBytes                       = "[]byte"
)

var protoFieldTypeNameToGolangFieldTypeNameMapping = map[ProtoFieldTypeName]GolangFieldTypeName{
	ProtoFieldTypeNameDouble:   GolangFieldTypeNameFloat64,
	ProtoFieldTypeNameFloat:    GolangFieldTypeNameFloat32,
	ProtoFieldTypeNameInt64:    GolangFieldTypeNameInt64,
	ProtoFieldTypeNameUInt64:   GolangFieldTypeNameUInt64,
	ProtoFieldTypeNameInt32:    GolangFieldTypeNameInt32,
	ProtoFieldTypeNameFixed64:  GolangFieldTypeNameUInt64,
	ProtoFieldTypeNameFixed32:  GolangFieldTypeNameUInt32,
	ProtoFieldTypeNameBool:     GolangFieldTypeNameBool,
	ProtoFieldTypeNameString:   GolangFieldTypeNameString,
	ProtoFieldTypeNameBytes:    GolangFieldTypeNameBytes,
	ProtoFieldTypeNameUInt32:   GolangFieldTypeNameUInt32,
	ProtoFieldTypeNameSFixed32: GolangFieldTypeNameInt32,
	ProtoFieldTypeNameSFixed64: GolangFieldTypeNameInt64,
	ProtoFieldTypeNameSInt32:   GolangFieldTypeNameInt32,
	ProtoFieldTypeNameSInt64:   GolangFieldTypeNameInt64,
}

var protoFieldTypeNameMapping = map[descriptorpb.FieldDescriptorProto_Type]ProtoFieldTypeName{
	1: ProtoFieldTypeNameDouble,
	2: ProtoFieldTypeNameFloat,
	3: ProtoFieldTypeNameInt64,
	4: ProtoFieldTypeNameUInt64,
	5: ProtoFieldTypeNameInt32,
	6: ProtoFieldTypeNameFixed64,
	7: ProtoFieldTypeNameFixed32,
	8: ProtoFieldTypeNameBool,
	9: ProtoFieldTypeNameString,
	//10: "TYPE_GROUP",
	//11: "TYPE_MESSAGE",
	12: ProtoFieldTypeNameBytes,
	13: ProtoFieldTypeNameUInt32,
	//14: "TYPE_ENUM",
	15: ProtoFieldTypeNameSFixed32,
	16: ProtoFieldTypeNameSFixed64,
	17: ProtoFieldTypeNameSInt32,
	18: ProtoFieldTypeNameSInt64,
}

type ProtoField struct {
	fd            *desc.FieldDescriptor
	Name          string   // proto field name
	Comment       *Comment // 注释
	Label         Label    // Label类型
	KeyTypeName   string   // key的proto type name
	ValueTypeName string   // value的proto type name（如果是map的话）
}

func NewProtoField(pf *ProtoFile, fd *desc.FieldDescriptor) *ProtoField {
	return &ProtoField{
		Name: fd.GetName(),
		fd:   fd,
	}
}

func (p *Parser) BuildProtoField(pf *ProtoFile, fd *desc.FieldDescriptor) *ProtoField {
	f := NewProtoField(pf, fd)
	f.Comment = p.comments[fd.AsFieldDescriptorProto()]
	switch fd.GetLabel() {
	case descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_LABEL_REQUIRED:
		f.Label = LabelOptional
		f.KeyTypeName = p.MustGetFieldTypeName(fd)
	case descriptorpb.FieldDescriptorProto_LABEL_REPEATED:
		if !fd.IsMap() {
			f.Label = LabelSlice
			f.KeyTypeName = p.MustGetFieldTypeName(fd)
		} else {
			f.Label = LabelMap
			f.KeyTypeName = p.MustGetFieldTypeName(fd.GetMapKeyType())
			f.ValueTypeName = p.MustGetFieldTypeName(fd.GetMapValueType())
		}
	}
	return f
}

func (pf *ProtoField) AsFieldDescriptor() *desc.FieldDescriptor { return pf.fd }

func ProtoFieldTypeNameToGolangFieldTypeName(protoFieldTypeName string) string {
	return protoFieldTypeNameToGolangFieldTypeNameMapping[protoFieldTypeName]
}

func FieldDescriptorProtoTypeToProtoFieldTypeName(t descriptorpb.FieldDescriptorProto_Type) string {
	return protoFieldTypeNameMapping[t]
}

func FieldDescriptorProtoTypeToGolangFieldTypeName(t descriptorpb.FieldDescriptorProto_Type) string {
	return ProtoFieldTypeNameToGolangFieldTypeName(protoFieldTypeNameMapping[t])
}

func (pf *ProtoField) KeyGoTypeName() string {
	return ProtoFieldTypeNameToGolangFieldTypeName(pf.KeyTypeName)
}

func (pf *ProtoField) ValueGoTypeName() string {
	return ProtoFieldTypeNameToGolangFieldTypeName(pf.ValueTypeName)
}
