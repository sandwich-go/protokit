package protokit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("parse all proto files", t, func() {
		var nsList []*Namespace
		nsList = append(nsList, NewNamespace(NamespaceProtokit, "../protokitgo/sdk/protokit"))
		nsList = append(nsList, NewNamespace(NamespaceValidate, "../protokitgo/sdk/proto_validate"))
		nsList = append(nsList, NewNamespace(NamespaceClaim, "../protokitgo/sdk/proto_claim"))
		nsList = append(nsList, NewNamespace(NamespaceGoogle, "../protokitgo/sdk/proto_google"))
		nsList = append(nsList, NewNamespace(NamespaceNetutils, "../protokitgo/sdk/proto_netutils_queue"))
		nsList = append(nsList, NewNamespace(NamespaceUser, "../protokitgo/example/protos"))
		m := NewParser(WithProtoFileAccessor(MustGetFileAccessorWithNamespace(nsList...)), WithGolangBasePackagePath("example/gen/golang"))
		m.Parse(nsList...)
	})
}

func TestParserParseEReturnsErrorForInvalidProto(t *testing.T) {
	dir := t.TempDir()
	protoPath := filepath.Join(dir, "bad.proto")
	content := `syntax = "proto3";

message Bad {
  string name = ;
}
`
	if err := os.WriteFile(protoPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write bad proto: %v", err)
	}

	ns := NewNamespace(NamespaceUser, dir)
	parser := NewParser(WithProtoFileAccessor(MustGetFileAccessorWithNamespace(ns)))
	err := parser.ParseE(ns)
	if err == nil {
		t.Fatalf("ParseE() err = nil, want parse error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "parse proto files under dir") {
		t.Fatalf("ParseE() err = %q, want proto parse context", msg)
	}
	if !strings.Contains(msg, "bad.proto") {
		t.Fatalf("ParseE() err = %q, want file name", msg)
	}
}
