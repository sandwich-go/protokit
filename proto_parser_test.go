package protokit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("parse all proto files", t, func() {
		var nsList []*Namespace
		nsList = append(nsList, NewNamespace(NamespaceGoogle, "/Users/wh/prjs/funplus/protokitgo/sdk/proto_google"))
		nsList = append(nsList, NewNamespace(NamespaceNetutils, "/Users/wh/prjs/funplus/protokitgo/sdk/proto_netutils_queue"))
		nsList = append(nsList, NewNamespace(NamespaceUser, "/Users/wh/prjs/funplus/protokitgo/example/protos"))
		m := NewParser(WithProtoFileAccessor(MustGetFileAccessorWithNamespace(nsList...)), WithGolangBasePackagePath("example/gen/golang"))
		m.Parse(nsList...)
	})
}
