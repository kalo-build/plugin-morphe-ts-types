package compile

import (
	"fmt"
	"strings"

	"github.com/kalo/go-util/core"
	"github.com/kalo/go-util/strcase"
	"github.com/kalo/morphe-go/pkg/registry"
	"github.com/kalo/morphe-go/pkg/yaml"
	"github.com/kalo/morphe-go/pkg/yamlops"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo/plugin-morphe-ts-types/pkg/typemap"
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

	allRelatedEntityNames := core.MapKeysSorted(entityRelations)
	for _, relatedEntityName := range allRelatedEntityNames {
		entityRelation := entityRelations[relatedEntityName]
		relatedEntityDef, relatedEntityDefErr := r.GetEntity(relatedEntityName)
		if relatedEntityDefErr != nil {
			return nil, relatedEntityDefErr
		}

		tsIDField, tsIDErr := getRelatedTsFieldForMorpheEntityPrimaryID(r, entityRelation.Type, relatedEntityName, relatedEntityDef)

		if tsIDErr != nil {
			return nil, tsIDErr
		}
		allFields = append(allFields, tsIDField)

		tsRelatedField := getRelatedTsFieldForMorpheEntityOptionalObject(entityRelation.Type, relatedEntityName)
		allFields = append(allFields, tsRelatedField)

	}
	return allFields, nil
}

func getRelatedTsFieldForMorpheEntityPrimaryID(r *registry.Registry, relationType string, relatedEntityName string, relatedEntityDef yaml.Entity) (tsdef.ObjectField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetEntityPrimaryIdentifierFieldName(relatedEntityDef)
	if relatedIDFieldNameErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}
	idFieldName := fmt.Sprintf("%s%s", relatedEntityName, relatedPrimaryIDFieldName)

	relatedPrimaryIDFieldDef, relatedIDFieldDefErr := yamlops.GetEntityFieldDefinitionByName(relatedEntityDef, relatedPrimaryIDFieldName)
	if relatedIDFieldDefErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w (primary identifier)", relatedIDFieldDefErr)
	}
	idFieldType, typeErr := getTsTypeForEntityField(r, relatedPrimaryIDFieldDef)
	if typeErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w (primary identifier)", typeErr)
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

func getRelatedTsFieldForMorpheEntityOptionalObject(relationType string, relatedEntityName string) tsdef.ObjectField {
	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relatedEntityName + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeObject{
						ModulePath: "./" + strcase.ToKebabCaseLower(relatedEntityName),
						Name:       relatedEntityName,
					},
				},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relatedEntityName,
		Type: tsdef.TsTypeOptional{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./" + strcase.ToKebabCaseLower(relatedEntityName),
				Name:       relatedEntityName,
			},
		},
	}
	return tsRelatedField

}
