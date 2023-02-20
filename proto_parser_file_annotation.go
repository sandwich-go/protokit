package protokit

import (
	"strings"

	"github.com/sandwich-go/boost/misc/annotation"
	"github.com/sandwich-go/boost/xpanic"
)

const file_annotation_magic = AnnotationPrefix + AnnotationSpaceAnnotationType

func (p *Parser) parseAnnotation() {
	for _, protoFile := range p.protoFilePathToProtoFile {
		if !strings.Contains(protoFile.Content, file_annotation_magic) {
			continue
		}
		lines := strings.Split(protoFile.Content, "\n")
		var annotationLines []string
		for _, line := range lines {
			if strings.Contains(line, AnnotationPrefix) {
				commentIndex := strings.Index(line, "//")
				if commentIndex != -1 {
					comment := line[commentIndex:]
					annotationLines = append(annotationLines, comment)
				}
			}
		}
		if len(annotationLines) > 0 {
			var err error
			protoFile.Annotations, err = annotation.New().ResolveNoDuplicate(annotationLines...)
			xpanic.WhenErrorAsFmtFirst(err, "got error:%w while parse file annotation for file:%s", protoFile.FilePath)
		}
	}
}
