package protokit

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProto(t *testing.T) {
	Convey("proto field name converter", t, func() {
		for _, test := range []struct {
			rawName    string
			golangName string
		}{
			{rawName: "_tags_", golangName: "XTags_"},
			{rawName: "_tags", golangName: "XTags"},
			{rawName: "param1_args1", golangName: "Param1Args1"},
			{rawName: "param1__args1", golangName: "Param1_Args1"},
		} {
			So(GoFieldName(test.rawName), ShouldEqual, test.golangName)
		}
	})
	Convey("proto struct name converter", t, func() {

		//(*OuterTestT1)(nil),             // 2: msg.Outer.test_t1
		//(*OuterTest_T2)(nil),            // 3: msg.Outer.test_T2
		//(*Outer_TestT3)(nil),            // 4: msg.Outer.Test_t3
		//(*Outer_Test_T4)(nil),           // 5: msg.Outer.Test_T4
		//(*Outer_Test_T5)(nil),           // 6: msg.Outer.Test__t5
		//(*Outer_Test__T6)(nil),          // 7: msg.Outer.Test__T6
		//(*Outer_TestT7__)(nil),          // 8: msg.Outer.Test_t7__
		//(*Outer_TestT8__)(nil),          // 9: msg.Outer.Test_t8__
		//(*Outer_Test____T9__)(nil),      // 10: msg.Outer.Test_____t9__
		//(*Outer_Test____T9__)(nil),      // 10: msg.Outer.Test_____t9__
		// ActivityData  .FunPlus.ServerCommon.Config.ActivityData

		// _sym_db.RegisterMessage(Outer)
		// _sym_db.RegisterMessage(Outer.test_t1)
		// _sym_db.RegisterMessage(Outer.test_T2)
		// _sym_db.RegisterMessage(Outer.Test_t3)
		// _sym_db.RegisterMessage(Outer.Test_T4)
		// _sym_db.RegisterMessage(Outer.Test__t5)
		// _sym_db.RegisterMessage(Outer.Test__T6)
		// _sym_db.RegisterMessage(Outer.Test_t7__)
		// _sym_db.RegisterMessage(Outer.Test_t8__)
		// _sym_db.RegisterMessage(Outer.Test_____t9__)
		for _, test := range []struct {
			// input
			fullyQualifiedName string
			// output
			goStructNameWithPackage     string
			pythonStructNameWithPackage string
			protoPackage                string
			golangPackage               string
		}{
			{fullyQualifiedName: ".msg.Outer.test_t1", goStructNameWithPackage: "msg.OuterTestT1", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: ".msg.Outer.test_T2", goStructNameWithPackage: "msg.OuterTest_T2", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: ".msg.Outer.Test_t3", goStructNameWithPackage: "msg.Outer_TestT3", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: ".msg.Outer.Test_T4", goStructNameWithPackage: "msg.Outer_Test_T4", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: ".msg.Outer.Test__t5", goStructNameWithPackage: "msg.Outer_Test_T5", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: "msg.Outer.Test__T6", goStructNameWithPackage: "msg.Outer_Test__T6", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: "msg.Outer.Test_t7__", goStructNameWithPackage: "msg.Outer_TestT7__", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: "msg.Outer.Test_t8__", goStructNameWithPackage: "msg.Outer_TestT8__", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: "msg.Outer.Test_____t9__", goStructNameWithPackage: "msg.Outer_Test____T9__", protoPackage: "msg", golangPackage: "msg"},
			{fullyQualifiedName: ".FunPlus.ServerCommon.Config.ActivityData", goStructNameWithPackage: "bus.ActivityData", protoPackage: "FunPlus.ServerCommon.Config", golangPackage: "bus"},
			{fullyQualifiedName: ".test.PBBIGve2ndBattle", goStructNameWithPackage: "test.PBBIGve2NdBattle", protoPackage: "test", golangPackage: "test"},
		} {
			test.pythonStructNameWithPackage = strings.TrimPrefix(test.fullyQualifiedName, ".")
			So(GoStructNameWithGolangPackage(test.fullyQualifiedName, test.protoPackage, test.golangPackage), ShouldEqual, test.goStructNameWithPackage)
		}
	})
}
