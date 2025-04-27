package compile

import "github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"

// CompiledEnums maps Enum.Name -> CompiledEnum
type CompiledEnums map[string]CompiledEnum

func (enums CompiledEnums) AddCompiledEnum(enumDef *tsdef.Enum, enumContents []byte) {
	enums[enumDef.Name] = CompiledEnum{
		Enum:         enumDef,
		EnumContents: enumContents,
	}
}

func (enums CompiledEnums) GetCompiledEnum(enumName string) CompiledEnum {
	compiledEnum, compiledEnumExists := enums[enumName]
	if !compiledEnumExists {
		return CompiledEnum{}
	}
	return compiledEnum
}
