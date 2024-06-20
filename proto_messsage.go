package protokit

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/jhump/protoreflect/desc"
	protokit2 "github.com/sandwich-go/protokit/option/gen/golang/protokit"
	"google.golang.org/protobuf/proto"
	"strings"
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
	ValidateOptions               *validate.MessageConstraints
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
	if o, ok := proto.GetExtension(md.AsDescriptorProto().GetOptions(), validate.E_Message).(*validate.MessageConstraints); ok {
		pm.ValidateOptions = o
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

func (pm *ProtoMessage) HasValidateOption() bool {
	if pm.ValidateOptions != nil {
		return true
	}
	for _, f := range pm.Fields {
		if f.ValidateOptions != nil {
			return true
		}
	}
	return false
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

func (pm *ProtoMessage) OrmNoLog() bool {
	opts := pm.GetOrmOption()
	if opts == nil {
		return false
	}
	return opts.GetNoLog()
}

func (pm *ProtoMessage) OrmAlias() *OrmFieldAlias {
	opts := pm.GetOrmOption()
	if opts == nil {
		return nil
	}
	bean := opts.GetBean()
	if bean == nil {
		return nil
	}
	alias := bean.GetDynamicGeneric()
	if alias == nil {
		return nil
	}
	return (*OrmFieldAlias)(alias)
}
