package protokit

import (
	"fmt"
	"strings"

	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xstrings"
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
		isActorService := an.GetBool("actor", false)
		// URI使用是否GRPC模式
		methodAllAliasAllAsGRPC := p.cc.URIUsingGRPC
		if !methodAllAliasAllAsGRPC {
			methodAllAlias := an.GetString("alias")
			xpanic.WhenTrue(methodAllAlias != "" && methodAllAlias != "grpc", "service annotation alias only support grpc now, got:%s", methodAllAlias)
			methodAllAliasAllAsGRPC = methodAllAlias == "grpc"
		}

		isActorServiceAllTell := an.GetBool("tell", false)
		service.LangOffTag = strings.Split(an.GetString("lang_off"), ",")
		for j, protoMethod := range protoService.Method {
			// actor参数，是否为actor是否为tell
			isActorMethod := isActorService
			isAsk := true
			isTell := isActorServiceAllTell
			anMethod := GetAnnotation(p.comments[protoMethod], AnnotationService)
			isActorMethod = anMethod.GetBool("actor", isActorMethod)
			isTell = anMethod.GetBool("tell", isTell)
			if isTell {
				isAsk = false
			}
			// Note:
			// 这里只是简单的换算一次格式合法的名称，具体请求名要通过ImportSet进行纠正
			reqTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetInputType()), ".")
			rspTypeName := strings.TrimPrefix(p.typeStr(protoMethod.GetOutputType()), ".")

			method := &Method{
				md:                             protoFile.fd.GetServices()[i].GetMethods()[j],
				Name:                           xstrings.CamelCase(protoMethod.GetName()),
				TypeInputDotFullQualifiedName:  protoMethod.GetInputType(),
				TypeOutputDotFullQualifiedName: protoMethod.GetOutputType(),
				TypeInputWithSelfPackage:       reqTypeName,
				TypeOutputWithSelfPackage:      rspTypeName,
				IsActor:                        isActorMethod,
				IsAsk:                          isAsk,
				IsTell:                         isTell,
			}
			// 只保留rpc逻辑
			if isActorMethod && !needActor {
				continue
			}
			// 只保留Acor逻辑
			if !isActorMethod && !needRPC {
				continue
			}
			if methodComment, exist := p.comments[protoMethod]; exist && methodComment != nil {
				method.Comment = methodComment.Content
			}
			method.TypeInputGRPC = fmt.Sprintf("/%s.%s/%s", fdp.GetPackage(), service.Name, method.Name)
			// 请求别名逻辑，允许proto中设定input类型别名，在请求的proto中uri将使用此名称
			nameAlias := ""
			if methodAllAliasAllAsGRPC {
				nameAlias = method.TypeInputGRPC
			}
			if anMethod.Has("alias") {
				nameAlias = fmt.Sprintf("%s_%s", service.Name, method.Name) // 默认alias
				if autoAlias := anMethod.GetBool("alias", false); !autoAlias {
					// proto中指定了alias名称
					nameAlias = anMethod.GetString("alias")
				}
				if strings.EqualFold(nameAlias, "grpc") {
					// 如果指定为grpc，则使用grpc的路由名称
					nameAlias = method.TypeInputGRPC
				} else {
					// name alias 必须有namespace前缀，以便于激活自动转发功能，如果没有指定，则使用与TypeInput想听的前缀
					if !strings.Contains(nameAlias, ".") {
						nameAlias = fmt.Sprintf("%s.%s", strings.Split(method.TypeInputWithSelfPackage, ".")[0], nameAlias)
					}
				}
			}
			// 允许逻辑层强制指定别名，此时不再进行namepace的添加逻辑
			if anMethod.Has("alias_force") {
				nameAlias = anMethod.GetString("alias_force")
				if strings.EqualFold(nameAlias, "grpc") {
					// 如果指定为grpc，则使用grpc的路由名称
					nameAlias = method.TypeInputGRPC
				}
			}
			method.TypeInputAlias = strings.TrimSpace(nameAlias)
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
				method.HTTPPathComment = "from proto, user defined"
			}

			method.LangOffTag = strings.Split(anMethod.GetString("lang_off"), ",")

			service.Methods = append(service.Methods, method)
			if method.IsActor {
				service.HasActorMethod = true
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
