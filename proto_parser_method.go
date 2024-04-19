package protokit

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/sandwich-go/boost/xstrings"
	"google.golang.org/protobuf/types/descriptorpb"
	"path"
	"strings"
)

const actorPathSuffix = "/actor"

func (p *Parser) method(
	protoFile *ProtoFile,
	serviceName string,
	protoMethod *descriptorpb.MethodDescriptorProto,
	md *desc.MethodDescriptor,
	isActorMethod bool,
	isAsk bool,
	fixActorMethodName bool,
	serviceUriAutoAlias bool,
	isERPCMethod bool,
	queryPath string,
	isJob bool,
) *Method {
	// Note:
	// 这里只是简单的换算一次格式合法的名称，具体请求名要通过ImportSet进行纠正
	reqTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetInputType()), ".")
	rspTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetOutputType()), ".")
	methodName := xstrings.CamelCase(protoMethod.GetName())
	if isActorMethod && (fixActorMethodName || isERPCMethod) {
		// actor 有rpc或者erpc方法
		methodName += "ForActor"
	}
	if isERPCMethod {
		// erpc 固定带这种歌后缀
		methodName += "ForERPC"
	}
	if isJob {
		methodName += "ForJob"
	}
	method := &Method{
		md:                             md,
		RpcOption:                      getRpcMethodOption(protoMethod),
		BackOfficeOption:               getBackOfficeMethodOption(protoMethod),
		JobOption:                      getJobMethodOption(protoMethod),
		Name:                           methodName,
		TypeInputDotFullQualifiedName:  protoMethod.GetInputType(),
		TypeOutputDotFullQualifiedName: protoMethod.GetOutputType(),
		TypeInputWithSelfPackage:       reqTypeName,
		TypeOutputWithSelfPackage:      rspTypeName,
		IsActor:                        isActorMethod,
		IsERPC:                         isERPCMethod,
		IsAsk:                          isAsk,
		IsTell:                         !isAsk,
	}
	if methodComment, exist := p.comments[protoMethod]; exist && methodComment != nil {
		method.Comment = methodComment.Content
	}
	fdp := protoFile.fd.AsFileDescriptorProto()
	if p.cc.URIUsingGRPCWithoutPackage {
		method.TypeInputGRPC = fmt.Sprintf("/%s/%s", serviceName, method.Name)
	} else {
		method.TypeInputGRPC = fmt.Sprintf("/%s.%s/%s", fdp.GetPackage(), serviceName, method.Name)
	}

	// 请求别名逻辑，允许proto中设定input类型别名，在请求的proto中uri将使用此名称
	// URI使用是否GRPC模式
	nameAlias := ""
	if serviceUriAutoAlias || p.cc.URIUsingGRPC {
		nameAlias = "grpc"
	}
	aliasCheckPrefer := []string{Alias}
	if isActorMethod {
		aliasCheckPrefer = []string{ActorAlias, Alias}
	}
	var anMethod methodeAnnotation
	if method.RpcOption != nil {
		anMethod = &methodOptionAnnotation{method.RpcOption}
	} else {
		anMethod = GetAnnotation(p.comments[protoMethod], AnnotationService)
	}
	for _, aliasKey := range aliasCheckPrefer {
		if anMethod.Contains(aliasKey) {
			nameAlias = anMethod.String(aliasKey)
			if nameAlias == "" {
				nameAlias = fmt.Sprintf("%s_%s", serviceName, method.Name) // 默认alias
			}
			if strings.EqualFold(nameAlias, "grpc") {
				// 如果指定为grpc，则使用grpc的路由名称
				nameAlias = method.TypeInputGRPC
			} else {
				// name alias 必须有namespace前缀，以便于激活自动转发功能，如果没有指定，则使用与TypeInput相同的package前缀
				if !strings.Contains(nameAlias, ".") {
					nameAlias = fmt.Sprintf("%s.%s", strings.Split(method.TypeInputWithSelfPackage, ".")[0], nameAlias)
				}
			}
			if aliasKey == "actor_alias" {
				// 通过actor_alias指定的别名，不再进行/actor的修正
				fixActorMethodName = false
			}
			break
		}
	}
	if strings.EqualFold(nameAlias, "grpc") {
		// 如果指定为grpc，则使用grpc的路由名称
		nameAlias = method.TypeInputGRPC
	}
	// 默认的http请求路径
	if pathStr, err := HTTPPath(protoMethod); err == nil && pathStr != "" {
		if !strings.HasPrefix(pathStr, "/") {
			pathStr = "/" + pathStr
		}
		method.HTTPPath = pathStr
		method.HTTPPathComment = "from proto, user defined"
	}
	if anMethod.Contains("http_path") {
		method.HTTPPath = anMethod.String("http_path")
		if fixActorMethodName && !strings.HasSuffix(method.HTTPPath, actorPathSuffix) {
			method.HTTPPath = path.Clean(method.HTTPPath + actorPathSuffix)
		}
		// 如果通过标注指定了http path
		nameAlias = method.HTTPPath
		method.HTTPPathComment = "from proto, user defined"
	}

	if nameAlias != "" && fixActorMethodName && !strings.EqualFold(nameAlias, method.TypeInputGRPC) && !strings.HasSuffix(nameAlias, actorPathSuffix) {
		nameAlias = path.Clean(nameAlias + actorPathSuffix)
	}

	method.HTTPPathConstName = fmt.Sprintf("%s_%s_FullHTTPName", serviceName, method.Name)

	method.FullPathHTTP = standardFullPathHTTP(method.HTTPPath, queryPath)

	method.FullPathHTTPConstName = fmt.Sprintf("%s_%s_FullPathHTTP", serviceName, method.Name)

	method.TypeInputGRPCConstName = fmt.Sprintf("%s_%s_FullGRPCName", serviceName, method.Name)
	method.TypeInputAlias = strings.TrimSpace(nameAlias)
	// {service}_{method}_FullMethodName
	method.TypeInputAliasConstName = fmt.Sprintf("%s_%s_FullMethodName", serviceName, method.Name)

	method.LangOffTag = strings.Split(anMethod.String(LangOff), ",")
	return method
}
