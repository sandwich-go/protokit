package protokit

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
)

const (
	tagFileDescriptorMessageType int32 = 4
	tagFileDescriptorEnumType    int32 = 5
	tagFileDescriptorServiceType int32 = 6
	tagDescriptorField           int32 = 2
	tagDescriptorNestedType      int32 = 3
	tagDescriptorEnumType        int32 = 4
	tagDescriptorOneOfDecl       int32 = 8
	tagEnumDescriptorValue       int32 = 2
	tagMethodDescriptorValue     int32 = 2
)

// GetDefinitionAtPath 解析comment使用
func GetDefinitionAtPath(file *descriptor.FileDescriptorProto, path []int32) proto.Message {
	var pos proto.Message = file
	for step := 0; step < len(path); step++ {
		switch p := pos.(type) {
		case *descriptor.FileDescriptorProto:
			switch path[step] {
			case tagFileDescriptorMessageType:
				step++
				pos = p.MessageType[path[step]]
			case tagFileDescriptorEnumType:
				step++
				pos = p.EnumType[path[step]]
			case tagFileDescriptorServiceType:
				step++
				pos = p.Service[path[step]]
			default:
				return nil
			}
		case *descriptor.DescriptorProto:
			switch path[step] {
			case tagDescriptorField:
				step++
				pos = p.Field[path[step]]
			case tagDescriptorNestedType:
				step++
				pos = p.NestedType[path[step]]
			case tagDescriptorEnumType:
				step++
				pos = p.EnumType[path[step]]
			case tagDescriptorOneOfDecl:
				step++
				pos = p.OneofDecl[path[step]]
			default:
				return nil
			}

		case *descriptor.EnumDescriptorProto:
			switch path[step] {
			case tagEnumDescriptorValue:
				step++
				pos = p.Value[path[step]]
			default:
				return nil
			}
		case *descriptor.ServiceDescriptorProto:
			switch path[step] {
			case tagMethodDescriptorValue:
				step++
				pos = p.Method[path[step]]
			default:
				return nil
			}
		default:
			return nil
		}
	}
	return pos
}
