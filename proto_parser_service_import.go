package protokit

import (
	"fmt"
	"path"
	"sort"
	"strings"

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
	var ss = make([]string, 0, len(p.protoFilePathToProtoFile))
	for protoFilePath := range p.protoFilePathToProtoFile {
		ss = append(ss, protoFilePath)
	}
	sort.Strings(ss)
	for _, protoFilePath := range ss {
		protoFile := p.protoFilePathToProtoFile[protoFilePath]
		for _, sg := range protoFile.ServiceGroups {
			// 设定import忽略路径
			sg.ImportSet.ExcludeImportName = p.cc.ImportSetExclude
			for _, service := range sg.Services {
				for _, method := range service.Methods {
					// 名称需要做一次修复，根据import的package名称
					method.TypeInput, _ = p.addImportByDotFullyQualifiedTypeName(method.TypeInputDotFullQualifiedName, sg.ImportSet)
					method.TypeOutput, _ = p.addImportByDotFullyQualifiedTypeName(method.TypeOutputDotFullQualifiedName, sg.ImportSet)
					service.InputOutputTypes = xslice.StringsSetAdd(service.InputOutputTypes, method.TypeInput, method.TypeOutput)
					// 请求使用使用的uri名称, 需要用这个名字来作为http请求的路径，携带自由的package名称
					uriUsing := method.TypeInputWithSelfPackage
					if method.TypeInputAlias != "" {
						uriUsing = method.TypeInputAlias
					}
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
