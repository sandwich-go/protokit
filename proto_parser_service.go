package protokit

import (
	"fmt"
	"path"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/sandwich-go/boost/xstrings"
	"google.golang.org/protobuf/types/descriptorpb"
)

func nameMustHaveSuffix(s string, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}
func (p *Parser) parseService() {
	for _, protoFile := range p.protoFilePathToProtoFile {
		for _, serviceTag := range allServiceTags {
			protoFile.ServiceGroups[serviceTag] = &ServiceGroup{
				ProtoFilePath: protoFile.FilePath,
				Services:      p.parseServiceForProtoFile(protoFile, serviceTag),
				ImportSet:     NewImportSet(protoFile.GolangPackageName, protoFile.GolangPackagePath),
			}
		}
	}
}

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
) *Method {
	// Note:
	// 这里只是简单的换算一次格式合法的名称，具体请求名要通过ImportSet进行纠正
	reqTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetInputType()), ".")
	rspTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetOutputType()), ".")
	methodName := xstrings.CamelCase(protoMethod.GetName())
	if isActorMethod && fixActorMethodName {
		methodName += "ForActor"
	}
	method := &Method{
		md:                             md,
		Name:                           methodName,
		TypeInputDotFullQualifiedName:  protoMethod.GetInputType(),
		TypeOutputDotFullQualifiedName: protoMethod.GetOutputType(),
		TypeInputWithSelfPackage:       reqTypeName,
		TypeOutputWithSelfPackage:      rspTypeName,
		IsActor:                        isActorMethod,
		IsAsk:                          isAsk,
		IsTell:                         !isAsk,
	}
	if methodComment, exist := p.comments[protoMethod]; exist && methodComment != nil {
		method.Comment = methodComment.Content
	}
	fdp := protoFile.fd.AsFileDescriptorProto()
	method.TypeInputGRPC = fmt.Sprintf("/%s.%s/%s", fdp.GetPackage(), serviceName, method.Name)
	// 请求别名逻辑，允许proto中设定input类型别名，在请求的proto中uri将使用此名称
	// URI使用是否GRPC模式
	nameAlias := ""
	if serviceUriAutoAlias || p.cc.URIUsingGRPC {
		nameAlias = "grpc"
	}
	aliasCheckPrefer := []string{"alias"}
	if isActorMethod {
		aliasCheckPrefer = []string{"actor_alias", "alias"}
	}
	anMethod := GetAnnotation(p.comments[protoMethod], AnnotationService)
	for _, aliasKey := range aliasCheckPrefer {
		if anMethod.Has(aliasKey) {
			nameAlias = fmt.Sprintf("%s_%s", serviceName, method.Name) // 默认alias
			if autoAlias := anMethod.GetBool(aliasKey, false); !autoAlias {
				// proto中指定了alias名称
				nameAlias = anMethod.GetString(aliasKey)
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
	if anMethod.Has("http_path") {
		method.HTTPPath = anMethod.GetString("http_path")
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
	method.TypeInputAlias = strings.TrimSpace(nameAlias)
	method.LangOffTag = strings.Split(anMethod.GetString("lang_off"), ",")
	return method
}

func (p *Parser) parseServiceForProtoFile(protoFile *ProtoFile, st ServiceTag) (ret []*Service) {
	fdp := protoFile.fd.AsFileDescriptorProto()
	for i, protoService := range fdp.Service {
		name := protoService.GetName()
		service := &Service{
			sd:             protoFile.fd.GetServices()[i],
			Name:           name,
			DeprecatedName: xstrings.CamelCase(nameMustHaveSuffix(name, "Service")),
			DescName:       fmt.Sprintf("%s.%s", fdp.GetPackage(), name),
			DescProtoFile:  fdp.GetName(),
		}
		needActor := true
		needRPC := true
		if st == ServiceTagALL {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternServerHandler, name)
		} else if st == ServiceTagActor {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternActorClient, name)
			needRPC = false
		} else if st == ServiceTagRPC {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternRPCClient, name)
			needActor = false
		}
		service.ServerHandlerInterfaceName = fmt.Sprintf(p.cc.NamePatternServerHandler, name)
		service.RPCClientInterfaceName = fmt.Sprintf(p.cc.NamePatternRPCClient, name)
		service.ActorClientInterfaceName = fmt.Sprintf(p.cc.NamePatternActorClient, name)
		comment, ok := p.comments[protoService]
		if ok {
			service.Comment = comment.Content
		}
		an := GetAnnotation(comment, AnnotationService)
		// 整个service是否完全为actor方法
		serviceUriAutoAlias := an.GetBool("service_uri_auto_alias", false)
		// 整个service是否完全为actor方法
		isActorService := an.GetBool("actor", false)
		// 整个service是否完全为rpc方法
		isRPCService := an.GetBool("rpc", !isActorService)
		hasSpecifiedRPCService := an.Has("rpc")
		// 整个service是否完全为tell方法
		isServiceAllTell := an.GetBool("tell", false)

		service.LangOffTag = strings.Split(an.GetString("lang_off"), ",")
		for j, protoMethod := range protoService.Method {
			// actor参数，是否为actor是否为tell
			isAsk := true
			isTell := isServiceAllTell
			anMethod := GetAnnotation(p.comments[protoMethod], AnnotationService)
			isActorMethod := anMethod.GetBool("actor", isActorService)
			// 默认指定了actor方法则不再支持生成rpc逻辑，除非明确指定:
			// method级别的annotation指定生成RPC，service级别明确指定是rpc service
			isRPCMethod := anMethod.GetBool("rpc", !isActorMethod)
			if !isRPCMethod && hasSpecifiedRPCService && isRPCService {
				isRPCMethod = true
			}
			isTell = anMethod.GetBool("tell", isTell)
			if isTell {
				isAsk = false
			}
			if isActorMethod {
				if needActor {
					m := p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], true, isAsk, isRPCMethod, serviceUriAutoAlias)
					service.Methods = append(service.Methods, m)
					service.HasActorMethod = true
				}
			}
			if isRPCMethod {
				if needRPC {
					m := p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], false, isAsk, false, serviceUriAutoAlias)
					service.Methods = append(service.Methods, m)
				}
			}
		}
		if len(service.Methods) > 0 {
			ret = append(ret, service)
		}
	}
	return ret
}
func (p *Parser) typeStr(dotFullyQualifiedTypeName string) string {
	ss := strings.Split(dotFullyQualifiedTypeName, ".")
	if len(ss) == 1 {
		return dotFullyQualifiedTypeName
	}
	if protoFile, ok := p.dotFullyQualifiedTypeNameToProtoFile[dotFullyQualifiedTypeName]; ok {
		return strings.Join([]string{protoFile.GolangPackageName, ss[len(ss)-1]}, ".")
	}
	return strings.Join(ss[len(ss)-2:], ".")
}
