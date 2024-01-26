package protokit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("parse all proto files", t, func() {
		var nsList []*Namespace
		nsList = append(nsList, NewNamespace(NamespaceGoogle, "../protokitgo/sdk/proto_google"))
		nsList = append(nsList, NewNamespace(NamespaceNetutils, "../protokitgo/sdk/proto_netutils_queue"))
		nsList = append(nsList, NewNamespace(NamespaceUser, "../protokitgo/example/protos"))
		m := NewParser(WithProtoFileAccessor(MustGetFileAccessorWithNamespace(nsList...)), WithGolangBasePackagePath("example/gen/golang"))
		m.Parse(nsList...)
	})
}
