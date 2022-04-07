package protokit

const ValidateComment = "validate"

func (p *Parser) filterValidatorMessage() map[string]*ProtoMessage {
	msgs := make(map[string]*ProtoMessage)
	for _, protoFile := range p.protoFilePathToProtoFile {
		for _, msg := range protoFile.Messages {
			if msg.HasCommentField(ValidateComment) {
				msgs[msg.dotFullyQualifiedTypeName] = msg
			}
		}
	}
	return msgs
}

func (p *Parser) parseValidatorForMethod() {
	hms := p.filterValidatorMessage()
	for _, protoFile := range p.protoFilePathToProtoFile {
		for _, sg := range protoFile.ServiceGroups {
			for _, service := range sg.Services {
				for _, method := range service.Methods {
					if _, ok := hms[method.TypeInputDotFullQualifiedName]; ok {
						method.ValidatorInput = true
						service.HasValidator = true
					}
					if _, ok := hms[method.TypeOutputDotFullQualifiedName]; ok {
						method.ValidatorOutput = true
						service.HasValidator = true
					}
				}
			}
		}
	}
}
