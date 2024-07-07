package compile

import (
	"fmt"

	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/typemap"
)

func AllMorpheModelsToTsObjects(config MorpheCompileConfig, r *registry.Registry) (map[string][]*tsdef.Object, error) {
	allModelTypeDefs := map[string][]*tsdef.Object{}
	for modelName, model := range r.GetAllModels() {
		modelTypes, modelErr := MorpheModelToTsObjects(config.ModelHooks, config.MorpheModelsConfig, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelTypeDefs[modelName] = modelTypes
	}
	return allModelTypeDefs, nil
}

func MorpheModelToTsObjects(modelHooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model) ([]*tsdef.Object, error) {
	config, model, compileStartErr := triggerCompileMorpheModelStart(modelHooks, config, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileStartErr)
	}
	allModelTypes, structsErr := morpheModelToTsObjectTypes(config, model)
	if structsErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, structsErr)
	}

	allModelTypes, compileSuccessErr := triggerCompileMorpheModelSuccess(modelHooks, allModelTypes)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileSuccessErr)
	}
	return allModelTypes, nil
}

func morpheModelToTsObjectTypes(config cfg.MorpheModelsConfig, model yaml.Model) ([]*tsdef.Object, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := model.Validate()
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	modelType, modelTypeErr := getModelObjectType(model)
	if modelTypeErr != nil {
		return nil, modelTypeErr
	}
	allModelTypes := []*tsdef.Object{
		modelType,
	}

	modelIdentifiers := model.Identifiers
	allIdentifierNames := core.MapKeysSorted(modelIdentifiers)
	for _, identifierName := range allIdentifierNames {
		identifierDef := modelIdentifiers[identifierName]

		allIdentFieldDefs, identFieldDefsErr := getIdentifierObjectFieldSubset(*modelType, identifierName, identifierDef)
		if identFieldDefsErr != nil {
			return nil, identFieldDefsErr
		}

		identObject, identObjectErr := getIdentifierObjectType(modelType.Name, identifierName, allIdentFieldDefs)
		if identObjectErr != nil {
			return nil, identObjectErr
		}
		allModelTypes = append(allModelTypes, identObject)
	}
	return allModelTypes, nil
}

func getIdentifierObjectFieldSubset(modelType tsdef.Object, identifierName string, identifier yaml.ModelIdentifier) ([]tsdef.ObjectField, error) {
	identifierFieldDefs := []tsdef.ObjectField{}
	for _, fieldName := range identifier.Fields {
		identifierFieldDef := tsdef.ObjectField{}
		for _, modelFieldDef := range modelType.Fields {
			if modelFieldDef.Name != fieldName {
				continue
			}
			identifierFieldDef = tsdef.ObjectField{
				Name: modelFieldDef.Name,
				Type: modelFieldDef.Type,
			}
		}
		if identifierFieldDef.Name == "" {
			return nil, ErrMissingMorpheIdentifierField(modelType.Name, identifierName, fieldName)
		}
		identifierFieldDefs = append(identifierFieldDefs, identifierFieldDef)
	}
	return identifierFieldDefs, nil
}

func getIdentifierObjectType(modelName string, identifierName string, allIdentFieldDefs []tsdef.ObjectField) (*tsdef.Object, error) {
	identifierType := tsdef.Object{
		Name:   fmt.Sprintf("%sID%s", modelName, strcase.ToPascalCase(identifierName)),
		Fields: allIdentFieldDefs,
	}
	return &identifierType, nil
}

func getModelObjectType(model yaml.Model) (*tsdef.Object, error) {
	modelType := tsdef.Object{
		Name: model.Name,
	}
	typeFields, fieldsErr := getTsFieldsForMorpheModel(model.Fields)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	modelType.Fields = typeFields
	return &modelType, nil
}

func getTsFieldsForMorpheModel(modelFields map[string]yaml.ModelField) ([]tsdef.ObjectField, error) {
	allFields := []tsdef.ObjectField{}

	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]
		goFieldType, typeSupported := typemap.MorpheFieldToTsField[fieldDef.Type]
		if !typeSupported {
			return nil, ErrUnsupportedMorpheFieldType(fieldDef.Type)
		}
		goField := tsdef.ObjectField{
			Name: fieldName,
			Type: goFieldType,
		}
		allFields = append(allFields, goField)
	}
	return allFields, nil
}

func triggerCompileMorpheModelStart(hooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model) (cfg.MorpheModelsConfig, yaml.Model, error) {
	if hooks.OnCompileMorpheModelStart == nil {
		return config, model, nil
	}

	updatedConfig, updatedModel, startErr := hooks.OnCompileMorpheModelStart(config, model)
	if startErr != nil {
		return cfg.MorpheModelsConfig{}, yaml.Model{}, startErr
	}

	return updatedConfig, updatedModel, nil
}

func triggerCompileMorpheModelSuccess(hooks hook.CompileMorpheModel, allModelObjects []*tsdef.Object) ([]*tsdef.Object, error) {
	if hooks.OnCompileMorpheModelSuccess == nil {
		return allModelObjects, nil
	}
	if allModelObjects == nil {
		return nil, ErrNoModelObjects
	}
	allModelObjectsClone := clone.DeepCloneSlicePointers(allModelObjects)

	allModelObjects, successErr := hooks.OnCompileMorpheModelSuccess(allModelObjectsClone)
	if successErr != nil {
		return nil, successErr
	}
	return allModelObjects, nil
}

func triggerCompileMorpheModelFailure(hooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model, failureErr error) error {
	if hooks.OnCompileMorpheModelFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheModelFailure(config, model.DeepClone(), failureErr)
}
