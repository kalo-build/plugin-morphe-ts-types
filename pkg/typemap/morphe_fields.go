package typemap

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"

	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

var MorpheModelFieldToTsField = map[yaml.ModelFieldType]tsdef.TsType{
	yaml.ModelFieldTypeUUID:          tsdef.TsTypeString,
	yaml.ModelFieldTypeAutoIncrement: tsdef.TsTypeNumber,
	yaml.ModelFieldTypeString:        tsdef.TsTypeString,
	yaml.ModelFieldTypeInteger:       tsdef.TsTypeNumber,
	yaml.ModelFieldTypeFloat:         tsdef.TsTypeNumber,
	yaml.ModelFieldTypeBoolean:       tsdef.TsTypeBoolean,
	yaml.ModelFieldTypeTime:          tsdef.TsTypeDate,
	yaml.ModelFieldTypeDate:          tsdef.TsTypeDate,
	yaml.ModelFieldTypeProtected:     tsdef.TsTypeString,
	yaml.ModelFieldTypeSealed:        tsdef.TsTypeString,
}

var MorpheStructureFieldToTsField = map[yaml.StructureFieldType]tsdef.TsType{
	yaml.StructureFieldTypeUUID:          tsdef.TsTypeString,
	yaml.StructureFieldTypeAutoIncrement: tsdef.TsTypeNumber,
	yaml.StructureFieldTypeString:        tsdef.TsTypeString,
	yaml.StructureFieldTypeInteger:       tsdef.TsTypeNumber,
	yaml.StructureFieldTypeFloat:         tsdef.TsTypeNumber,
	yaml.StructureFieldTypeBoolean:       tsdef.TsTypeBoolean,
	yaml.StructureFieldTypeTime:          tsdef.TsTypeDate,
	yaml.StructureFieldTypeDate:          tsdef.TsTypeDate,
	yaml.StructureFieldTypeProtected:     tsdef.TsTypeString,
	yaml.StructureFieldTypeSealed:        tsdef.TsTypeString,
}
