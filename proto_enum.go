package protokit

import "github.com/jhump/protoreflect/desc"

type ProtoEnum struct {
	dotFullyQualifiedTypeName string
	ed                        *desc.EnumDescriptor
	ProtoFile                 *ProtoFile
	Name                      string // proto enum name
}

func NewProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	return &ProtoEnum{
		ed:        ed,
		ProtoFile: pf,
		Name:      ed.GetName(),
	}
}

func (pe *ProtoEnum) DotFullyQualifiedTypeName() string      { return pe.dotFullyQualifiedTypeName }
func (pe *ProtoEnum) AsEnumDescriptor() *desc.EnumDescriptor { return pe.ed }

func (p *Parser) BuildProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	pe := NewProtoEnum(pf, ed)
	pe.dotFullyQualifiedTypeName = p.descriptor2DotFullyQualifiedTypeName[pe.ed]
	return pe
}
