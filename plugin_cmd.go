package protokit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sandwich-go/boost/xexec"
	"github.com/sandwich-go/boost/xpanic"
)

// MustRun pluginPath目前只支持本地文件, 后续加入远程版本支持
func MustRun(pluginPath string, args string, parameter *Parameter) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pluginPath, err = getCommand(pluginPath)
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while check plugin")
	bb, err := json.Marshal(parameter)
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while marshal parameter")
	content, err := xexec.Run(pluginPath+" "+args, filepath.Dir(ex), bytes.NewBuffer(bb))
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while run plugin :%s with args:%s ", pluginPath, args)
	fmt.Println(content)
}

var _ = MustRun

func getCommand(pluginPath string) (string, error) {
	p, err := exec.LookPath(pluginPath)
	if err == nil {
		abs, err := filepath.Abs(p)
		if err != nil {
			return abs, err
		}
		return abs, nil
	}
	defaultErr := errors.New("invalid plugin value " + pluginPath)
	return pluginPath, defaultErr
}
