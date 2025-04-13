package compile

import (
	"fmt"

	"github.com/kalo/go-util/core"
	"github.com/kalo/go-util/strcase"
	"github.com/kalo/morphe-go/pkg/registry"
	"github.com/kalo/morphe-go/pkg/yaml"
	"github.com/kalo/morphe-go/pkg/yamlops"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheModel(r *registry.Registry, modelFields map[string]yaml.ModelField, modelRelations map[string]yaml.ModelRelation) ([]tsdef.ObjectField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}
	allFields, fieldErr := getDirectTsFieldsForMorpheModel(r.GetAllEnums(), modelFields)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedTsFieldsForMorpheModel(r, modelRelations)
	if relatedErr != nil {
		return nil, relatedErr
	}

	allFields = append(allFields, allRelatedFields...)
	return allFields, nil
}

func getDirectTsFieldsForMorpheModel(allEnums map[string]yaml.Enum, modelFields map[string]yaml.ModelField) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]

		tsEnumField := getEnumFieldAsTsFieldType(allEnums, fieldName, string(fieldDef.Type))
		if tsEnumField.Name != "" && tsEnumField.Type != nil {
			allFields = append(allFields, tsEnumField)
			continue
		}

		tsFieldType, typeSupported := typemap.MorpheModelFieldToTsField[fieldDef.Type]
		if !typeSupported {
			return nil, ErrUnsupportedMorpheFieldType(fieldDef.Type)
		}
		tsField := tsdef.ObjectField{
			Name: fieldName,
			Type: tsFieldType,
		}
		allFields = append(allFields, tsField)
	}

	return allFields, nil
}

func getRelatedTsFieldsForMorpheModel(r *registry.Registry, modelRelations map[string]yaml.ModelRelation) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}

	allRelatedModelNames := core.MapKeysSorted(modelRelations)
	for _, relatedModelName := range allRelatedModelNames {
		modelRelation := modelRelations[relatedModelName]
		relatedModelDef, relatedModelDefErr := r.GetModel(relatedModelName)
		if relatedModelDefErr != nil {
			return nil, relatedModelDefErr
		}

		tsIDField, tsIDErr := getRelatedTsFieldForMorpheModelPrimaryID(modelRelation.Type, relatedModelName, relatedModelDef)
		if tsIDErr != nil {
			return nil, tsIDErr
		}
		allFields = append(allFields, tsIDField)

		tsRelatedField := getRelatedTsFieldForMorpheModelOptionalObject(modelRelation.Type, relatedModelName)
		allFields = append(allFields, tsRelatedField)
	}
	return allFields, nil
}

func getEnumFieldAsTsFieldType(allEnums map[string]yaml.Enum, fieldName string, enumName string) tsdef.ObjectField {
	if len(allEnums) == 0 {
		return tsdef.ObjectField{}
	}

	_, enumTypeExists := allEnums[enumName]
	if !enumTypeExists {
		return tsdef.ObjectField{}
	}

	tsFieldType := tsdef.TsTypeObject{
		ModulePath: "../enums/" + strcase.ToKebabCaseLower(enumName),
		Name:       enumName,
	}
	tsField := tsdef.ObjectField{
		Name: fieldName,
		Type: tsFieldType,
	}
	return tsField
}

func getRelatedTsFieldForMorpheModelPrimaryID(relationType string, relatedModelName string, relatedModelDef yaml.Model) (tsdef.ObjectField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetModelPrimaryIdentifierFieldName(relatedModelDef)
	if relatedIDFieldNameErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}
	idFieldName := fmt.Sprintf("%s%s", relatedModelName, relatedPrimaryIDFieldName)

	relatedPrimaryIDFieldDef, relatedIDFieldDefErr := yamlops.GetModelFieldDefinitionByName(relatedModelDef, relatedPrimaryIDFieldName)
	if relatedIDFieldDefErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w (primary identifier)", relatedIDFieldDefErr)
	}
	idFieldType, typeSupported := typemap.MorpheModelFieldToTsField[relatedPrimaryIDFieldDef.Type]
	if !typeSupported {
		return tsdef.ObjectField{}, ErrUnsupportedMorpheFieldType(relatedPrimaryIDFieldDef.Type)
	}

	if yamlops.IsRelationMany(relationType) {
		tsIDField := tsdef.ObjectField{
			Name: idFieldName + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: idFieldType,
				},
			},
		}
		return tsIDField, nil
	}

	tsIDField := tsdef.ObjectField{
		Name: idFieldName,
		Type: tsdef.TsTypeOptional{
			ValueType: idFieldType,
		},
	}
	return tsIDField, nil
}

func getRelatedTsFieldForMorpheModelOptionalObject(relationType string, relatedModelName string) tsdef.ObjectField {
	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relatedModelName + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeObject{
						ModulePath: "./" + strcase.ToKebabCaseLower(relatedModelName),
						Name:       relatedModelName,
					},
				},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relatedModelName,
		Type: tsdef.TsTypeOptional{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./" + strcase.ToKebabCaseLower(relatedModelName),
				Name:       relatedModelName,
			},
		},
	}
	return tsRelatedField
}
