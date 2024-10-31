package tsdef

type TsTypePrimitive struct {
	Syntax string
}

func (t TsTypePrimitive) IsPrimitive() bool {
	return true
}

func (t TsTypePrimitive) IsFunction() bool {
	return false
}

func (t TsTypePrimitive) IsArray() bool {
	return false
}

func (t TsTypePrimitive) IsObject() bool {
	return false
}

func (t TsTypePrimitive) IsInterface() bool {
	return false
}

func (t TsTypePrimitive) IsPromise() bool {
	return false
}

func (t TsTypePrimitive) IsOptional() bool {
	return false
}

func (t TsTypePrimitive) GetSyntax() string {
	return t.Syntax
}

func (t TsTypePrimitive) DeepClone() TsTypePrimitive {
	return TsTypePrimitive{
		Syntax: t.Syntax,
	}
}
