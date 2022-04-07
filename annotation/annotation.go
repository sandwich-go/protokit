package annotation

import (
	"strings"

	"github.com/sandwich-go/protokit/util"
)

type Register interface {
	ResolveAnnotations(annotationLines []string) []Annotation
	ResolveAnnotationByName(annotationLines []string, name string) Annotation
	ResolveAnnotation(annotationLines string) (Annotation, bool)
}

func ResolveAnnotationByName(methodCommentLines []string, key string) Annotation {
	registry := NewRegistry()
	return registry.ResolveAnnotationByName(methodCommentLines, key)
}

type annotationRegistry struct {
	descriptors []*Descriptor
}

const all = "*"

func NewRegistry(descriptors ...*Descriptor) Register {
	v := &annotationRegistry{
		descriptors: descriptors,
	}
	if len(v.descriptors) == 0 {
		v.descriptors = append(v.descriptors, &Descriptor{Name: all})
	}
	return v
}

type Annotation struct {
	Name       string
	Attributes map[string]string
}

func (a Annotation) Has(key string) bool {
	_, ok := a.Attributes[key]
	return ok
}

func (a Annotation) GetBool(key string, defaultVal ...bool) bool {
	v, ok := a.Attributes[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return false
	}
	lowerV := util.Trim(strings.ToLower(v))
	if lowerV == "1" || lowerV == "true" || lowerV == "y" || lowerV == "yes" {
		return true
	}
	return false
}

func (a Annotation) GetString(key string, defaultVal ...string) string {
	v, ok := a.Attributes[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return util.Trim(v)
}

type validationFunc func(a Annotation) bool

type Descriptor struct {
	Name      string
	Validator validationFunc
}

func (ar *annotationRegistry) ResolveAnnotations(annotationLines []string) []Annotation {
	annotations := make([]Annotation, 0)
	for _, line := range annotationLines {
		if ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line)); ok {
			annotations = append(annotations, ann)
		}
	}
	return annotations
}

func (ar *annotationRegistry) ResolveAnnotationByName(annotationLines []string, name string) Annotation {
	for _, line := range annotationLines {
		ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line))
		if ok && ann.Name == name {
			return ann
		}
	}
	return Annotation{}
}

func (ar *annotationRegistry) ResolveAnnotation(annotationLines string) (Annotation, bool) {
	for _, descriptor := range ar.descriptors {
		if !strings.Contains(annotationLines, magicPrefix) {
			continue
		}

		ann, err := parseAnnotation(annotationLines)
		if err != nil {
			panic(err)
		}

		if descriptor.Name != all && ann.Name != descriptor.Name {
			continue
		}

		if descriptor.Validator != nil {
			ok := descriptor.Validator(ann)
			if !ok {
				continue
			}

		}
		return ann, true
	}
	return Annotation{}, false
}
