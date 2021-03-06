package protokit

import (
	"github.com/jhump/protoreflect/desc"
)

func (p *Parser) getDottedPackage(fd *desc.FileDescriptor) string {
	dottedPkg := "." + fd.GetPackage()
	if dottedPkg != "." {
		dottedPkg += "."
	}
	return dottedPkg
}

func (p *Parser) parseMessages() {
	for _, pf := range p.protoFilePathToProtoFile {
		p.parseProtoFileMessages(pf)
	}
}

func getProtoForMessageDescriptor(md *desc.MessageDescriptor) *desc.FileDescriptor {
	return md.GetFile()
}

func (p *Parser) parsedMessageOrEnumGuard(d desc.Descriptor) bool {
	_, ok := p.tmpParsedMessageOrEnumMapping[d]
	if !ok {
		p.tmpParsedMessageOrEnumMapping[d] = struct{}{}
	}
	return ok
}

func (p *Parser) parseProtoFileMessage(pf *ProtoFile, md *desc.MessageDescriptor) {
	if p.parsedMessageOrEnumGuard(md) {
		return
	}
	pd := getProtoForMessageDescriptor(md)
	if pd != pf.fd {
		return
	}
	pm := p.BuildProtoMessage(pf, md)
	for _, f := range md.GetFields() {
		protoField := p.BuildProtoField(pf, f)
		pm.Fields = append(pm.Fields, protoField)
		// 解析import
		_, _ = p.addImportByDotFullyQualifiedTypeName(protoField.KeyTypeName, pm.ImportSet)
		_, _ = p.addImportByDotFullyQualifiedTypeName(protoField.ValueTypeName, pm.ImportSet)
	}
	pf.Messages = append(pf.Messages, pm)
	p.dotFullyQualifiedTypeNameToProtoMessage[pm.dotFullyQualifiedTypeName] = pm
	for _, mt := range md.GetNestedMessageTypes() {
		if mt.IsMapEntry() {
			continue
		}
		p.parseProtoFileMessage(pf, mt)
	}
	for _, et := range md.GetNestedEnumTypes() {
		p.parseProtoFileEnum(pf, et)
	}
}

func (p *Parser) parseProtoFileEnum(pf *ProtoFile, ed *desc.EnumDescriptor) {
	if p.parsedMessageOrEnumGuard(ed) {
		return
	}
	pf.Enums = append(pf.Enums, p.BuildProtoEnum(pf, ed))
}

func (p *Parser) parseProtoFileMessages(pf *ProtoFile) {
	for _, mt := range pf.fd.GetMessageTypes() {
		p.parseProtoFileMessage(pf, mt)
	}
	for _, et := range pf.fd.GetEnumTypes() {
		p.parseProtoFileEnum(pf, et)
	}
}
