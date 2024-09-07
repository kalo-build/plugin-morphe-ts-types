package tsdef

type ObjectField struct {
	Name string
	Type TsType
}

func (f ObjectField) DeepClone() ObjectField {
	return ObjectField{
		Name: f.Name,
		Type: DeepCloneTsType(f.Type),
	}
}
