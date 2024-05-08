package protokit

import (
	"github.com/jhump/protoreflect/desc"
	protokit2 "github.com/sandwich-go/protokit/option/gen/golang/protokit"
	"google.golang.org/protobuf/proto"
	"strings"
)

type ProtoEnumField struct {
	Name   string
	Number int32
	Field  *desc.EnumValueDescriptor
}
type ProtoEnum struct {
	dotFullyQualifiedTypeName string
	ed                        *desc.EnumDescriptor
	Fields                    []*ProtoEnumField
	ProtoFile                 *ProtoFile
	goNameWithGolangPackage   string
	Name                      string // proto enum name
}

func NewProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	v := &ProtoEnum{
		ed:                      ed,
		ProtoFile:               pf,
		Name:                    ed.GetName(),
		goNameWithGolangPackage: GoStructNameWithGolangPackage(ed.GetFullyQualifiedName(), pf.Package, pf.GolangPackageName),
	}
	for _, field := range v.ed.GetValues() {
		v.Fields = append(v.Fields, &ProtoEnumField{
			Name:   field.GetName(),
			Number: field.GetNumber(),
			Field:  field,
		})
	}
	return v
}

func (pe *ProtoEnum) DotFullyQualifiedTypeName() string      { return pe.dotFullyQualifiedTypeName }
func (pe *ProtoEnum) AsEnumDescriptor() *desc.EnumDescriptor { return pe.ed }
func (pe *ProtoEnum) GoNameWithGolangPackage() string {
	return pe.goNameWithGolangPackage
}
func (pe *ProtoEnum) GoNameWithoutGolangPackage() string {
	ss := strings.Split(pe.goNameWithGolangPackage, ".")
	return ss[len(ss)-1]
}
func (p *Parser) BuildProtoEnum(pf *ProtoFile, ed *desc.EnumDescriptor) *ProtoEnum {
	pe := NewProtoEnum(pf, ed)
	pe.dotFullyQualifiedTypeName = p.descriptor2DotFullyQualifiedTypeName[pe.ed]
	return pe
}

func (pe *ProtoEnum) GetOrmOption() *protokit2.OrmEnumOptions {
	opts, ok := proto.GetExtension(pe.AsEnumDescriptor().GetEnumOptions(), protokit2.E_OrmEnum).(*protokit2.OrmEnumOptions)
	if ok {
		return opts
	}
	return nil
}

func (pe *ProtoEnumField) GetEnumValueOptions() *protokit2.EnumValueOptions {
	opts, ok := proto.GetExtension(pe.Field.GetEnumValueOptions(), protokit2.E_EnumValue).(*protokit2.EnumValueOptions)
	if ok {
		return opts
	}
	return nil
}
