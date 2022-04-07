package protokit

import (
	"io"
)

type FileExcludeFilter = func(string) bool
type FileAccessor = func(fielRelativePath string) (io.ReadCloser, error)

//go:generate optionGen   --xconf=true --usage_tag_name=usage --xconf=true
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@GolangBasePackagePath(comment="golang基础package path")
		"GolangBasePackagePath": string(""),
		// annotation@GolangRelative(comment="是否启用golang relative模式")
		"GolangRelative": true,
		// annotation@ProtoImportPath(comment="proto import路径")
		"ProtoImportPath": []string{},
		// annotation@ProtoFileAccessor(comment="proto import路径")
		"ProtoFileAccessor": FileAccessor(nil),
		// annotation@ProfoExcludeFilter(comment="proto过滤")
		"ProtoFileExcludeFilter": FileExcludeFilter(func(string) bool { return false }),
		// annotation@ZapLogMapKeyTypes(comment="以类型为key的map的MarshalLogObject实现，使得可以直接使用zap.Object函数打印map数据")
		"ZapLogMapKeyTypes": []string{"int", "int32", "int64", "uint32", "uint64", "string"},
		// annotation@ZapLogBytesMode(comment="zap以何种方式输出[]byte, 可以使用base64或者bytes, 默认bytes")
		"ZapLogBytesMode": "bytes",
		// annotation@NamePattern(xconf="namepattern",comment="NamePattern",inline="true")
		"NamePattern": (*NamePattern)(NewNamePattern()),
	}
}

//go:generate optionGen   --xconf=true --usage_tag_name=usage --xconf=true
func NamePatternOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@NamePatternServerHandler(comment="code server handler名称格式化")
		"NamePatternServerHandler": "ServerHandler%s",
		// annotation@NamePatternRpcClient(comment="code rpc client名称格式化")
		"NamePatternRpcClient": "RpcClient%s",
		// annotation@NamePatternActorClient(comment="code actor client名称格式化")
		"NamePatternActorClient": "ActorClient%s",
		// annotation@NamePatternActorServerHandlerMethod(comment="actor server handler的方法名称格式化")
		"NamePatternActorServerHandlerMethod": "%s",
		// annotation@NamePatternSingleConfPackage(comment="生成golang conf时控制package名称，可用参数#file、#sheet")
		"NamePatternSingleConfPackage": "%s",
	}
}

// fixme 主要是为了摒除netutils部分的影响，这里可以考虑将template中预知的import移到ImportSet
var excludeImportName = []string{
	"netutils",
}
