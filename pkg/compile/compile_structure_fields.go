package compile

import (
	"github.com/kalo-build/go-util/core"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheStructure(r *registry.Registry, structureFields map[string]yaml.StructureField, fieldCasing cfg.Casing) ([]tsdef.ObjectField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(structureFields)
	for _, fieldName := range allFieldNames {
		field := structureFields[fieldName]
		fieldType, fieldTypeErr := getTsTypeForStructureField(r.GetAllEnums(), field, fieldCasing)
		if fieldTypeErr != nil {
			return nil, fieldTypeErr
		}

		// Structure fields are required by default; wrap in TsTypeOptional for "optional" attribute
		if hasAttribute(field.Attributes, "optional") {
			fieldType = tsdef.TsTypeOptional{ValueType: fieldType}
		}

		allFields = append(allFields, tsdef.ObjectField{
			Name: fieldCasing.Apply(fieldName),
			Type: fieldType,
		})
	}

	return allFields, nil
}

func hasAttribute(attributes []string, target string) bool {
	for _, attr := range attributes {
		if attr == target {
			return true
		}
	}
	return false
}

func getTsTypeForStructureField(allEnums map[string]yaml.Enum, field yaml.StructureField, fieldCasing cfg.Casing) (tsdef.TsType, error) {
	tsEnumType := getEnumFieldAsTsFieldType(allEnums, "", string(field.Type), fieldCasing)
	if tsEnumType.Type != nil {
		return tsEnumType.Type, nil
	}

	tsType, tsTypeExists := typemap.MorpheStructureFieldToTsField[field.Type]
	if !tsTypeExists {
		return nil, ErrUnsupportedMorpheFieldType(field.Type)
	}

	return tsType, nil
}
