package protokit

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sandwich-go/boost/xstrings"
)

func nameMustHaveSuffix(s string, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}
func (p *Parser) parseService() {
	reqMap := make(map[string]map[string]string)
	for _, serviceTag := range allServiceTags {
		reqMap[serviceTag] = make(map[string]string)
	}
	for _, protoFile := range p.protoFilePathToProtoFile {
		for _, serviceTag := range allServiceTags {
			// 请求的uri校验应该在整个proto包级别 不应该在独立的文件内
			protoFile.ServiceGroups[serviceTag] = &ServiceGroup{
				ProtoFilePath: protoFile.FilePath,
				Services:      p.parseServiceForProtoFile(protoFile, serviceTag, reqMap[serviceTag]),
				ImportSet:     NewImportSet(protoFile.GolangPackageName, protoFile.GolangPackagePath),
			}
		}
	}
}

func (p *Parser) parseServiceForProtoFile(protoFile *ProtoFile, st ServiceTag, reqMap map[string]string) (ret []*Service) {
	fdp := protoFile.fd.AsFileDescriptorProto()
	for i, protoService := range fdp.Service {
		name := protoService.GetName()
		service := &Service{
			Parser:         p,
			sd:             protoFile.fd.GetServices()[i],
			Name:           name,
			DeprecatedName: xstrings.CamelCase(nameMustHaveSuffix(name, "Service")),
			DescName:       fmt.Sprintf("%s.%s", fdp.GetPackage(), name),
			DescProtoFile:  fdp.GetName(),
		}
		service.RpcOption = getRpcServiceOption(service.sd)
		service.BackOfficeOption = getBackOfficeServiceOption(service.sd)
		service.IsJob = isJobService(service.sd)
		needActor := true
		needRPC := true
		needERPC := true
		needJob := true
		if st == ServiceTagALL {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternServerHandler, name)
			if service.IsJob {
				needJob = false
			}
		} else if st == ServiceTagActor {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternActorClient, name)
			needRPC = false
			needERPC = false
			needJob = false
		} else if st == ServiceTagRPC {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternRPCClient, name)
			needActor = false
			needERPC = false
			needJob = false
		} else if st == ServiceTagERPC {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternERPCClient, name)
			needActor = false
			needRPC = false
			needJob = false
		} else if st == ServiceTagJob {
			service.ServiceName = fmt.Sprintf(p.cc.NamePatternERPCClient, name)
			needActor = false
			needRPC = false
			needERPC = false
			needJob = true
		}
		service.ServerHandlerInterfaceName = fmt.Sprintf(p.cc.NamePatternServerHandler, name)
		service.RPCClientInterfaceName = fmt.Sprintf(p.cc.NamePatternRPCClient, name)
		service.ActorClientInterfaceName = fmt.Sprintf(p.cc.NamePatternActorClient, name)
		service.ERPCClientInterfaceName = fmt.Sprintf(p.cc.NamePatternERPCClient, name)
		service.JobClientInterfaceName = fmt.Sprintf(p.cc.NamePatternJobClient, name)
		service.JobServiceInterfaceName = fmt.Sprintf(p.cc.NamePatternJobService, name)
		comment, ok := p.comments[protoService]
		if ok {
			service.Comment = comment.Content
		}
		var an serviceAnnotation
		if service.RpcOption != nil {
			an = &serviceOptionAnnotation{service.RpcOption}
		} else {
			an = GetAnnotation(comment, AnnotationService)
		}
		snakeCase, _ := an.Bool(QueryPathSnakeCase, true)

		service.QueryPath = standardQueryPath(an.String(QueryPath, p.cc.DefaultQueryPath), snakeCase, p.cc.QueryPathMapping)

		for _, v := range p.cc.GetInvalidServiceAnnotations() {
			if an.Contains(strings.TrimSpace(v)) {
				log.Fatal().Msg(fmt.Sprintf("invalid annotation: %s", v))
			}
		}
		serviceUriAutoAlias, _ := an.Bool(ServiceUriAutoAlias, false)
		// 整个service是否完全为actor方法
		isActorService, _ := an.Bool(ServiceTagActor, false)
		// 整个service是否完全为erpc方法
		isERPCService, _ := an.Bool(ServiceTagERPC, false)
		// 整个service是否完全为rpc方法
		isRPCService, _ := an.Bool(ServiceTagRPC, false)
		// 整个service是否完全为tell方法
		isServiceAllTell, _ := an.Bool(Tell, false)

		var anMethod methodeAnnotation
		if !isERPCService && !isActorService && !isRPCService && !service.IsJob {
			// service级别没有任何定义，则如果任意一个方法既不是actor也不是erpc那么这个service就是rpc
			// 否则这个就不是rpc的service（没有任何一个方法是rpc）
			for _, protoMethod := range protoService.Method {
				rpcMethodOption := getRpcMethodOption(protoMethod)
				if rpcMethodOption != nil {
					anMethod = &methodOptionAnnotation{rpcMethodOption}
				} else {
					anMethod = GetAnnotation(p.comments[protoMethod], AnnotationService)
				}
				if isRpcMethod, _ := anMethod.Bool(ServiceTagRPC, false); isRpcMethod {
					// 任意一个方法是rpc，那么这个service就是rpc service
					isRPCService = true
					break
				}
				isActorMethod, _ := anMethod.Bool(ServiceTagActor, false)
				isErpcMethod, _ := anMethod.Bool(ServiceTagERPC, false)
				if !isActorMethod && !isErpcMethod {
					// 任意一个方法不是actor也不是erpc，那么这个service就是rpc service
					isRPCService = true
					break
				}
			}
		}

		service.LangOffTag = strings.Split(an.String(LangOff), ",")
		for j, protoMethod := range protoService.Method {
			// actor参数，是否为actor是否为tell
			isAsk := true
			isTell := isServiceAllTell
			rpcMethodOption := getRpcMethodOption(protoMethod)
			if rpcMethodOption != nil {
				anMethod = &methodOptionAnnotation{rpcMethodOption}
			} else {
				anMethod = GetAnnotation(p.comments[protoMethod], AnnotationService)
			}
			isActorMethod, _ := anMethod.Bool(ServiceTagActor, isActorService)
			isERPCMethod, _ := anMethod.Bool(ServiceTagERPC, isERPCService)
			isRPCMethod, _ := anMethod.Bool(ServiceTagRPC, isRPCService)
			isTell, _ = anMethod.Bool(Tell, isTell)
			if isTell {
				isAsk = false
			}
			var m *Method
			if service.IsJob {
				jobMethodOption := getJobMethodOption(protoMethod)
				if jobMethodOption != nil && jobMethodOption.Creator != nil {
					if needJob {
						m = p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], false, false, false, serviceUriAutoAlias, false, service.QueryPath, true)
						service.HasJobCreatorMethod = true
						service.Methods = append(service.Methods, m)
					}
				}
			}
			if isActorMethod {
				if needActor {
					m = p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], true, isAsk, isActorMethod, serviceUriAutoAlias, isERPCMethod, service.QueryPath, false)
					service.Methods = append(service.Methods, m)
					service.HasActorMethod = true
				}
			}
			if isERPCMethod {
				if needERPC {
					m = p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], isActorMethod, isAsk, isRPCMethod, serviceUriAutoAlias, isERPCMethod, service.QueryPath, false)
					service.Methods = append(service.Methods, m)
					service.HasERPCMethod = true
				}
			}
			if isRPCMethod {
				if needRPC {
					m = p.method(protoFile, service.Name, protoMethod, protoFile.fd.GetServices()[i].GetMethods()[j], false, isAsk, false, serviceUriAutoAlias, false, service.QueryPath, false)
					service.Methods = append(service.Methods, m)
				}
			}
			if m != nil {
				checkName := m.TypeInputDotFullQualifiedName
				if m.TypeInputAlias != "" {
					checkName = m.TypeInputAlias
				}
				// 校验uriUsing是否已经被使用过
				// 如果为严格模式，才会去校验
				if v, ok0 := reqMap[checkName]; ok0 && p.cc.StrictMode {
					log.Fatal().
						Str("req", m.TypeInputDotFullQualifiedName).
						Str("method_now", m.TypeInputGRPC).
						Str("method_last", v).
						Str("uri", checkName).
						Msg("duplicated request uri")
				}
				reqMap[checkName] = m.TypeInputGRPC
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
