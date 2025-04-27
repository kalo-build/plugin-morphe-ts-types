package tsdef

import "github.com/kalo-build/clone"

type Enum struct {
	Name    string
	Type    TsType
	Entries []EnumEntry
}

func (s Enum) DeepClone() Enum {
	return Enum{
		Name:    s.Name,
		Type:    DeepCloneTsType(s.Type),
		Entries: clone.DeepCloneSlice(s.Entries),
	}
}
