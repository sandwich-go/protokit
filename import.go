package protokit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sandwich-go/protokit/util"
)

// AddWithDotFullQualifiedName 返回pd在本次import中应该使用的struct名称
func (e *ImportSet) AddWithDotFullQualifiedName(dotFullyQualifiedName string, pf *ProtoFile) (string, *Import) {
	item := &Import{
		MessageDotFullQualifiedName: []string{dotFullyQualifiedName},
		ProtoFilePath:               pf.FilePath,
		GolangPackageName:           pf.GolangPackageName,
		GolangPackagePath:           pf.GolangPackagePath,
		PythonModuleName:            pf.GolangPackageName, // python使用golang的package name
		PythonModulePath:            pythonModulePath(pf.FilePath),
	}
	e.Add(item)
	// 根据item的名字调节使用的struct名称
	structName := GoStructNameWithGolangPackage(dotFullyQualifiedName, pf.Package, item.GolangPackageName)
	e.MessageDotFullQualifiedNameToGolangType[dotFullyQualifiedName] = structName
	return structName, item
}

func (e *ImportSet) Add(add *Import) {
	if util.ContainsString(e.ExcludeImportName, add.GolangPackageName) {
		return
	}
	duplicated := false
	// 本包内
	if add.GolangPackagePath == e.GolangPackagePath {
		// 上层需要根据Name重新定位package
		add.GolangPackageName = "."
		return
	}
	if add.GolangPackageName == e.GolangPackageName {
		// path不同但是package name相同，起别名
		n := e.importAliasMappingCount[add.GolangPackageName]
		n++
		add.GolangPackageName = fmt.Sprintf("%s%d", add.GolangPackageName, n)
		e.importAliasMappingCount[add.GolangPackageName] = n
	}
	originalName := add.GolangPackageName
	for _, i := range e.Set {
		if i.GolangPackagePath == add.GolangPackagePath {
			duplicated = true
			add.GolangPackageName = i.GolangPackageName
			i.MessageDotFullQualifiedName = util.StringSetAdd(i.MessageDotFullQualifiedName, add.MessageDotFullQualifiedName...)
			break
		}
		// path不同但是package name相同，起别名
		if i.GolangPackageName == add.GolangPackageName {
			add.GolangPackageName = fmt.Sprintf("%s%d", add.GolangPackageName, e.importAliasMappingCount[originalName])
		}
	}
	if !duplicated {
		e.Set = append(e.Set, add)
		e.importAliasMappingCount[originalName] += 1
	}

	sort.SliceStable(e.Set, func(i, j int) bool {
		return e.Set[i].GolangPackagePath < e.Set[j].GolangPackagePath
	})
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
