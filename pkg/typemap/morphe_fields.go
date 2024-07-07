package typemap

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"

	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

var MorpheFieldToTsField = map[yaml.ModelFieldType]tsdef.TsType{
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
