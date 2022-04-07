package protokit

import (
	"sort"
	"strings"
)

type PythonNestedMessage struct {
	PythonMessage string
	GolangMessage string
}
type PythonModule struct {
	From                    string
	ModuleName              string
	FullModuleName          string
	MessageToRegister       []string
	NestedMessageToRegister []*PythonNestedMessage
}

func PythonStructNameWithGolangPackage(fullyQualifiedName string, protoPackagePath, goPackageName string) string {
	protoPackageWithDot := strings.ReplaceAll(protoPackagePath, "/", ".")
	fullyQualifiedName = strings.TrimPrefix(fullyQualifiedName, ".")
	nameWithoutProtoPackage := strings.TrimPrefix(fullyQualifiedName, protoPackageWithDot)
	structName := strings.TrimPrefix(nameWithoutProtoPackage, ".")
	structName = strings.TrimPrefix(structName, "/")
	if goPackageName == "." || goPackageName == "" {
		return strings.TrimPrefix(structName, ".")
	}
	return goPackageName + "." + strings.TrimPrefix(structName, ".")
}

func pythonPB2FileName(fileName string) string {
	fileName = pythonPB2FilePath(fileName)
	return nameWithoutPacket(fileName)
}

func nameWithoutPacket(name string) string {
	if pos := strings.LastIndexByte(name, '.'); pos != -1 {
		return name[pos+1:]
	}
	return name
}

func (p *Parser) pythonModule(dotFullyQualifiedTypeName string, structName string, item *Import, set *ImportSet, isGlobal bool, protoFile *ProtoFile) {
	pythonStructName := PythonStructNameWithGolangPackage(dotFullyQualifiedTypeName, protoFile.Package, item.GolangPackageName)
	fullModuleName := pythonPB2FileName(protoFile.FilePath)
	if isGlobal {
		// 全路径，包含文件名
		fullModuleName = pythonModulePath(protoFile.FilePath)
	}

	var module *PythonModule
	for _, v := range set.PythonModules {
		if v.FullModuleName == fullModuleName {
			module = v
		}
	}
	if module == nil {
		module = &PythonModule{FullModuleName: fullModuleName}
		if isGlobal {
			moduleNameSlice := strings.Split(fullModuleName, ".")
			if len(moduleNameSlice) > 0 {
				module.From = strings.Join(moduleNameSlice[:len(moduleNameSlice)-1], ".")
				module.ModuleName = moduleNameSlice[len(moduleNameSlice)-1]
			}
		}
		set.PythonModules = append(set.PythonModules, module)
	}
	if !strings.Contains(pythonStructName, ".") {
		module.MessageToRegister = append(module.MessageToRegister, pythonStructName)
	} else {
		module.NestedMessageToRegister = append(module.NestedMessageToRegister, &PythonNestedMessage{
			PythonMessage: pythonStructName,
			GolangMessage: structName,
		})

		sort.Slice(module.NestedMessageToRegister, func(i, j int) bool {
			return module.NestedMessageToRegister[i].GolangMessage < module.NestedMessageToRegister[j].GolangMessage
		})
	}
	sort.Strings(module.MessageToRegister)
	sort.SliceStable(set.PythonModules, func(i, j int) bool {
		return set.PythonModules[i].FullModuleName < set.PythonModules[j].FullModuleName
	})
}
