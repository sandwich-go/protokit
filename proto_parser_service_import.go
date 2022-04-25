package protokit

import (
	"fmt"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sandwich-go/boost/xslice"
)

func (p *Parser) addImportByDotFullyQualifiedTypeName(dotFullyQualifiedTypeName string, set *ImportSet) (string, *Import) {
	if dotFullyQualifiedTypeName == ".google.protobuf.Empty" {
		return "protobufEmpty.Empty", nil
	}
	protoFile, ok := p.dotFullyQualifiedTypeNameToProtoFile[dotFullyQualifiedTypeName]
	if !ok {
		return dotFullyQualifiedTypeName, nil
	}
	structName, item := set.AddWithDotFullQualifiedName(dotFullyQualifiedTypeName, protoFile)
	return structName, item
}

func (p *Parser) parseImport() {
	// fixme 校验req rsp映射关系,TCP需要严格校验，HTTP缺可以不严格校验
	for _, protoFile := range p.protoFilePathToProtoFile {
		for _, sg := range protoFile.ServiceGroups {
			reqMap := make(map[string]string)
			// 设定import忽略路径
			sg.ImportSet.ExcludeImportName = p.cc.ImportSetExclude
			for _, service := range sg.Services {
				for _, method := range service.Methods {
					// 名称需要做一次修复，根据import的package名称
					method.TypeInput, _ = p.addImportByDotFullyQualifiedTypeName(method.TypeInputDotFullQualifiedName, sg.ImportSet)
					method.TypeOutput, _ = p.addImportByDotFullyQualifiedTypeName(method.TypeOutputDotFullQualifiedName, sg.ImportSet)
					service.InputOutputTypes = xslice.StringSetAdd(service.InputOutputTypes, method.TypeInput, method.TypeOutput)
					// 请求使用使用的uri名称
					uriUsing := method.TypeInput
					if method.TypeInputAlias != "" {
						uriUsing = method.TypeInputAlias
					}
					// 校验uriUsing是否已经被使用过
					if v, ok := reqMap[uriUsing]; ok {
						log.Fatal().
							Str("req", method.TypeInput).
							Str("method_now", method.TypeInputGRPC).
							Str("method_last", v).
							Msg("duplicated request uri")
					}
					reqMap[uriUsing] = method.TypeInputGRPC
					// http请求path逻辑校验，需要依赖纠正过后的TypeInput
					if method.HTTPPath == "" {
						method.HTTPPathComment = "auto generate by ProtoKitGo"
						pathUsing := strings.TrimLeft(uriUsing, ".")
						if !strings.HasPrefix(pathUsing, "/") {
							pathUsing = "/" + pathUsing
						}
						method.HTTPPath = path.Clean(fmt.Sprintf(p.cc.NamePatternHTTPPath, pathUsing))
					}
				}
			}
		}
	}
}
