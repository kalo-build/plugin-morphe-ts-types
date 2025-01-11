package compile

import (
	"strings"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheEntity(r *registry.Registry, entityFields map[string]yaml.EntityField, entityRelations map[string]yaml.EntityRelation) ([]tsdef.ObjectField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	allFields, fieldErr := getDirectTsFieldsForMorpheEntity(r, entityFields)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedTsFieldsForMorpheEntity(r, entityRelations)
	if relatedErr != nil {
		return nil, relatedErr
	}

	allFields = append(allFields, allRelatedFields...)
	return allFields, nil
}

func getDirectTsFieldsForMorpheEntity(r *registry.Registry, entityFields map[string]yaml.EntityField) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(entityFields)

	for _, fieldName := range allFieldNames {
		fieldDef := entityFields[fieldName]
		tsType, typeErr := getTsTypeForEntityField(r, fieldDef)
		if typeErr != nil {
			return nil, typeErr
		}

		typeField := tsdef.ObjectField{
			Name: fieldName,
			Type: tsType,
		}
		allFields = append(allFields, typeField)
	}

	return allFields, nil
}

func getTsTypeForEntityField(r *registry.Registry, field yaml.EntityField) (tsdef.TsType, error) {
	fieldPath := strings.Split(string(field.Type), ".")
	if len(fieldPath) < 2 {
		return nil, ErrInvalidEntityFieldPath(string(field.Type))
	}

	rootModelName := fieldPath[0]
	currentModel, modelErr := r.GetModel(rootModelName)
	if modelErr != nil {
		return nil, ErrRootModelNotFound(rootModelName)
	}

	for segmentIdx := 1; segmentIdx < len(fieldPath)-1; segmentIdx++ {
		relatedName := fieldPath[segmentIdx]
		_, exists := currentModel.Related[relatedName]
		if !exists {
			return nil, ErrRelatedModelNotFound(relatedName, string(field.Type))
		}

		relatedModel, relatedErr := r.GetModel(relatedName)
		if relatedErr != nil {
			return nil, ErrFailedToGetRelatedModel(relatedName, string(field.Type))
		}
		currentModel = relatedModel
	}

	terminalFieldName := fieldPath[len(fieldPath)-1]
	terminalField, exists := currentModel.Fields[terminalFieldName]
	if !exists {
		return nil, ErrTerminalFieldNotFound(terminalFieldName, string(field.Type))
	}

	tsEnumField := getEnumFieldAsTsFieldType(r.GetAllEnums(), terminalFieldName, string(terminalField.Type))
	if tsEnumField.Name != "" && tsEnumField.Type != nil {
		return tsEnumField.Type, nil
	}

	tsFieldType, typeSupported := typemap.MorpheModelFieldToTsField[terminalField.Type]
	if !typeSupported {
		return nil, ErrUnsupportedMorpheFieldType(terminalField.Type)
	}
	return tsFieldType, nil
}

func getRelatedTsFieldsForMorpheEntity(r *registry.Registry, entityRelations map[string]yaml.EntityRelation) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}

	allRelatedModelNames := core.MapKeysSorted(entityRelations)
	for _, relatedModelName := range allRelatedModelNames {
		entityRelation := entityRelations[relatedModelName]
		relatedModelDef, relatedModelDefErr := r.GetModel(relatedModelName)
		if relatedModelDefErr != nil {
			return nil, relatedModelDefErr
		}

		tsIDField, tsIDErr := getRelatedTsFieldForMorpheModelPrimaryID(entityRelation.Type, relatedModelName, relatedModelDef)
		if tsIDErr != nil {
			return nil, tsIDErr
		}
		allFields = append(allFields, tsIDField)

		tsRelatedField := getRelatedTsFieldForMorpheModelOptionalObject(entityRelation.Type, relatedModelName)
		allFields = append(allFields, tsRelatedField)
	}
	return allFields, nil
}
