package protokit

import (
	"path"
	"sort"

	"github.com/jhump/protoreflect/desc"
	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xslice"
)

type IsMapEntry interface {
	IsMapEntry() bool
}

func namespace(nsList []*Namespace, name string) *Namespace {
	for _, v := range nsList {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (p *Parser) parsePackage(nsList []*Namespace) {
	var upPackage = func(protoPackage *Package, dotFullyQualifiedTypeName string, protoFile *ProtoFile) {
		if protoPackage == nil {
			return
		}
		structName, item := p.addImportByDotFullyQualifiedTypeName(dotFullyQualifiedTypeName, protoPackage.ImportSet)
		if item != nil {
			p.pythonModule(dotFullyQualifiedTypeName, structName, item, protoPackage.ImportSet, protoPackage.IsGlobal, protoFile)
		}
	}

	// 处理所有注册进来的消息
	var ss = make([]string, 0, len(p.dotFullyQualifiedTypeNameToProtoFile))
	for dotFullyQualifiedTypeName := range p.dotFullyQualifiedTypeNameToProtoFile {
		ss = append(ss, dotFullyQualifiedTypeName)
	}
	sort.Strings(ss)
	for _, dotFullyQualifiedTypeName := range ss {
		protoFile := p.dotFullyQualifiedTypeNameToProtoFile[dotFullyQualifiedTypeName]
		tt := p.dotFullyQualifiedTypeNameToDescriptor[dotFullyQualifiedTypeName]
		if mapEntry, ok := tt.(IsMapEntry); ok && mapEntry.IsMapEntry() {
			// map entry不做处理
			continue
		}
		if _, ok := tt.(*desc.EnumDescriptor); ok {
			// enum不做处理
			continue
		}
		// 根据protoFile获取namespace
		ns := namespace(nsList, protoFile.Namespace)
		xpanic.WhenTrue(ns == nil, "can not got namspace with name: %s", protoFile.Namespace)

		golangPackagePath := protoFile.GolangPackagePath
		pi, ok := ns.Packages[golangPackagePath]
		if !ok {
			pi = NewPackageWithPackageName(protoFile.GolangPackageName, protoFile.GolangPackagePath)
			pi.FilePath, _ = path.Split(protoFile.FilePath)
			pi.Package = protoFile.Package
			pi.GolangRelative = p.cc.GolangRelative
			// 设定import忽略路径
			pi.ImportSet.ExcludeImportName = p.cc.ImportSetExclude
			ns.Packages[golangPackagePath] = pi
		}
		// 本保内的消息
		upPackage(pi, dotFullyQualifiedTypeName, protoFile)
		// 注册全局消息
		globalPackage := ns.Packages[NamespaceMessageRegistryPackageName]
		if globalPackage != nil {
			// 重新赋值一次，默认的message registry package在创建的时候没有这个参数
			globalPackage.GolangRelative = p.cc.GolangRelative
			// 设定import忽略路径
			globalPackage.ImportSet.ExcludeImportName = p.cc.ImportSetExclude
			upPackage(globalPackage, dotFullyQualifiedTypeName, protoFile)
		}
	}

	for _, protoFile := range p.protoFilePathToProtoFile {
		sg := protoFile.ServiceGroups[ServiceTagALL]
		if sg == nil || len(sg.Services) == 0 {
			continue
		}
		ns := namespace(nsList, protoFile.Namespace)
		xpanic.WhenTrue(ns == nil, "can not got namspace with name: %s", protoFile.Namespace)
		pp := ns.Packages[NamespaceMessageRegistryPackageName]
		if pp == nil {
			continue
		}
		messageDotFullQualifiedNameToGolangType := p.getMessageDotFullQualifiedNameToGolangType(nsList, ns)
		for _, service := range sg.Services {
			for _, method := range service.Methods {
				// 找出类型在当前ImportSet下的类型名
				golangInputType := messageDotFullQualifiedNameToGolangType[method.TypeInputDotFullQualifiedName]
				golangOutputType := messageDotFullQualifiedNameToGolangType[method.TypeOutputDotFullQualifiedName]
				if golangInputType != "" && method.TypeInputAlias != "" {
					pp.AliasToGolangType[method.TypeInputAlias] = golangInputType
				}
				if golangInputType != "" && method.FullPathHttpBackOffice != "" {
					pp.AliasToGolangType[method.FullPathHttpBackOffice] = golangInputType
				}
				if method.IsActor {
					golangInputType = pp.ImportSet.MessageDotFullQualifiedNameToGolangType[method.TypeInputDotFullQualifiedName]
					golangOutputType = pp.ImportSet.MessageDotFullQualifiedNameToGolangType[method.TypeOutputDotFullQualifiedName]
					if golangInputType != "" {
						pp.ActorMessageGolangType = xslice.StringsSetAdd(pp.ActorMessageGolangType, golangInputType)
					}
					if golangOutputType != "" {
						pp.ActorMessageGolangType = xslice.StringsSetAdd(pp.ActorMessageGolangType, golangOutputType)
					}
				}
			}
		}
	}
}

func (p *Parser) getMessageDotFullQualifiedNameToGolangType(nsList []*Namespace, ns *Namespace) map[string]string {
	if ns == nil {
		return nil
	}
	pp := ns.Packages[NamespaceMessageRegistryPackageName]
	if pp == nil {
		return nil
	}
	switch ns.Name {
	case NamespaceGoogle, NamespaceNetutils:
		return pp.ImportSet.MessageDotFullQualifiedNameToGolangType
	}
	out := make(map[string]string)
	for k, v := range pp.ImportSet.MessageDotFullQualifiedNameToGolangType {
		out[k] = v
	}
	for _, name := range []string{NamespaceGoogle, NamespaceNetutils} {
		for k, v := range p.getMessageDotFullQualifiedNameToGolangType(nsList, namespace(nsList, name)) {
			out[k] = v
		}
	}
	return out
}
