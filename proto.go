package protokit

import (
	"path"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
)

type Import struct {
	ProtoFilePath               string   // 引入的proto文件名
	GolangPackageName           string   // golang 引用的名称
	GolangPackagePath           string   // golang 引用的路径
	PythonModuleName            string   // python 引用的名称
	PythonModulePath            string   // python 引用的文件路径
	MessageDotFullQualifiedName []string // 当前import下引入的类型列表
}

type ImportSet struct {
	GolangPackageName                       string            // 宿主文件的package name
	GolangPackagePath                       string            // 宿主文件的package path
	Set                                     []*Import         // 当前import set
	MessageDotFullQualifiedNameToGolangType map[string]string // GolangType是经过import纠正过package名称的，可能带着1，2之类的标记
	PythonModules                           []*PythonModule   // python module 辅助python代码生成
	importAliasMappingCount                 map[string]int    // 构建中使用的临时数据，记录同名但不同路径的import
	ExcludeImportName                       []string
}

func NewImportSet(golangPackageName, golangPackagePath string) *ImportSet {
	return &ImportSet{
		GolangPackageName:                       golangPackageName,
		GolangPackagePath:                       golangPackagePath,
		MessageDotFullQualifiedNameToGolangType: make(map[string]string),
		importAliasMappingCount:                 make(map[string]int),
	}
}

type Method struct {
	md              *desc.MethodDescriptor
	Name            string // 方法名称，proto中获取到的原始名称
	Comment         string // method注释
	ValidatorInput  bool   // 是否检验输入
	ValidatorOutput bool   // 是否校验输出
	// Note: golang与python使用相同的名称，类型名是golang规则,对于嵌套结构,为python生成一套类型别名
	TypeInput                      string // Import校正后的名称，携带package信息
	TypeOutput                     string // Import校正后的名称，携带package信息
	TypeInputAlias                 string // Input别名
	TypeInputGRPC                  string // GRPC模式下的Input路径
	HTTPPath                       string // HTTP模式下的请求路径
	HTTPPathComment                string // HTTP模式下的请求路径注释，来源
	IsAsk                          bool   // 是否为Ask方法
	IsTell                         bool   // 是否为Tell方法
	IsActor                        bool   // 是否为Actor方法
	TypeInputDotFullQualifiedName  string // proto原始Input，也就是DotFullQualifiedName
	TypeOutputDotFullQualifiedName string // proto原始Output，也就是DotFullQualifiedName
	TypeInputWithSelfPackage       string //  只携带自身package信息
	TypeOutputWithSelfPackage      string //  只携带自身package信息
}

func (m *Method) AsMethodDescriptor() *desc.MethodDescriptor { return m.md }

type Service struct {
	sd                         *desc.ServiceDescriptor
	Name                       string    // 通过proto获取到的原始名字
	ServiceName                string    // 当前服务的名称，格式化后的，数据源:RPCClientInterfaceName/ServerHandlerInterfaceName/ActorClientInterfaceName
	ServerHandlerInterfaceName string    // Server Handler名称
	RPCClientInterfaceName     string    // RPC Client名称
	ActorClientInterfaceName   string    // Actor Client 名称
	Comment                    string    // 注释信息
	DeprecatedName             string    // 兼容数据，弃用的结构名称
	Methods                    []*Method // 当前服务中的Method列表
	InputOutputTypes           []string  // 当前服务内使用的消息列表，用于加速uri生成，rpc actor中使用
	HasActorMethod             bool      // 辅助生成的时候是否import actor包
	HasValidator               bool      // 辅助生成的时候是否携带validator包
	DescName                   string    // fdp.GetPackage().Name
	DescProtoFile              string    // fdp.GetName() 应该是ProtoFile.FilePath
}

func (s *Service) AsServiceDescriptor() *desc.ServiceDescriptor { return s.sd }

type ServiceGroup struct {
	ProtoFilePath string     // 当前group所属的ProtoFilePath,可以通过ProtoFilePath唯一索引到ProtoFile
	Services      []*Service // 当前group内的service列表，一个proto文件内可能有多个service
	ImportSet     *ImportSet // 同一个ServiceGroup内的service共享同一个ImportSet，目的是生成到同一个文件
}

type ServiceTag = string

const (
	ServiceTagALL   ServiceTag = "all"
	ServiceTagRPC   ServiceTag = "rpc"
	ServiceTagActor ServiceTag = "actor"
)

var allServiceTags = []ServiceTag{ServiceTagALL, ServiceTagRPC, ServiceTagActor}

type ProtoFile struct {
	Namespace           string                   // 当前文件所属的NameSpace名称，在构建package信息的时候需要使用
	FilePath            string                   // fd.GetName() proto文件路径，相对路径，也是引用时使用的路径
	Package             string                   // package名称
	GolangPackageName   string                   // protokitgo逻辑内计算的golang package name
	GolangPackagePath   string                   // protokitgo逻辑内计算的golang package path
	OptionGolangPackage string                   // 原始proto文件Option中的数据，可能为空
	OptionCSNamespace   string                   // 原始proto文件Option中的数据，可能为空
	Content             []byte                   // proto文件内容
	Messages            []*ProtoMessage          // proto文件中所有的ProtoMessage
	Enums               []*ProtoEnum             // proto文件中所有的ProtoEnum
	ServiceGroups       map[string]*ServiceGroup // 带tag的service，用于处理用于处理RPC/Actor Client等逻辑，service可能不是全量的
	GolangRelative      bool                     // 是否使用的是golang relative模式,数据来源于option
	fd                  *desc.FileDescriptor     // 文件对应的FileDescriptor
}

func (p *ProtoFile) AsFileDescriptor() *desc.FileDescriptor { return p.fd }

func (p *ProtoFile) GetFullPathWithSuffix(suffix string) string {
	name := p.FilePath
	if ext := filepath.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}
	if suffix == "" {
		suffix = ".go"
	}
	name += suffix
	if p.GolangRelative {
		return name
	}
	// Replace the existing dirname with the declared import path.
	_, name = path.Split(name)
	return path.Join(p.GolangPackagePath, name)
}

type Package struct {
	GolangRelative          bool              // 使用使用golang relative模式,数据来源于option
	GolangPackageName       string            // protokitgo逻辑内计算的golang package name
	GolangPackagePath       string            // protokitgo逻辑内计算的golang package path
	Package                 string            // package名称，proto原始数据，来源于package下的某一个ProtoFile
	FilePath                string            // package所在目录，来源于package下的某一个ProtoFile.FilePath的目录部分
	ImportSet               *ImportSet        // 当前保内的所有消息，主要用于辅助golang message registry生成
	MessageRegistryAutoInit bool              // golang message registry是否自动生成init
	AliasToGolangType       map[string]string // alias到golang typ的名字映射，golang type名称是经过package修正的
	ActorMessageGolangType  []string          // actor消息列表，golang type名称是经过package修正的
	IsGlobal                bool              // 是否为虚拟的全局package
}

func (p *Package) GetFullPathWithFileName(fileName string) string {
	name := path.Join(p.FilePath, fileName)
	if p.GolangRelative {
		return name
	}
	// Replace the existing dirname with the declared import path.
	_, name = path.Split(name)
	return path.Join(p.GolangPackagePath, name)
}

// GetFullPathWithFileNameIgnoreGlocalPackageDir global package忽略pakcage名称
func (p *Package) GetFullPathWithFileNameIgnoreGlocalPackageDir(fileName string) string {
	name := path.Join(p.FilePath, fileName)
	if p.IsGlobal {
		name = fileName
	}
	if p.GolangRelative {
		return name
	}
	// Replace the existing dirname with the declared import path.
	_, name = path.Split(name)
	return path.Join(p.GolangPackagePath, name)
}

func NewPackageWithPackageName(golangPackageName, golangPackagePath string) *Package {
	p := &Package{}
	p.GolangPackageName = golangPackageName
	p.GolangPackagePath = golangPackagePath
	p.AliasToGolangType = make(map[string]string)
	p.ImportSet = NewImportSet(golangPackageName, golangPackagePath)
	return p
}

func NewProtoFile(golangPackageName, golangPackagePath string) *ProtoFile {
	return &ProtoFile{
		GolangPackageName: golangPackageName,
		GolangPackagePath: golangPackagePath,
		ServiceGroups:     make(map[string]*ServiceGroup),
	}
}

const (
	NamespaceGoogle   = "google"   // google sdk
	NamespaceNetutils = "netutils" // netutils sdk
	NamespaceUser     = "user"     // user proto files
)

// NamespaceMessageRegistryPackageName namespace根目录下聚合message注册的包名
const NamespaceMessageRegistryPackageName = "message_registry"

type DescriptorAccessor func(dotName string) (desc.Descriptor, bool)
type CommentAccessor func(proto.Message) (*Comment, bool)
type DotFullyQualifiedTypeNameAccessor func(desc.Descriptor) (string, bool)
type DotFullyQualifiedTypeNameToProtoFileAccessor func(name string) (*ProtoFile, bool)

// Namespace 一个proto file的聚合，为逻辑层单独抽象的概念
type Namespace struct {
	Options  *Options
	Name     string                // 用户指定的namespace，不允许重复,逻辑上的概念用于区分不同的根目录
	Path     string                // 加载时的路径
	Files    map[string]*ProtoFile // proto文件
	Packages map[string]*Package   // package列表,key为golang package path

	// 临时方案，兼容部分未调整的export逻辑
	ParserVisitor
}

// NewNamespace 新建一个namespace
func NewNamespace(name string, path string) *Namespace {
	messageRegistryPackage := NewPackageWithPackageName(NamespaceMessageRegistryPackageName, NamespaceMessageRegistryPackageName)
	messageRegistryPackage.MessageRegistryAutoInit = true
	messageRegistryPackage.FilePath = NamespaceMessageRegistryPackageName
	messageRegistryPackage.IsGlobal = true
	return &Namespace{
		Name:     name,
		Path:     path,
		Files:    make(map[string]*ProtoFile),
		Packages: map[string]*Package{NamespaceMessageRegistryPackageName: messageRegistryPackage},
	}
}
