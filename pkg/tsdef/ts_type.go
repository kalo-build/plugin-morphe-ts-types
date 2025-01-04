package tsdef

type TsType interface {
	IsPrimitive() bool
	IsFunction() bool
	IsArray() bool
	IsObject() bool
	IsInterface() bool
	IsPromise() bool
	IsOptional() bool

	GetImports() []ObjectImport
	GetSyntax() string
}
