package protokit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("parse all proto files", t, func() {
		cc := NewOptions()
		var nsList []*Namespace
		nsList = append(nsList, NewNamespace(NamespaceGoogle, "../sdk/proto_google", cc))
		nsList = append(nsList, NewNamespace(NamespaceNetutils, "../sdk/proto_netutils_queue", cc))
		nsList = append(nsList, NewNamespace(NamespaceUser, "../example/protos", cc))
		m := NewParser(WithProtoFileAccessor(MustGetFileAccessorWithNamespace(nsList...)), WithGolangBasePackagePath("example/gen/golang"))
		m.Parse(nsList...)
	})
}
