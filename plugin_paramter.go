package protokit

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sandwich-go/boost/xpanic"
)

// Parameter protokito传递过来的参数数据
type Parameter struct {
	NameSpaces           map[string]string // namespace到根路径的映射,需要手动加载
	WorkingDir           string
	Outpath              string
	OutpathLua           string
	OutpathGolang        string
	OutpathJS            string
	OutpathCSharp        string
	OutpathPython        string
	OutpathRawdataServer string
	OutpathRawdataClient string

	GolangBasePackagePath    string
	GolangRelative           bool
	GolangRawDataPackageName string
	CSBaseNamespace          string
	CSConfNamespace          string
	CSRawDataNamespace       string
}

type Plugin struct {
	Parameter  *Parameter // 自动解析得到的参数数据，由protokitgo传递而来
	Namespaces map[string]*Namespace
}

// NewPlugin 返回插件
func MustNewPlugin(opts ...Option) *Plugin {
	var p Plugin
	p.Parameter = &Parameter{}
	p.Namespaces = make(map[string]*Namespace)
	content, err := ioutil.ReadAll(os.Stdin)
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while read stdin")
	err = json.Unmarshal(content, p.Parameter)
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while unmarshal to Parameter")
	xpanic.PanicIfTrue(p.Parameter.WorkingDir == "", "WorkingDir should not empty")
	// 解析namespaces
	var nsList []*Namespace
	for k, v := range p.Parameter.NameSpaces {
		nsList = append(nsList, NewNamespace(k, v))
	}
	if len(nsList) > 0 {
		optList := []Option{
			WithProtoFileAccessor(MustGetFileAccessorWithNamespace(nsList...)),
			WithGolangBasePackagePath(p.Parameter.GolangBasePackagePath),
			WithGolangRelative(p.Parameter.GolangRelative),
		}
		parser := NewParser(append(optList, opts...)...)
		parser.Parse(nsList...)
	}
	for _, v := range nsList {
		p.Namespaces[v.Name] = v
	}
	return &p
}
