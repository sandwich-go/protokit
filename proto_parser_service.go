package protokit

import (
	"fmt"
	"strings"

	"github.com/sandwich-go/protokit/util"
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
	for _, protoService := range fdp.Service {
		name := protoService.GetName()
		service := &Service{
			Name:           name,
			DeprecatedName: util.CamelCase(nameMustHaveSuffix(name, "Service")),
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
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternRpcClient, name)
			needActor = false
		}
		service.ServerHandlerInterfaceName = fmt.Sprintf(p.cc.NamePatternServerHandler, name)
		service.RPCClientInterfaceName = fmt.Sprintf(p.cc.NamePatternRpcClient, name)
		service.ActorClientInterfaceName = fmt.Sprintf(p.cc.NamePatternActorClient, name)
		comment, ok := p.comments[protoService]
		if ok {
			service.Comment = comment.Content
		}
		an := GetAnnotation(comment, AnnotationService)
		isActorService := an.GetBool("actor", false)
		isActorServiceAllTell := an.GetBool("tell", false)

		for _, protoMethod := range protoService.Method {
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
				Name:                           util.CamelCase(protoMethod.GetName()),
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
			if anMethod.Has("alias") {
				nameAlias := fmt.Sprintf("%s_%s", service.Name, method.Name) // 默认alias
				if autoAlias := anMethod.GetBool("alias", false); !autoAlias {
					// proto中指定了alias名称
					nameAlias = anMethod.GetString("alias")
				}
				if !strings.Contains(nameAlias, ".") {
					nameAlias = fmt.Sprintf("%s.%s", strings.Split(method.TypeInput, ".")[0], nameAlias)
				}
				method.TypeInputAlias = strings.TrimSpace(nameAlias)
			}
			// 默认的http请求路径
			if pathStr, err := HTTPPath(protoMethod); err == nil && pathStr != "" {
				if !strings.HasPrefix(pathStr, "/") {
					pathStr = "/" + pathStr
				}
				method.HTTPPath = pathStr
				method.HTTPPathComment = "from proto, user defined"
			}

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
	if profoFile, ok := p.dotFullyQualifiedTypeNameToProtoFile[dotFullyQualifiedTypeName]; ok {
		return strings.Join([]string{profoFile.GolangPackageName, ss[len(ss)-1]}, ".")
	}
	return strings.Join(ss[len(ss)-2:], ".")
}
