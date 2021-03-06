package protokit

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sandwich-go/boost/xpanic"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Parameter protokito传递过来的参数数据
type Parameter struct {
	NameSpaces map[string]string // namespace到根路径的映射,需要手动加载
	WorkingDir string
	// 输出目录
	Outpath       string
	OutpathLua    string
	OutpathGolang string
	OutpathJS     string
	OutpathCSharp string
	OutpathPython string
	// raw data属性
	RawDataPackageName   string
	OutpathRawdataServer string
	OutpathRawdataClient string
	// golang 基础配置
	GolangBasePackagePath string
	GolangRelative        bool
	// cs基础配置
	CSBaseNamespace    string
	CSConfNamespace    string
	CSRawDataNamespace string
	// 日志配置
	LogLevel int
	LogColor bool
	Raw      []byte
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
	xpanic.WhenErrorAsFmtFirst(err, "got err:%w while read stdin")
	err = json.Unmarshal(content, p.Parameter)
	xpanic.WhenErrorAsFmtFirst(err, "got err:%w while unmarshal to Parameter")
	xpanic.WhenTrue(p.Parameter.WorkingDir == "", "WorkingDir should not empty")

	zerolog.SetGlobalLevel(zerolog.Level(p.Parameter.LogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339, NoColor: !p.Parameter.LogColor})

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
