package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/morphe-go/pkg/yamlops"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheModel(r *registry.Registry, modelFields map[string]yaml.ModelField, modelRelations map[string]yaml.ModelRelation) ([]tsdef.ObjectField, error) {
	allFields, fieldErr := getDirectTsFieldsForMorpheModel(modelFields)
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

func getDirectTsFieldsForMorpheModel(modelFields map[string]yaml.ModelField) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]
		tsFieldType, typeSupported := typemap.MorpheFieldToTsField[fieldDef.Type]
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

		if yamlops.IsRelationFor(modelRelation.Type) {
			tsIDField, tsIDErr := getRelatedTsFieldForMorpheModelPrimaryID(modelRelation.Type, relatedModelName, relatedModelDef)
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)
		}

		tsRelatedField := getRelatedTsFieldForMorpheModelOptionalObject(modelRelation.Type, relatedModelName)
		allFields = append(allFields, tsRelatedField)
	}
	return allFields, nil
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
	idFieldType, typeSupported := typemap.MorpheFieldToTsField[relatedPrimaryIDFieldDef.Type]
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
					ValueType: tsdef.TsTypeObject{Name: relatedModelName},
				},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relatedModelName,
		Type: tsdef.TsTypeOptional{
			ValueType: tsdef.TsTypeObject{Name: relatedModelName},
		},
	}
	return tsRelatedField
}
