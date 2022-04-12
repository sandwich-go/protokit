package protokit

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

func IsRepeatedField(f *descriptorpb.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Type != nil && f.Label != nil && *f.Label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return true
	}

	return false
}

func IsEnumField(f *descriptorpb.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Type != nil && f.Label != nil && *f.Type == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
		return true
	}

	return false
}

func IsMessageField(f *descriptorpb.FieldDescriptorProto) bool {
	if f == nil {
		return false
	}
	if f.Type != nil && f.Label != nil && *f.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return true
	}

	return false
}
func shortType(s string) string {
	t := strings.Split(s, ".")
	return t[len(t)-1]
}
func IsMapField(f *descriptorpb.FieldDescriptorProto, m *descriptorpb.DescriptorProto) bool {
	if f.TypeName == nil {
		return false
	}

	shortName := shortType(*f.TypeName)
	var nt *descriptorpb.DescriptorProto
	for _, t := range m.NestedType {
		if *t.Name == shortName {
			nt = t
			break
		}
	}

	if nt == nil {
		return false
	}

	for _, f := range nt.Field {
		switch *f.Name {
		case "key":
			if *f.Number != 1 {
				return false
			}
		case "value":
			if *f.Number != 2 {
				return false
			}
		default:
			return false
		}
	}

	return true
}
