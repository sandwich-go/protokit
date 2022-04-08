package protokit

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/proto"

	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xpanic"
)

type ParserVisitor interface {
	// NewMessage 新建message
	NewMessage(name string) (*dynamic.Message, bool)
	// RangeDotFullyQualifiedTypeNameToDescriptor 遍历dot名称到Descriptor的映射
	RangeDotFullyQualifiedTypeNameToDescriptor(f func(string, desc.Descriptor) bool)
	// Comment 访问类型对应的注释信息
	Comment(m proto.Message) (*Comment, bool)
	// DotFullyQualifiedTypeNameToDescriptor 由dot名称获取Descriptor
	DotFullyQualifiedTypeNameToDescriptor(dotName string) (desc.Descriptor, bool)
	// DotFullyQualifiedTypeNameToProtoFile 由dot名称获取对应的ProtoFile
	DotFullyQualifiedTypeNameToProtoFile(dotName string) (*ProtoFile, bool)
	// DotFullyQualifiedTypeNameToProtoMessage 由dot名称获取对应的ProtoMessage
	DotFullyQualifiedTypeNameToProtoMessage(dotName string) (*ProtoMessage, bool)
	// DescriptorToDotFullyQualifiedTypeName 由Descriptor获取dot名称
	DescriptorToDotFullyQualifiedTypeName(d desc.Descriptor) (string, bool)
}

type Parser struct {
	cc                                      *Options
	comments                                map[proto.Message]*Comment
	protoFilePathToProtoFile                map[string]*ProtoFile        // proto 文件路径(相对路径) => *ProtoFile
	dotFullyQualifiedTypeNameToProtoFile    map[string]*ProtoFile        // 类型名(.package_path.message) => *ProtoFile
	dotFullyQualifiedTypeNameToDescriptor   map[string]desc.Descriptor   // 类型名(.package_path.message) => desc.Descriptor
	dotFullyQualifiedTypeNameToProtoMessage map[string]*ProtoMessage     // 类型名(.package_path.message) => *ProtoMessage
	descriptor2DotFullyQualifiedTypeName    map[desc.Descriptor]string   // desc.Descriptor => 类型名(.package_path.message)
	tmpParsedMessageOrEnumMapping           map[desc.Descriptor]struct{} // 已经解析过的message或者enum，临时的
}

func NewParser(opts ...Option) *Parser {
	p := &Parser{cc: NewOptions(opts...)}
	p.Clean()
	return p
}

func (p *Parser) Options() *Options { return p.cc }

func (p *Parser) NewMessage(name string) (*dynamic.Message, bool) {
	nameWithDot := name
	if !strings.HasPrefix(nameWithDot, ".") {
		nameWithDot = "." + nameWithDot
	}
	protoFile, ok := p.dotFullyQualifiedTypeNameToProtoFile[nameWithDot]
	if !ok {
		return nil, ok
	}
	nameWithoutDot := strings.TrimPrefix(nameWithDot, ".")
	typeDesc := p.protoFilePathToProtoFile[protoFile.FilePath].fd.FindMessage(nameWithoutDot)
	return dynamic.NewMessage(typeDesc), ok
}

func (p *Parser) Clean() {
	p.comments = make(map[proto.Message]*Comment)
	p.protoFilePathToProtoFile = make(map[string]*ProtoFile)
	p.dotFullyQualifiedTypeNameToProtoFile = make(map[string]*ProtoFile)
	p.dotFullyQualifiedTypeNameToDescriptor = make(map[string]desc.Descriptor)
	p.dotFullyQualifiedTypeNameToProtoMessage = make(map[string]*ProtoMessage)
	p.descriptor2DotFullyQualifiedTypeName = make(map[desc.Descriptor]string)
	p.tmpParsedMessageOrEnumMapping = make(map[desc.Descriptor]struct{})
}

func (p *Parser) Parse(nsList ...*Namespace) {
	parser := &protoparse.Parser{
		IncludeSourceCodeInfo: true,
		Accessor:              p.cc.ProtoFileAccessor,
	}
	parser.ImportPaths = append(parser.ImportPaths, p.cc.ProtoImportPath...)
	for _, ns := range nsList {
		ns.ParserVisitor = p
		// clean一次移除path前缀中的诸如./之类的字符便于获取文件的相对路径，稳妥起见，外层直接使用绝对路径最为安全
		pathProtoRoot := path.Clean(ns.Path)
		// 获取文件列表
		fileListAbs := make([]string, 0)
		err := xos.FilePathWalkFollowLink(pathProtoRoot, xos.FileWalkFuncWithExcludeFilter(&fileListAbs, p.cc.ProtoFileExcludeFilter, ".proto"))
		xpanic.PanicIfErrorAsFmtFirst(err, "got error: %w while walk dir:%s", pathProtoRoot)
		// 路径替换为相对路径，Parser需求，按照相对proto文件名查找依赖
		fileList := make([]string, len(fileListAbs))
		for index, filePath := range fileListAbs {
			fileList[index] = strings.TrimLeft(strings.Replace(filePath, pathProtoRoot, "", 1), string(os.PathSeparator))
			// proto import使用/, 兼容windows，替换路径中的\为/
			fileList[index] = strings.ReplaceAll(fileList[index], `\`, `/`)
		}
		// 按照FileDescriptor解析所有proto文件
		var fds []*desc.FileDescriptor
		fds, err = parser.ParseFiles(fileList...)
		xpanic.PanicIfErrorAsFmtFirst(err, "got error: %w while parse files under dir:%s", pathProtoRoot)
		for index, fd := range fds {
			golangPackagePath, golangPackageName := GolangPackagePathAndName(fd, p.cc.GolangBasePackagePath, p.cc.GolangRelative)
			pf := NewProtoFile(golangPackageName, golangPackagePath)
			pf.GolangRelative = p.cc.GolangRelative
			pf.Namespace = ns.Name
			pf.FilePath = fd.GetName()
			filePath := fileListAbs[index]
			bb, err := xos.FileGetContents(filePath)
			xpanic.PanicIfErrorAsFmtFirst(err, "got error: %w while load file content:%s", filePath)
			pf.Content = string(bb)
			pf.Package = fd.GetPackage()
			pf.OptionGolangPackage = fd.AsFileDescriptorProto().GetOptions().GetGoPackage()
			pf.OptionCSNamespace = fd.AsFileDescriptorProto().GetOptions().GetCsharpNamespace()
			pf.fd = fd
			// comment
			p.parseComments(fd)
			// 存储文件名到fd的映射，便于根据文件查找fd
			p.protoFilePathToProtoFile[pf.FilePath] = pf
			// 填充当前ns内的文件
			ns.Files[pf.FilePath] = pf
		}
	}
	// type 映射关系
	p.buildTypeNameMap()
	// 解析message
	p.parseMessages()
	// 解析service,依赖于buildTypeNameMap
	p.parseService()
	// 解析import
	p.parseImport()
	// 填充字段Validator信息
	p.parseValidatorForMethod()
	// 解析package
	p.parsePackage(nsList)
	p.parseAnnotation()
}

func (p *Parser) setType(mdp *desc.MessageDescriptor, name string, pf *ProtoFile) {
	p.dotFullyQualifiedTypeNameToProtoFile[name] = pf
	p.dotFullyQualifiedTypeNameToDescriptor[name] = mdp
	for _, v := range mdp.GetNestedMessageTypes() {
		// 过滤掉map entry,
		if v.IsMapEntry() {
			continue
		}
		p.setType(v, "."+v.GetFullyQualifiedName(), pf)
	}
	for _, v := range mdp.GetNestedEnumTypes() {
		dotName := "." + v.GetFullyQualifiedName()
		p.dotFullyQualifiedTypeNameToProtoFile[dotName] = pf
		p.dotFullyQualifiedTypeNameToDescriptor[dotName] = v
	}
}

// BuildTypeNameMap builds the map from fully qualified type names to objects.
// The key names for the map come from the input data, which puts a period at the beginning.
// It should be called after SetPackageNames and before GenerateAllFiles.
func (p *Parser) buildTypeNameMap() {
	for _, pf := range p.protoFilePathToProtoFile {
		// The names in this loop are defined by the proto world, not us, so the
		// package name may be empty.  If so, the dotted package name of X will
		// be ".X"; otherwise it will be ".pkg.X".
		dottedPkg := p.getDottedPackage(pf.fd)
		for _, mt := range pf.fd.GetMessageTypes() {
			p.setType(mt, dottedPkg+mt.GetName(), pf)
		}
		for _, v := range pf.fd.GetEnumTypes() {
			dotName := dottedPkg + v.GetName()
			p.dotFullyQualifiedTypeNameToProtoFile[dotName] = pf
			p.dotFullyQualifiedTypeNameToDescriptor[dotName] = v
		}
	}
	for k, v := range p.dotFullyQualifiedTypeNameToDescriptor {
		p.descriptor2DotFullyQualifiedTypeName[v] = k
	}
}

func (p *Parser) MustGetFieldTypeName(fd *desc.FieldDescriptor) string {
	if tn, ok := protoFieldTypeNameMapping[fd.GetType()]; ok {
		return tn
	}
	if mt := fd.GetMessageType(); mt != nil {
		return p.descriptor2DotFullyQualifiedTypeName[mt]
	}
	if et := fd.GetEnumType(); et != nil {
		return p.descriptor2DotFullyQualifiedTypeName[et]
	}
	// fixme, group?
	panic(fmt.Sprintf("Unknown field type, %s", fd.GetType()))
}

func (p *Parser) DotFullyQualifiedTypeNameToDescriptor(dotName string) (desc.Descriptor, bool) {
	v, ok := p.dotFullyQualifiedTypeNameToDescriptor[dotName]
	if ok {
		return v, true
	}
	return nil, false
}

func (p *Parser) DotFullyQualifiedTypeNameToProtoFile(dotName string) (*ProtoFile, bool) {
	v, ok := p.dotFullyQualifiedTypeNameToProtoFile[dotName]
	if ok {
		return v, true
	}
	return nil, false
}

func (p *Parser) DotFullyQualifiedTypeNameToProtoMessage(dotName string) (*ProtoMessage, bool) {
	v, ok := p.dotFullyQualifiedTypeNameToProtoMessage[dotName]
	if ok {
		return v, true
	}
	return nil, false
}

func (p *Parser) DescriptorToDotFullyQualifiedTypeName(d desc.Descriptor) (string, bool) {
	v, ok := p.descriptor2DotFullyQualifiedTypeName[d]
	if ok {
		return v, true
	}
	return "", false
}

func (p *Parser) Comment(m proto.Message) (*Comment, bool) {
	v, ok := p.comments[m]
	if ok {
		return v, true
	}
	return nil, false
}

func (p *Parser) RangeDotFullyQualifiedTypeNameToDescriptor(f func(string, desc.Descriptor) bool) {
	for k, v := range p.dotFullyQualifiedTypeNameToDescriptor {
		if !f(k, v) {
			break
		}
	}
}
