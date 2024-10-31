package tsdef

type TsTypeOptional struct {
	ValueType TsType
}

func (t TsTypeOptional) IsPrimitive() bool {
	return false
}

func (t TsTypeOptional) IsFunction() bool {
	return false
}

func (t TsTypeOptional) IsArray() bool {
	return false
}

func (t TsTypeOptional) IsObject() bool {
	return false
}

func (t TsTypeOptional) IsInterface() bool {
	return false
}

func (t TsTypeOptional) IsPromise() bool {
	return false
}

func (t TsTypeOptional) IsOptional() bool {
	return true
}

func (t TsTypeOptional) GetSyntax() string {
	return t.ValueType.GetSyntax()
}

func (t TsTypeOptional) DeepClone() TsTypeOptional {
	return TsTypeOptional{
		ValueType: DeepCloneTsType(t.ValueType),
	}
}
