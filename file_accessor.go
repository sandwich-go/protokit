package protokit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xpanic"
)

func ProtoFileContents(dirs ...string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	for _, pathProtoRoot := range dirs {
		pathProtoRoot = path.Clean(pathProtoRoot)
		fileList := make([]string, 0)
		err := xos.FilePathWalkFollowLink(pathProtoRoot, xos.FileWalkFuncWithExcludeFilter(&fileList, nil, ".proto"))
		if err != nil {
			return nil, err
		}
		for _, filePath := range fileList {
			content, err := xos.FileGetContents(filePath)
			if err != nil {
				return nil, err
			}
			relativePath := strings.TrimLeft(strings.Replace(filePath, pathProtoRoot, "", 1), string(os.PathSeparator))
			// proto import使用/, 兼容windows，替换路径中的\为/
			relativePath = strings.ReplaceAll(relativePath, `\`, `/`)
			if _, ok := files[relativePath]; ok {
				return nil, fmt.Errorf("got duplicate proto file with path:%s under:%s", relativePath, pathProtoRoot)
			}
			files[relativePath] = content
		}
	}
	return files, nil
}

// MustGetFileAccessorWithNamespace 获取文件加载器，会主动加载并缓存目录下的所有文件夹爱你内容
func MustGetFileAccessorWithNamespace(nsList ...*Namespace) FileAccessor {
	var dirList []string
	for _, v := range nsList {
		dirList = append(dirList, v.Path)
	}
	return MustGetFileAccessorWithDirs(dirList...)
}

// MustGetFileAccessorWithDirs 获取文件加载器，会主动加载并缓存目录下的所有文件夹爱你内容
func MustGetFileAccessorWithDirs(dirs ...string) FileAccessor {
	contents, err := ProtoFileContents(dirs...)
	xpanic.WhenErrorAsFmtFirst(err, "got error:%s while load contents with:%s ", strings.Join(dirs, ","))
	return GetFileAccessor(contents)
}

// GetFileAccessor 获取文件加载器
func GetFileAccessor(contents map[string][]byte) FileAccessor {
	return func(fielRelativePath string) (io.ReadCloser, error) {
		content, ok := contents[fielRelativePath]
		if !ok {
			return nil, errors.New("accessor can not read proto content with filename: " + fielRelativePath)
		}
		return ioutil.NopCloser(bytes.NewReader(content)), nil
	}
}
