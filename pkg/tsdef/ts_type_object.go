package tsdef

type TsTypeObject struct {
	Name string
}

func (t TsTypeObject) IsPrimitive() bool {
	return false
}

func (t TsTypeObject) IsFunction() bool {
	return false
}

func (t TsTypeObject) IsArray() bool {
	return false
}

func (t TsTypeObject) IsObject() bool {
	return true
}

func (t TsTypeObject) IsInterface() bool {
	return false
}

func (t TsTypeObject) IsPromise() bool {
	return false
}

func (t TsTypeObject) IsOptional() bool {
	return false
}

func (t TsTypeObject) GetSyntax() string {
	return t.Name
}

func (t TsTypeObject) DeepClone() TsTypeObject {
	return TsTypeObject{
		Name: t.Name,
	}
}
