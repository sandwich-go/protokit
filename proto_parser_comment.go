package protokit

import "github.com/jhump/protoreflect/desc"

// parseComments解析fd的注释
func (p *Parser) parseComments(fd *desc.FileDescriptor) {
	for _, loc := range fd.AsFileDescriptorProto().GetSourceCodeInfo().GetLocation() {
		declaration := GetDefinitionAtPath(fd.AsFileDescriptorProto(), loc.Path)
		if c := NewComment(CommentLines(loc)); declaration != nil && c != nil {
			p.comments[declaration] = c
		}
	}
}
