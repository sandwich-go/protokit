package protokit

type EnumOrMessage interface {
	DotFullyQualifiedTypeName() string
	GoNameWithoutGolangPackage() string
	GetProtoFile() *ProtoFile
}
