package protokit

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// 获取fd这个文件的golang package path，relative模式下有特殊处理
func GolangPackagePathAndName(fd *desc.FileDescriptor, basePackagePath string, golangRelative bool) (string, string) {
	packageName := ""
	packagePath := ""
	rootPath := basePackagePath
	fileBasePath := filepath.Dir(fd.GetName())

	// for windows filepath.dir return a\\b
	fileBasePath = strings.ReplaceAll(fileBasePath, "\\", "/")

	// golang package 优先proto中的golang package定义 => proto package => 依赖目录结构
	protoGolangPackage := fd.AsFileDescriptorProto().GetOptions().GetGoPackage()
	protoProtoPackageAsGolangPackage := strings.ReplaceAll(fd.AsFileDescriptorProto().GetPackage(), ".", "_")
	if protoGolangPackage == "" {
		protoGolangPackage = strings.ToLower(protoProtoPackageAsGolangPackage)
	}

	if protoGolangPackage == "" {
		// 以文件名作准
		protoGolangPackage = (strings.Split(fd.GetFile().GetName(), "."))[0]
	}

	packageSlice := strings.Split(protoGolangPackage, "/")
	packageName = packageSlice[len(packageSlice)-1]
	// protoGolangPackage的最后一个字段作为package名称,但是path是目录相关的
	if golangRelative && fileBasePath != "google/protobuf" {
		if fileBasePath == "" || fileBasePath == "." {
			// fixme FunPlus Puzzle Game使用问题兼容
			// fixme 如果获取到的proto在根目录下，relative模式下，package path不再完全依赖于file path，将package name纳package path
			fileBasePath = packageName
		}
		packagePath = path.Join(rootPath, fileBasePath)
		return packagePath, packageName
	}

	return protoGolangPackage, packageName
}

//(*OuterTestT1)(nil),             // 2: msg.Outer.test_t1
//(*OuterTest_T2)(nil),            // 3: msg.Outer.test_T2
//(*Outer_TestT3)(nil),            // 4: msg.Outer.Test_t3
//(*Outer_Test_T4)(nil),           // 5: msg.Outer.Test_T4
//(*Outer_Test_T5)(nil),           // 6: msg.Outer.Test__t5
//(*Outer_Test__T6)(nil),          // 7: msg.Outer.Test__T6
//(*Outer_TestT7__)(nil),          // 8: msg.Outer.Test_t7__
//(*Outer_TestT8__)(nil),          // 9: msg.Outer.Test_t8__
//(*Outer_Test____T9__)(nil),      // 10: msg.Outer.Test_____t9__

// goPackageName通过方法GetGolangPackageName获取,如果传入.或者空，则返回struct名称不带package名字
func GoStructNameWithGolangPackage(fullyQualifiedName string, protoPackagePath, goPackageName string) string {
	protoPackageWithDot := strings.ReplaceAll(protoPackagePath, "/", ".")
	fullyQualifiedName = strings.TrimPrefix(fullyQualifiedName, ".")
	nameWithoutProtoPackage := strings.TrimPrefix(fullyQualifiedName, protoPackageWithDot)
	structName := GoStructNameFromFullyQualifiedNameTrimProtoPackage(nameWithoutProtoPackage)
	structName = strings.TrimPrefix(structName, "/")
	if goPackageName == "." || goPackageName == "" {
		return strings.TrimPrefix(structName, ".")
	}
	return goPackageName + "." + strings.TrimPrefix(structName, ".")
}

var keywords = map[string]struct{}{
	"Reset":        {},
	"String":       {},
	"ProtoMessage": {},
	"ProtoReflect": {},
	"Descriptor":   {},
}

func GoFieldName(s string) string {
	s = strings.Title(s)
	if _, ok := keywords[s]; ok {
		return s + "_"
	}
	return goFieldName(s)
}

func goFieldName(s string) string {
	var ns string
	var skip, toUpper bool
	for j, c := range s {
		if skip {
			skip = false
			continue
		}
		if toUpper {
			ns += strings.ToUpper(string(c))
			toUpper = c == '_'
			continue
		}
		if c == '_' {
			if j == 0 {
				ns += string('X')
				toUpper = true
			} else if j < len(s)-1 && ('a' <= s[j+1] && s[j+1] <= 'z') {
				toUpper = true
			} else {
				ns += string(c)
			}
		} else if '0' <= s[j] && s[j] <= '9' {
			ns += string(c)
			toUpper = true
			if j < len(s)-1 && '_' == s[j+1] {
				skip = true
			}
		} else {
			ns += string(c)
		}
	}
	return ns
}

func goName(s string) string {
	var ns string
	var toUpper bool
	for j, c := range s {
		if toUpper {
			ns += strings.ToUpper(string(c))
			toUpper = false
			continue
		}
		if c == '_' {
			if j == 0 {
				ns += string('X')
			} else if j < len(s)-1 && ('a' <= s[j+1] && s[j+1] <= 'z') {
				toUpper = true
			} else {
				ns += string(c)
			}
		} else if '0' <= s[j] && s[j] <= '9' {
			ns += string(c)
			toUpper = true
		} else {
			ns += string(c)
		}
	}
	return ns
}

// 由fullyQualifiedName 转换到 golang struct名称，fullyQualifiedName需要去除掉proto package的名称
// 底层无法自动判定proto package名称依赖上层传递
// FunPlus.ServerCommon.Config.ActivityData : proto package名称为FunPlus.ServerCommon.Config
// msg.Outer.Test_____t9__ : proto package名称为msg
func GoStructNameFromFullyQualifiedNameTrimProtoPackage(fullyQualifiedNameWithoutProtoPackage string) string {
	fullyQualifiedNameWithoutProtoPackage = strings.TrimPrefix(fullyQualifiedNameWithoutProtoPackage, ".")
	nameParts := strings.Split(fullyQualifiedNameWithoutProtoPackage, ".")
	ret := ""
	for i, s := range nameParts {
		if 'A' <= s[0] && s[0] <= 'Z' && i != 0 {
			ret += "_"
		}
		ret += goName(strings.Title(s))
	}
	return ret
}
