package tsdef

import "fmt"

type TsTypeArray struct {
	ValueType TsType
}

func (t TsTypeArray) IsPrimitive() bool {
	return false
}

func (t TsTypeArray) IsFunction() bool {
	return false
}

func (t TsTypeArray) IsArray() bool {
	return true
}

func (t TsTypeArray) IsObject() bool {
	return false
}

func (t TsTypeArray) IsInterface() bool {
	return false
}

func (t TsTypeArray) IsPromise() bool {
	return false
}

func (t TsTypeArray) IsOptional() bool {
	return false
}

func (t TsTypeArray) GetSyntax() string {
	return fmt.Sprintf("%s[]", t.ValueType.GetSyntax())
}

func (t TsTypeArray) DeepClone() TsTypeArray {
	return TsTypeArray{
		ValueType: DeepCloneTsType(t.ValueType),
	}
}
