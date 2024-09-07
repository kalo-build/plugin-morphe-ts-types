package tsdef

import "github.com/kaloseia/clone"

type Object struct {
	Name   string
	Fields []ObjectField
}

func (s Object) DeepClone() Object {
	return Object{
		Name:   s.Name,
		Fields: clone.DeepCloneSlice(s.Fields),
	}
}
