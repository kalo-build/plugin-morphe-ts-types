package tsdef

type EnumEntry struct {
	Name  string
	Value any
}

func (f EnumEntry) DeepClone() EnumEntry {
	return EnumEntry{
		Name:  f.Name,
		Value: f.Value,
	}
}
