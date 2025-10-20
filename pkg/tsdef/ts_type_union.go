package tsdef

import "strings"

type TsTypeUnion struct {
	Types []TsType
}

func (t TsTypeUnion) IsPrimitive() bool {
	return false
}

func (t TsTypeUnion) IsFunction() bool {
	return false
}

func (t TsTypeUnion) IsArray() bool {
	return false
}

func (t TsTypeUnion) IsObject() bool {
	return false
}

func (t TsTypeUnion) IsInterface() bool {
	return false
}

func (t TsTypeUnion) IsPromise() bool {
	return false
}

func (t TsTypeUnion) IsOptional() bool {
	return false
}

func (t TsTypeUnion) GetImports() []ObjectImport {
	importMap := map[string]ObjectImport{}
	for _, unionType := range t.Types {
		for _, imp := range unionType.GetImports() {
			importMap[imp.ModulePath] = imp
		}
	}

	imports := []ObjectImport{}
	for _, imp := range importMap {
		imports = append(imports, imp)
	}
	return imports
}

func (t TsTypeUnion) GetSyntax() string {
	syntaxes := []string{}
	for _, unionType := range t.Types {
		syntaxes = append(syntaxes, unionType.GetSyntax())
	}
	return strings.Join(syntaxes, " | ")
}

func (t TsTypeUnion) DeepClone() TsTypeUnion {
	return TsTypeUnion{
		Types: DeepCloneTsTypeSlice(t.Types),
	}
}
