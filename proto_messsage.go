package protokit

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
)

type ProtoMessage struct {
	dotFullyQualifiedTypeName     string
	md                            *desc.MessageDescriptor
	Name                          string // proto message name
	goStructNameWithGolangPackage string
	ProtoFile                     *ProtoFile
	Fields                        []*ProtoField // 所有的属性
	Comment                       *Comment      // 注释
	ImportSet                     *ImportSet
}

func NewProtoMessage(pf *ProtoFile, md *desc.MessageDescriptor, cc *Options) *ProtoMessage {
	pm := &ProtoMessage{
		md:                            md,
		ProtoFile:                     pf,
		Name:                          md.GetName(),
		Fields:                        make([]*ProtoField, 0, len(md.GetFields())),
		goStructNameWithGolangPackage: GoStructNameWithGolangPackage(md.GetFullyQualifiedName(), pf.Package, pf.GolangPackageName),
		ImportSet:                     NewImportSet(pf.GolangPackageName, pf.GolangPackagePath, cc.ImportSetExclude),
	}
	return pm
}

func (p *Parser) BuildProtoMessage(pf *ProtoFile, md *desc.MessageDescriptor) *ProtoMessage {
	pm := NewProtoMessage(pf, md, p.cc)
	pm.dotFullyQualifiedTypeName = p.descriptor2DotFullyQualifiedTypeName[pm.md]
	pm.Comment = p.comments[pm.md.AsDescriptorProto()]
	return pm
}

func (pm *ProtoMessage) AsMessageDescriptor() *desc.MessageDescriptor { return pm.md }
func (pm *ProtoMessage) DotFullyQualifiedTypeName() string            { return pm.dotFullyQualifiedTypeName }
func (pm *ProtoMessage) GoStructNameWithGolangPackage() string {
	return pm.goStructNameWithGolangPackage
}

func (pm *ProtoMessage) GoStructNameWithoutGolangPackage() string {
	ss := strings.Split(pm.goStructNameWithGolangPackage, ".")
	return ss[len(ss)-1]
}

func (pm *ProtoMessage) HasCommentField(comment string) bool {
	comment = strings.ToLower(comment)
	for _, f := range pm.Fields {
		if f.Comment == nil {
			continue
		}
		for k := range f.Comment.Tags {
			if k == comment {
				return true
			}
		}
	}
	return false
}
