package protokit

import (
	"io"
	"strings"
)

type FileExcludeFilter = func(string) bool
type FileAccessor = func(fielRelativePath string) (io.ReadCloser, error)

// 默认过滤器
func defaultFileExcludeFilter(filePath string) bool {
	return strings.Contains(filePath, "rawdata_client")
}

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
		// annotation@ProtoFileExcludeFilter(comment="proto过滤")
		"ProtoFileExcludeFilter": FileExcludeFilter(defaultFileExcludeFilter),
		// annotation@ZapLogMapKeyTypes(comment="以类型为key的map的MarshalLogObject实现，使得可以直接使用zap.Object函数打印map数据")
		"ZapLogMapKeyTypes": []string{"int", "int32", "int64", "uint32", "uint64", "string"},
		// annotation@ZapLogBytesMode(comment="zap以何种方式输出[]byte, 可以使用base64或者bytes, 默认bytes")
		"ZapLogBytesMode": "bytes",
		// annotation@NamePattern(comment="名称格式化空自己",inline="true")
		"NamePattern": (*NamePattern)(NewNamePattern()),
		// annotation@ImportSetExclude(comment="import set忽略指定name的package")
		"ImportSetExclude": []string{"netutils"},
		// annotation@URIUsingGRPC(comment="service的uri是否使用GRPC模式")
		"URIUsingGRPC": false,
		// annotation@InvalidServiceAnnotations(comment="非法的 service annotations")
		"InvalidServiceAnnotations": []string{},
		// annotation@URIUsingGRPCWithoutPackage(comment="service的uri使用GRPC模式时，是否带package名")
		"URIUsingGRPCWithoutPackage": false,
		// annotation@StrictMode(comment="是否为严格模式")
		"StrictMode": true,
		// annotation@QueryPathMapping(comment="query path映射关系,通过 {{key}} 方式访问值")
		"QueryPathMapping": (map[string]string)(map[string]string{
			"root": "/",
		}),
		// annotation@DefaultQueryPath(comment="默认query path，支持配置 {{key}}的方式索引QueryPathMapping的key")
		"DefaultQueryPath": "/",
		// annotation@ForceGrpcStyle(comment="rpc actor erpc改用grpc风格的名称，不可同时存在")
		"ForceGrpcStyle": false,
	}
}

//go:generate optionGen --xconf=true --usage_tag_name=usage --xconf=true
func NamePatternOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@NamePatternServerHandler(comment="code server handler名称格式化")
		"NamePatternServerHandler": "ServerHandler%s",
		// annotation@NamePatternRPCClient(comment="code rpc client名称格式化")
		"NamePatternRPCClient": "RPCClient%s",
		// annotation@NamePatternActorClient(comment="code actor client名称格式化")
		"NamePatternActorClient": "ActorClient%s",
		// annotation@NamePatternERPCClient(comment="code erpc client名称格式化")
		"NamePatternERPCClient": "ERPCClient%s",
		// annotation@NamePatternHTTPPath(comment="自动生成的HTTP PATHG格式")
		"NamePatternHTTPPath": "%s",
		// annotation@NamePatternJobClient(comment="job client 的名称格式化")
		"NamePatternJobClient": "JobClient%s",
		// annotation@NamePatternJobService(comment="job service 的名称格式化")
		"NamePatternJobService": "JobService%s",
		// annotation@NamePatternJobMethod(comment="job method 的名称格式化")
		"NamePatternJobMethod": "%sForJob",
	}
}
