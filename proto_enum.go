package protokit

import "github.com/jhump/protoreflect/desc"

type ProtoEnumField struct {
	Name   string
	Number int32
}
type ProtoEnum struct {
	dotFullyQualifiedTypeName string
	ed                        *desc.EnumDescriptor
	Fields                    []*ProtoEnumField
	ProtoFile                 *ProtoFile
	Name                      string // proto enum name
}

func NewProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	v := &ProtoEnum{
		ed:        ed,
		ProtoFile: pf,
		Name:      ed.GetName(),
	}
	for _, field := range v.ed.AsEnumDescriptorProto().Value {
		v.Fields = append(v.Fields, &ProtoEnumField{
			Name:   field.GetName(),
			Number: field.GetNumber(),
		})
	}
	return v
}

func (pe *ProtoEnum) DotFullyQualifiedTypeName() string      { return pe.dotFullyQualifiedTypeName }
func (pe *ProtoEnum) AsEnumDescriptor() *desc.EnumDescriptor { return pe.ed }

func (p *Parser) BuildProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	pe := NewProtoEnum(pf, ed)
	pe.dotFullyQualifiedTypeName = p.descriptor2DotFullyQualifiedTypeName[pe.ed]
	return pe
}
