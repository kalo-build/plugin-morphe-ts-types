package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/typemap"
)

func getTsFieldsForMorpheEntity(r *registry.Registry, entityFields map[string]yaml.EntityField, entityRelations map[string]yaml.EntityRelation, fieldCasing cfg.Casing) ([]tsdef.ObjectField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	allFields, fieldErr := getDirectTsFieldsForMorpheEntity(r, entityFields, fieldCasing)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedTsFieldsForMorpheEntity(r, entityRelations, fieldCasing)
	if relatedErr != nil {
		return nil, relatedErr
	}

	allFields = append(allFields, allRelatedFields...)
	return allFields, nil
}

func getDirectTsFieldsForMorpheEntity(r *registry.Registry, entityFields map[string]yaml.EntityField, fieldCasing cfg.Casing) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}
	allFieldNames := core.MapKeysSorted(entityFields)

	for _, fieldName := range allFieldNames {
		fieldDef := entityFields[fieldName]
		tsType, typeErr := getTsTypeForEntityField(r, fieldDef)
		if typeErr != nil {
			return nil, typeErr
		}

		// Entity fields are required by default; wrap in TsTypeOptional for "optional" attribute
		var finalType tsdef.TsType = tsType
		if hasAttribute(fieldDef.Attributes, "optional") {
			finalType = tsdef.TsTypeOptional{ValueType: tsType}
		}

		typeField := tsdef.ObjectField{
			Name: fieldCasing.Apply(fieldName),
			Type: finalType,
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
		relation, exists := currentModel.Related[relatedName]
		if !exists {
			return nil, ErrRelatedModelNotFound(relatedName, string(field.Type))
		}

		targetModelName := relatedName
		if strings.TrimSpace(relation.Aliased) != "" {
			targetModelName = strings.TrimSpace(relation.Aliased)
		}

		relatedModel, relatedErr := r.GetModel(targetModelName)
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

	tsEnumField := getEnumFieldAsTsFieldType(r.GetAllEnums(), terminalFieldName, string(terminalField.Type), cfg.Casing(""))
	if tsEnumField.Name != "" && tsEnumField.Type != nil {
		return tsEnumField.Type, nil
	}

	tsFieldType, typeSupported := typemap.MorpheModelFieldToTsField[terminalField.Type]
	if !typeSupported {
		return nil, ErrUnsupportedMorpheFieldType(terminalField.Type)
	}
	return tsFieldType, nil
}

func getRelatedTsFieldsForMorpheEntity(r *registry.Registry, entityRelations map[string]yaml.EntityRelation, fieldCasing cfg.Casing) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}

	allRelatedEntityNames := core.MapKeysSorted(entityRelations)
	for _, relationshipName := range allRelatedEntityNames {
		entityRelation := entityRelations[relationshipName]

		// Handle different relationship types
		switch entityRelation.Type {
		case "ForOnePoly", "ForManyPoly":
			// For polymorphic "For" relationships, we need ID, type, and union fields
			polyFields, polyErr := getPolymorphicForTsFieldsForEntity(r, relationshipName, entityRelation, fieldCasing, hasAttribute(entityRelation.Attributes, "optional"))
			if polyErr != nil {
				return nil, polyErr
			}
			allFields = append(allFields, polyFields...)

		case "HasOnePoly", "HasManyPoly":
			// For polymorphic "Has" relationships, use the aliased entity if provided, otherwise use relationship name
			targetEntityName := relationshipName
			if entityRelation.Aliased != "" {
				targetEntityName = entityRelation.Aliased
			}

			targetEntityDef, targetEntityDefErr := r.GetEntity(targetEntityName)
			if targetEntityDefErr != nil {
				return nil, targetEntityDefErr
			}

			// Generate regular ID and object fields with the relationship name
			tsIDField, tsIDErr := getRelatedTsFieldForMorpheEntityPrimaryID(r, entityRelation.Type, relationshipName, targetEntityDef, fieldCasing, hasAttribute(entityRelation.Attributes, "optional"))
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheEntityObjectWithTargetName(entityRelation.Type, relationshipName, targetEntityName, fieldCasing)
			allFields = append(allFields, tsRelatedField)

		default:
			// Regular relationships
			targetEntityName := relationshipName
			if entityRelation.Aliased != "" {
				targetEntityName = entityRelation.Aliased
			}

			relatedEntityDef, relatedEntityDefErr := r.GetEntity(targetEntityName)
			if relatedEntityDefErr != nil {
				return nil, relatedEntityDefErr
			}

			tsIDField, tsIDErr := getRelatedTsFieldForMorpheEntityPrimaryID(r, entityRelation.Type, relationshipName, relatedEntityDef, fieldCasing, hasAttribute(entityRelation.Attributes, "optional"))
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheEntityObjectWithTargetName(entityRelation.Type, relationshipName, targetEntityName, fieldCasing)
			allFields = append(allFields, tsRelatedField)
		}
	}
	return allFields, nil
}

func getRelatedTsFieldForMorpheEntityPrimaryID(r *registry.Registry, relationType string, relatedEntityName string, relatedEntityDef yaml.Entity, fieldCasing cfg.Casing, isOptional bool) (tsdef.ObjectField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetEntityPrimaryIdentifierFieldName(relatedEntityDef)
	if relatedIDFieldNameErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}
	idFieldName := fieldCasing.Apply(fmt.Sprintf("%s%s", relatedEntityName, relatedPrimaryIDFieldName))

	relatedPrimaryIDFieldDef, relatedIDFieldDefErr := yamlops.GetEntityFieldDefinitionByName(relatedEntityDef, relatedPrimaryIDFieldName)
	if relatedIDFieldDefErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w (primary identifier)", relatedIDFieldDefErr)
	}
	idFieldType, typeErr := getTsTypeForEntityField(r, relatedPrimaryIDFieldDef)
	if typeErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w (primary identifier)", typeErr)
	}

	if yamlops.IsRelationMany(relationType) {
		arrayType := tsdef.TsTypeArray{ValueType: idFieldType}
		tsIDField := tsdef.ObjectField{
			Name: idFieldName + "s",
			Type: arrayType,
		}
		if isOptional {
			tsIDField.Type = tsdef.TsTypeOptional{ValueType: arrayType}
		}
		return tsIDField, nil
	}

	tsIDField := tsdef.ObjectField{
		Name: idFieldName,
		Type: idFieldType,
	}
	if isOptional {
		tsIDField.Type = tsdef.TsTypeOptional{ValueType: idFieldType}
	}
	return tsIDField, nil
}

func getRelatedTsFieldForMorpheEntityObject(relationType string, relatedEntityName string, isOptional bool) tsdef.ObjectField {
	objType := tsdef.TsTypeObject{
		ModulePath: "./" + strcase.ToKebabCaseLower(relatedEntityName),
		Name:       relatedEntityName,
	}
	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relatedEntityName + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{ValueType: objType},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relatedEntityName,
		Type: tsdef.TsTypeOptional{ValueType: objType},
	}
	return tsRelatedField
}

func getRelatedTsFieldForMorpheEntityObjectWithTargetName(relationType string, relationshipName string, targetEntityName string, fieldCasing cfg.Casing) tsdef.ObjectField {
	relationshipNameCamel := fieldCasing.Apply(relationshipName)
	objType := tsdef.TsTypeObject{
		ModulePath: "./" + strcase.ToKebabCaseLower(targetEntityName),
		Name:       targetEntityName,
	}

	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relationshipNameCamel + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{ValueType: objType},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relationshipNameCamel,
		Type: tsdef.TsTypeOptional{ValueType: objType},
	}
	return tsRelatedField
}

func getPolymorphicForTsFieldsForEntity(r *registry.Registry, relationshipName string, entityRelation yaml.EntityRelation, fieldCasing cfg.Casing, isOptional bool) ([]tsdef.ObjectField, error) {
	if len(entityRelation.For) == 0 {
		return nil, fmt.Errorf("polymorphic relation '%s' must have at least one entity in 'for' property", relationshipName)
	}

	relationshipNameCamel := fieldCasing.Apply(relationshipName)
	allFields := []tsdef.ObjectField{}

	wrapOptional := func(t tsdef.TsType) tsdef.TsType {
		if isOptional {
			return tsdef.TsTypeOptional{ValueType: t}
		}
		return t
	}

	// Add ID field(s)
	if yamlops.IsRelationMany(entityRelation.Type) {
		arrayType := tsdef.TsTypeArray{ValueType: tsdef.TsTypeString}
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel + "IDs",
			Type: wrapOptional(arrayType),
		})
	} else {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel + "ID",
			Type: wrapOptional(tsdef.TsTypeString),
		})
	}

	// Add type field
	allFields = append(allFields, tsdef.ObjectField{
		Name: relationshipNameCamel + "Type",
		Type: wrapOptional(tsdef.TsTypeString),
	})

	// Add union type field — always optional per ADR-003 (optionally loaded data)
	unionTypes := []tsdef.TsType{}
	for _, targetEntityName := range entityRelation.For {
		unionTypes = append(unionTypes, tsdef.TsTypeObject{
			ModulePath: "./" + strcase.ToKebabCaseLower(targetEntityName),
			Name:       targetEntityName,
		})
	}

	if yamlops.IsRelationMany(entityRelation.Type) {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeUnion{Types: unionTypes},
				},
			},
		})
	} else {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel,
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeUnion{Types: unionTypes},
			},
		})
	}

	return allFields, nil
}
