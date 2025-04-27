package tsdef

import "github.com/kalo-build/clone"

type ObjectImport struct {
	ModuleNames     []string
	ModulePath      string
	IsDefaultExport bool
}

func (i ObjectImport) DeepClone() ObjectImport {
	return ObjectImport{
		ModuleNames:     clone.Slice(i.ModuleNames),
		ModulePath:      i.ModulePath,
		IsDefaultExport: i.IsDefaultExport,
	}
}
