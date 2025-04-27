package tsdef

import "github.com/kalo-build/clone"

type Object struct {
	Name    string
	Imports []ObjectImport
	Fields  []ObjectField
}

func (s Object) DeepClone() Object {
	return Object{
		Name:    s.Name,
		Imports: clone.DeepCloneSlice(s.Imports),
		Fields:  clone.DeepCloneSlice(s.Fields),
	}
}
