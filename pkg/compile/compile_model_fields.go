package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/typemap"
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
			Name: strcase.ToCamelCase(fieldName),
			Type: tsFieldType,
		}
		allFields = append(allFields, tsField)
	}

	return allFields, nil
}

func getRelatedTsFieldsForMorpheModel(r *registry.Registry, modelRelations map[string]yaml.ModelRelation) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}

	allRelatedModelNames := core.MapKeysSorted(modelRelations)
	for _, relationshipName := range allRelatedModelNames {
		modelRelation := modelRelations[relationshipName]

		// Handle different relationship types
		switch modelRelation.Type {
		case "ForOnePoly", "ForManyPoly":
			// For polymorphic "For" relationships, we need ID, type, and union fields
			polyFields, polyErr := getPolymorphicForTsFields(r, relationshipName, modelRelation)
			if polyErr != nil {
				return nil, polyErr
			}
			allFields = append(allFields, polyFields...)

		case "HasOnePoly", "HasManyPoly":
			// For polymorphic "Has" relationships, use the aliased model
			targetModelName := modelRelation.Aliased
			if targetModelName == "" {
				return nil, fmt.Errorf("polymorphic Has* relationship '%s' must specify 'aliased' property", relationshipName)
			}

			targetModelDef, targetModelDefErr := r.GetModel(targetModelName)
			if targetModelDefErr != nil {
				return nil, targetModelDefErr
			}

			// Generate regular ID and object fields with the relationship name
			tsIDField, tsIDErr := getRelatedTsFieldForMorpheModelPrimaryID(modelRelation.Type, relationshipName, targetModelDef)
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheModelOptionalObjectWithTargetName(modelRelation.Type, relationshipName, targetModelName)
			allFields = append(allFields, tsRelatedField)

		default:
			// Regular relationships
			targetModelName := relationshipName
			if modelRelation.Aliased != "" {
				targetModelName = modelRelation.Aliased
			}

			targetModelDef, targetModelDefErr := r.GetModel(targetModelName)
			if targetModelDefErr != nil {
				return nil, targetModelDefErr
			}

			tsIDField, tsIDErr := getRelatedTsFieldForMorpheModelPrimaryID(modelRelation.Type, relationshipName, targetModelDef)
			if tsIDErr != nil {
				return nil, tsIDErr
			}
			allFields = append(allFields, tsIDField)

			tsRelatedField := getRelatedTsFieldForMorpheModelOptionalObjectWithTargetName(modelRelation.Type, relationshipName, targetModelName)
			allFields = append(allFields, tsRelatedField)
		}
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
		Name: strcase.ToCamelCase(fieldName),
		Type: tsFieldType,
	}
	return tsField
}

func getRelatedTsFieldForMorpheModelPrimaryID(relationType string, relatedModelName string, relatedModelDef yaml.Model) (tsdef.ObjectField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetModelPrimaryIdentifierFieldName(relatedModelDef)
	if relatedIDFieldNameErr != nil {
		return tsdef.ObjectField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}
	idFieldName := strcase.ToCamelCase(fmt.Sprintf("%s%s", relatedModelName, relatedPrimaryIDFieldName))

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

func getRelatedTsFieldForMorpheModelOptionalObjectWithTargetName(relationType string, relationshipName string, targetModelName string) tsdef.ObjectField {
	relationshipNameCamel := strcase.ToCamelCase(relationshipName)

	if yamlops.IsRelationMany(relationType) {
		tsRelatedField := tsdef.ObjectField{
			Name: relationshipNameCamel + "s",
			Type: tsdef.TsTypeOptional{
				ValueType: tsdef.TsTypeArray{
					ValueType: tsdef.TsTypeObject{
						ModulePath: "./" + strcase.ToKebabCaseLower(targetModelName),
						Name:       targetModelName,
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
				ModulePath: "./" + strcase.ToKebabCaseLower(targetModelName),
				Name:       targetModelName,
			},
		},
	}
	return tsRelatedField
}

func getPolymorphicForTsFields(r *registry.Registry, relationshipName string, modelRelation yaml.ModelRelation) ([]tsdef.ObjectField, error) {
	if len(modelRelation.For) == 0 {
		return nil, fmt.Errorf("polymorphic relation '%s' must have at least one model in 'for' property", relationshipName)
	}

	relationshipNameCamel := strcase.ToCamelCase(relationshipName)
	allFields := []tsdef.ObjectField{}

	// Add ID field(s)
	if yamlops.IsRelationMany(modelRelation.Type) {
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
	for _, targetModelName := range modelRelation.For {
		unionTypes = append(unionTypes, tsdef.TsTypeObject{
			ModulePath: "./" + strcase.ToKebabCaseLower(targetModelName),
			Name:       targetModelName,
		})
	}

	if yamlops.IsRelationMany(modelRelation.Type) {
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
