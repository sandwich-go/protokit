package protokit

import (
	"path"

	"github.com/jhump/protoreflect/desc"
	"github.com/sandwich-go/protokit/util"
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
		structName, item := p.addImportByDotFullyQualifiedTypeName(dotFullyQualifiedTypeName, protoPackage.ImportSet, protoFile)
		if item != nil {
			p.pythonModule(dotFullyQualifiedTypeName, structName, item, protoPackage.ImportSet, protoPackage.IsGlobal, protoFile)
		}
	}

	// 处理所有注册进来的消息
	for dotFullyQualifiedTypeName, protoFile := range p.dotFullyQualifiedTypeNameToProtoFile {
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
		util.PanicIfTrue(ns == nil, "can not got namspace with name: %s", protoFile.Namespace)

		golangPackagePath := protoFile.GolangPackagePath
		pi, ok := ns.Packages[golangPackagePath]
		if !ok {
			pi = NewPackageWithPackageName(protoFile.GolangPackageName, protoFile.GolangPackagePath, p.cc)
			pi.FilePath, _ = path.Split(protoFile.FilePath)
			pi.Package = protoFile.Package
			pi.GolangRelative = p.cc.GolangRelative
			ns.Packages[golangPackagePath] = pi
		}
		// 本保内的消息
		upPackage(pi, dotFullyQualifiedTypeName, protoFile)
		// 注册全局消息
		globalPackage := ns.Packages[NamespaceMessageRegistryPackageName]
		if globalPackage != nil {
			// 重新赋值一次，默认的message registry package在创建的时候没有这个参数
			globalPackage.GolangRelative = p.cc.GolangRelative
			upPackage(globalPackage, dotFullyQualifiedTypeName, protoFile)
		}
	}

	for _, protoFile := range p.protoFilePathToProtoFile {
		sg := protoFile.ServiceGroups[ServiceTagALL]
		if sg == nil || len(sg.Services) == 0 {
			continue
		}
		ns := namespace(nsList, protoFile.Namespace)
		util.PanicIfTrue(ns == nil, "can not got namspace with name: %s", protoFile.Namespace)
		pp := ns.Packages[NamespaceMessageRegistryPackageName]
		if pp == nil {
			continue
		}
		for _, service := range sg.Services {
			for _, method := range service.Methods {
				// 找出类型在当前ImportSet下的类型名
				golangInputType := pp.ImportSet.MessageDotFullQualifiedNameToGolangType[method.TypeInputDotFullQualifiedName]
				golangOutputType := pp.ImportSet.MessageDotFullQualifiedNameToGolangType[method.TypeOutputDotFullQualifiedName]
				if golangInputType != "" && method.TypeInputAlias != "" {
					pp.AliasToGolangType[method.TypeInputAlias] = golangInputType
				}
				if golangInputType != "" && method.IsActor {
					pp.ActorMessageGolangType = util.StringSetAdd(pp.ActorMessageGolangType, golangInputType)
				}
				if golangOutputType != "" && method.IsActor {
					pp.ActorMessageGolangType = util.StringSetAdd(pp.ActorMessageGolangType, golangOutputType)
				}
			}
		}
	}
}
