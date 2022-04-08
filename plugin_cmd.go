package protokit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sandwich-go/boost/xexec"
	"github.com/sandwich-go/boost/xpanic"
)

func MustRun(pluginPath string, args string, parameter *Parameter) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	bb, err := json.Marshal(parameter)
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while marshal parameter")
	content, err := xexec.Run(pluginPath+" "+args, filepath.Dir(ex), bytes.NewBuffer(bb))
	xpanic.PanicIfErrorAsFmtFirst(err, "got err:%w while run plugin :%s with args:%s ", pluginPath, args)
	fmt.Println(content)
}

var _ = MustRun
