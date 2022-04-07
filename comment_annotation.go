package protokit

import (
	"strings"

	"github.com/sandwich-go/protokit/annotation"
	"github.com/sandwich-go/protokit/util"
	"google.golang.org/protobuf/types/descriptorpb"
)

const AnnotationService = "service"
const AnnotationGlobal = "global"

type Comment struct {
	Content     string
	Tags        map[string]string
	Annotations []annotation.Annotation
}

func CommentLines(loc *descriptorpb.SourceCodeInfo_Location) []string {
	var lines []string
	for _, str := range loc.GetLeadingDetachedComments() {
		if s := strings.TrimSpace(str); s != "" {
			lines = append(lines, s)
		}
	}
	arr := strings.Split(loc.GetLeadingComments(), "\n")
	for _, v := range arr {
		if s := strings.TrimSpace(v); s != "" {
			lines = append(lines, s)
		}
	}
	if s := strings.TrimSpace(loc.GetTrailingComments()); s != "" {
		lines = append(lines, s)
	}
	return lines
}

func GetAnnotation(c *Comment, name string) annotation.Annotation {
	if c == nil {
		return annotation.Annotation{}
	}
	return c.Annotation(name)
}

func (c *Comment) Annotation(name ...string) annotation.Annotation {
	annotationName := AnnotationGlobal
	if len(name) != 0 {
		annotationName = name[0]
	}
	for _, v := range c.Annotations {
		if v.Name == annotationName {
			return v
		}
	}
	return annotation.Annotation{}
}

func NewComment(lines []string) *Comment {
	comment := &Comment{
		Tags: make(map[string]string),
	}
	if len(lines) == 0 {
		return comment
	}
	comment.Content = strings.Join(lines, "\n")
	comment.Annotations = annotation.NewRegistry().ResolveAnnotations(lines)
	for _, l := range lines {
		l = util.Trim(l)
		arr := strings.Split(l, "@")
		if len(arr) <= 1 {
			continue
		}
		for _, v := range arr {
			if len(v) == 0 {
				continue
			}
			// todo，目前的stamp只支持最多一个参数
			kv := strings.Split(v, "=")
			key := strings.Trim(strings.ToLower(kv[0]), " ")
			if len(kv) > 1 {
				comment.Tags[key] = kv[1]
			} else {
				comment.Tags[key] = "true"
			}
		}
	}

	global := annotation.Annotation{
		Name:       AnnotationGlobal,
		Attributes: make(map[string]string),
	}

	for k, v := range comment.Tags {
		global.Attributes[k] = v
	}
	comment.Annotations = append(comment.Annotations, global)

	return comment
}
