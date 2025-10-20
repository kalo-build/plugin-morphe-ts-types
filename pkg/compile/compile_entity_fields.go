package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/typemap"
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
			Name: strcase.ToCamelCase(fieldName),
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
	for _, relationshipName := range allRelatedEntityNames {
		entityRelation := entityRelations[relationshipName]

		// Handle different relationship types
		switch entityRelation.Type {
		case "ForOnePoly", "ForManyPoly":
			// For polymorphic "For" relationships, we need ID, type, and union fields
			polyFields, polyErr := getPolymorphicForTsFieldsForEntity(r, relationshipName, entityRelation)
			if polyErr != nil {
				return nil, polyErr
			}
			allFields = append(allFields, polyFields...)

		case "HasOnePoly", "HasManyPoly":
			// For polymorphic "Has" relationships, use the aliased entity
			targetEntityName := entityRelation.Aliased
			if targetEntityName == "" {
				return nil, fmt.Errorf("polymorphic Has* relationship '%s' must specify 'aliased' property", relationshipName)
			}

			targetEntityDef, targetEntityDefErr := r.GetEntity(targetEntityName)
			if targetEntityDefErr != nil {
				return nil, targetEntityDefErr
			}

			// Generate regular ID and object fields with the relationship name
			tsIDField, tsIDErr := getRelatedTsFieldForMorpheEntityPrimaryID(r, entityRelation.Type, relationshipName, targetEntityDef)
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheEntityOptionalObjectWithTargetName(entityRelation.Type, relationshipName, targetEntityName)
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

			tsIDField, tsIDErr := getRelatedTsFieldForMorpheEntityPrimaryID(r, entityRelation.Type, relationshipName, relatedEntityDef)
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheEntityOptionalObjectWithTargetName(entityRelation.Type, relationshipName, targetEntityName)
			allFields = append(allFields, tsRelatedField)
		}
	}
	return allFields, nil
}

func getRelatedTsFieldForMorpheEntityPrimaryID(r *registry.Registry, relationType string, relatedEntityName string, relatedEntityDef yaml.Entity) (tsdef.ObjectField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetEntityPrimaryIdentifierFieldName(relatedEntityDef)
	if relatedIDFieldNameErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}
	idFieldName := strcase.ToCamelCase(fmt.Sprintf("%s%s", relatedEntityName, relatedPrimaryIDFieldName))

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

func getRelatedTsFieldForMorpheEntityOptionalObjectWithTargetName(relationType string, relationshipName string, targetEntityName string) tsdef.ObjectField {
	relationshipNameCamel := strcase.ToCamelCase(relationshipName)

	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relationshipNameCamel + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeObject{
						ModulePath: "./" + strcase.ToKebabCaseLower(targetEntityName),
						Name:       targetEntityName,
					},
				},
			},
		}
		return tsRelatedField
	}

	tsRelatedField := tsdef.ObjectField{
		Name: relationshipNameCamel,
		Type: tsdef.TsTypeOptional{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./" + strcase.ToKebabCaseLower(targetEntityName),
				Name:       targetEntityName,
			},
		},
	}
	return tsRelatedField
}

func getPolymorphicForTsFieldsForEntity(r *registry.Registry, relationshipName string, entityRelation yaml.EntityRelation) ([]tsdef.ObjectField, error) {
	if len(entityRelation.For) == 0 {
		return nil, fmt.Errorf("polymorphic relation '%s' must have at least one entity in 'for' property", relationshipName)
	}

	relationshipNameCamel := strcase.ToCamelCase(relationshipName)
	allFields := []tsdef.ObjectField{}

	// Add ID field(s)
	if yamlops.IsRelationMany(entityRelation.Type) {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel + "IDs",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeString,
				},
			},
		})
	} else {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel + "ID",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeString,
			},
		})
	}

	// Add type field
	allFields = append(allFields, tsdef.ObjectField{
		Name: relationshipNameCamel + "Type",
		Type: tsdef.TsTypeOptional{
			ValueType: tsdef.TsTypeString,
		},
	})

	// Add union type field
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
					ValueType: tsdef.TsTypeUnion{
						Types: unionTypes,
					},
				},
			},
		})
	} else {
		allFields = append(allFields, tsdef.ObjectField{
			Name: relationshipNameCamel,
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeUnion{
					Types: unionTypes,
				},
			},
		})
	}

	return allFields, nil
}
