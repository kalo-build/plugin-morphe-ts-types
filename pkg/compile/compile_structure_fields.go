package compile

import (
	"github.com/kaloseia/go-util/core"

	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheStructure(r *registry.Registry, structureFields map[string]yaml.StructureField) ([]tsdef.ObjectField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(structureFields)
	for _, fieldName := range allFieldNames {
		field := structureFields[fieldName]
		fieldType, fieldTypeErr := getTsTypeForStructureField(r.GetAllEnums(), field)
		if fieldTypeErr != nil {
			return nil, fieldTypeErr
		}

		allFields = append(allFields, tsdef.ObjectField{
			Name: fieldName,
			Type: fieldType,
		})
	}

	return allFields, nil
}

func getTsTypeForStructureField(allEnums map[string]yaml.Enum, field yaml.StructureField) (tsdef.TsType, error) {
	tsEnumType := getEnumFieldAsTsFieldType(allEnums, "", string(field.Type))
	if tsEnumType.Type != nil {
		return tsEnumType.Type, nil
	}

	tsType, tsTypeExists := typemap.MorpheStructureFieldToTsField[field.Type]
	if !tsTypeExists {
		return nil, ErrUnsupportedMorpheFieldType(field.Type)
	}

	return tsType, nil
}
