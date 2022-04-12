package protokit

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Masterminds/sprig"
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xpanic"
)

func MustCallTemplate(name, templateStr string, data interface{}, fileName string, filer func([]byte) ([]byte, error)) []byte {
	var usingData interface{}
	usingData = data
	if usingData == nil {
		usingData = &MarkerInfo{}
	}
	if generated, ok := usingData.(Marker); ok {
		generated.Format()
	} else if generated, ok := usingData.(map[string]interface{}); ok {
		if generated == nil {
			generated = make(map[string]interface{})
		}
		usingData = generated
		tmp := MarkerInfo{}
		tmp.Format()
		generated["MarkerLeadingWithDoubleSlash"] = tmp.MarkerLeadingWithDoubleSlash
		generated["MarkerLeadingWithDoubleDash"] = tmp.MarkerLeadingWithDoubleDash
		generated["MarkerLeadingWithHexKey"] = tmp.MarkerLeadingWithHexKey
		generated["MarkerForHTML"] = tmp.MarkerForHTML
	} else {
		panic(fmt.Sprintf("CallTemplate got invalid data,type %v", reflect.TypeOf(data)))
	}
	t, err := template.New(name).Funcs(funcMap).Funcs(sprig.FuncMap()).Parse(templateStr)
	xpanic.PanicIfErrorAsFmtFirst(err, "got error:%w while parse template:%s ", name)
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, usingData)
	xpanic.PanicIfErrorAsFmtFirst(err, "got error:%w while Execute template:%s ", name)
	bytesUsing := buf.Bytes()
	if filer != nil {
		tmp, err := filer(buf.Bytes())
		xpanic.PanicIfErrorAsFmtFirst(err, "got error:%w while run user define filter,template:%s ", name)
		bytesUsing = tmp
	}
	if fileName != "" {
		dirName := filepath.Dir(fileName)
		err := os.MkdirAll(dirName, os.ModePerm)
		xpanic.PanicIfErrorAsFmtFirst(err, "got error:%w while MkdirAll:%s template:%s ", dirName, name)
		xos.MustFilePutContents(fileName, bytesUsing)
	}
	return bytesUsing
}
