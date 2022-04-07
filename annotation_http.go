package protokit

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// HTTPPath 获取指定的method描述中annotation中的的http访问路径
func HTTPPath(m *descriptor.MethodDescriptorProto) (string, error) {
	ext, err := proto.GetExtension(m.Options, annotations.E_Http)
	if err != nil {
		return "", err
	}
	opts, ok := ext.(*annotations.HttpRule)
	if !ok {
		return "", fmt.Errorf("extension is %T; want an HttpRule", ext)
	}

	switch t := opts.Pattern.(type) {
	default:
		return "", nil
	case *annotations.HttpRule_Get:
		return t.Get, nil
	case *annotations.HttpRule_Post:
		return t.Post, nil
	case *annotations.HttpRule_Put:
		return t.Put, nil
	case *annotations.HttpRule_Delete:
		return t.Delete, nil
	case *annotations.HttpRule_Patch:
		return t.Patch, nil
	case *annotations.HttpRule_Custom:
		return t.Custom.Path, nil
	}
}
