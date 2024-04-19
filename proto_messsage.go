package protokit

import (
	"strings"

	protokit2 "github.com/sandwich-go/protokit/option/gen/golang/protokit"
	"google.golang.org/protobuf/proto"

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
	Store                         map[interface{}]interface{}
	Parser                        *Parser
}

func NewProtoMessage(pf *ProtoFile, md *desc.MessageDescriptor) *ProtoMessage {
	pm := &ProtoMessage{
		md:                            md,
		ProtoFile:                     pf,
		Name:                          md.GetName(),
		Fields:                        make([]*ProtoField, 0, len(md.GetFields())),
		goStructNameWithGolangPackage: GoStructNameWithGolangPackage(md.GetFullyQualifiedName(), pf.Package, pf.GolangPackageName),
		ImportSet:                     NewImportSet(pf.GolangPackageName, pf.GolangPackagePath),
		Store:                         map[interface{}]interface{}{},
	}
	return pm
}

func (p *Parser) BuildProtoMessage(pf *ProtoFile, md *desc.MessageDescriptor) *ProtoMessage {
	pm := NewProtoMessage(pf, md)
	pm.dotFullyQualifiedTypeName = p.descriptor2DotFullyQualifiedTypeName[pm.md]
	pm.Comment = p.comments[pm.md.AsDescriptorProto()]
	pm.Parser = p
	return pm
}

func (pm *ProtoMessage) AddToStore(k, v interface{})                  { pm.Store[k] = v }
func (pm *ProtoMessage) GetFromStore(k interface{}) interface{}       { return pm.Store[k] }
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

func (pm *ProtoMessage) GetOrmOption() *protokit2.OrmMessageOptions {
	msgO := pm.AsMessageDescriptor().GetMessageOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_OrmMessage).(*protokit2.OrmMessageOptions)
	if ok {
		return opts
	}
	return nil
}
