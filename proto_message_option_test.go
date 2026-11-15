package protokit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse"
)

func TestCommonMessageOption(t *testing.T) {
	dir := t.TempDir()
	const source = `syntax = "proto3";
package example;
import "protokit/misc.proto";

message DefaultMessage {}
message ExplicitEnabled {
  option (protokit.message_opts).no_csharp = false;
}
message ServerOnly {
  option (protokit.message_opts).no_csharp = true;
}
message Outer {
  message Nested {
    option (protokit.message_opts).no_csharp = true;
  }
}
`
	if err := os.WriteFile(filepath.Join(dir, "message_option.proto"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	parser := protoparse.Parser{ImportPaths: []string{dir, "option/protos"}}
	files, err := parser.ParseFiles("message_option.proto")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name       string
		configured bool
		noCsharp   bool
	}{
		{"example.DefaultMessage", false, false},
		{"example.ExplicitEnabled", true, false},
		{"example.ServerOnly", true, true},
		{"example.Outer", false, false},
		{"example.Outer.Nested", true, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			md := files[0].FindMessage(tt.name)
			if md == nil {
				t.Fatal("message descriptor not found")
			}
			message := NewProtoMessage(&ProtoFile{}, md)
			opts := message.GetCommonMessageOption()
			if (opts != nil) != tt.configured {
				t.Fatalf("option configured = %v, want %v", opts != nil, tt.configured)
			}
			if got := opts.GetNoCsharp(); got != tt.noCsharp {
				t.Fatalf("no_csharp = %v, want %v", got, tt.noCsharp)
			}
		})
	}
}

func TestPhantomDomainOption(t *testing.T) {
	dir := t.TempDir()
	// 经过实际 proto 解析和扩展解码，确保 option 描述符与生成的 Go 类型一致。
	const source = `syntax = "proto3";
package example;
import "protokit/orm.proto";
message User {
  option (protokit.orm_message).bean = {};
  option (protokit.orm_message).phantom_domain = "user";
}
message TruckBrief {
  option (protokit.orm_message).bean = {phantom: {depend_bean: "truck"}};
  option (protokit.orm_message).phantom_domain = "user";
}
message Unconfigured {
  option (protokit.orm_message).bean = {};
}
`
	if err := os.WriteFile(filepath.Join(dir, "phantom_domain.proto"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	parser := protoparse.Parser{ImportPaths: []string{dir, "option/protos"}}
	files, err := parser.ParseFiles("phantom_domain.proto")
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{"User": "user", "TruckBrief": "user", "Unconfigured": ""} {
		message := NewProtoMessage(&ProtoFile{}, files[0].FindMessage("example."+name))
		if got := message.GetOrmOption().GetPhantomDomain(); got != want {
			t.Errorf("%s phantom_domain = %q, want %q", name, got, want)
		}
	}
}
