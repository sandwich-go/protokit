package protokit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sandwich-go/boost/xslice"
)

// AddWithDotFullQualifiedName 返回pd在本次import中应该使用的struct名称
func (e *ImportSet) AddWithDotFullQualifiedName(dotFullyQualifiedName string, pf *ProtoFile) (string, *Import) {
	item := &Import{
		MessageDotFullQualifiedName: []string{dotFullyQualifiedName},
		ProtoFilePath:               pf.FilePath,
		GolangPackageName:           pf.GolangPackageName,
		GolangPackagePath:           pf.GolangPackagePath,
		CSNamespace:                 pf.OptionCSNamespace,
		CSNamespaceName:             pf.GolangPackageName, // cs 使用golang的packagename作为别名
		PythonModuleName:            pf.GolangPackageName, // python使用golang的package name作为别名
		PythonModulePath:            pythonModulePath(pf.FilePath),
	}
	e.Add(item)
	// 根据item的名字调节使用的struct名称
	structName := GoStructNameWithGolangPackage(dotFullyQualifiedName, pf.Package, item.GolangPackageName)
	e.MessageDotFullQualifiedNameToGolangType[dotFullyQualifiedName] = structName
	return structName, item
}

func (e *ImportSet) Add(add *Import) {
	if xslice.ContainString(e.ExcludeImportName, add.GolangPackageName) {
		return
	}
	duplicated := false
	// 本包内
	if add.GolangPackagePath == e.GolangPackagePath {
		// 上层需要根据Name重新定位package
		add.GolangPackageName = "."
		return
	}
	for _, i := range e.Set {
		if i.GolangPackagePath == add.GolangPackagePath {
			duplicated = true
			add.GolangPackageName = i.GolangPackageName
			i.MessageDotFullQualifiedName = xslice.StringSetAdd(i.MessageDotFullQualifiedName, add.MessageDotFullQualifiedName...)
			break
		}
	}
	if duplicated {
		return
	}
	e.add(add)
	e.sort()
}

func (e *ImportSet) add(add *Import) {
	add.originGolangPackageName = add.GolangPackageName
	e.Set = append(e.Set, add)
}

func (e *ImportSet) sort() {
	sort.SliceStable(e.Set, func(i, j int) bool {
		return e.Set[i].GolangPackagePath < e.Set[j].GolangPackagePath
	})
	importAliasMappingCount := make(map[string]int)
	for _, i := range e.Set {
		originalName := i.originGolangPackageName
		if originalName == e.GolangPackageName {
			i.GolangPackageName = fmt.Sprintf("%s%d", originalName, importAliasMappingCount[originalName]+1)
		} else if _, ok := importAliasMappingCount[originalName]; ok {
			i.GolangPackageName = fmt.Sprintf("%s%d", originalName, importAliasMappingCount[originalName])
		}
		importAliasMappingCount[originalName]++
	}
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
func pythonPB2FilePath(fileName string) string {
	fileName = fileNameWithoutExtension(fileName)
	fp := strings.Replace(fmt.Sprintf(".%s_pb2", fileName), "/", ".", -1)
	fp = strings.TrimPrefix(fp, ".")
	return fp
}

func pythonModulePath(fileName string) string {
	fileName = pythonPB2FilePath(fileName)
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
