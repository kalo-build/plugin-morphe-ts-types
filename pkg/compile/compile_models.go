package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

func AllMorpheModelsToTsObjects(config MorpheCompileConfig, r *registry.Registry) (map[string][]*tsdef.Object, error) {
	allModelTypeDefs := map[string][]*tsdef.Object{}
	modelsConfig := config.MorpheModelsConfig
	modelsConfig.FieldCasing = config.FieldCasing
	for modelName, model := range r.GetAllModels() {
		modelTypes, modelErr := MorpheModelToTsObjects(config.ModelHooks, modelsConfig, r, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelTypeDefs[modelName] = modelTypes
	}
	return allModelTypeDefs, nil
}

func MorpheModelToTsObjects(modelHooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, r *registry.Registry, model yaml.Model) ([]*tsdef.Object, error) {
	if r == nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, ErrNoRegistry)
	}
	config, model, compileStartErr := triggerCompileMorpheModelStart(modelHooks, config, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileStartErr)
	}
	allModelTypes, objectsErr := morpheModelToTsObjectTypes(config, config.FieldCasing, r, model)
	if objectsErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, objectsErr)
	}

	allModelTypes, compileSuccessErr := triggerCompileMorpheModelSuccess(modelHooks, allModelTypes)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileSuccessErr)
	}
	return allModelTypes, nil
}

func morpheModelToTsObjectTypes(config cfg.MorpheModelsConfig, fieldCasing cfg.Casing, r *registry.Registry, model yaml.Model) ([]*tsdef.Object, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := model.Validate(r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	modelType, modelTypeErr := getModelObjectType(r, model, fieldCasing)
	if modelTypeErr != nil {
		return nil, modelTypeErr
	}
	allIdentifierTypes, identifierTypesErr := getAllModelIdentifierObjectTypes(model, modelType, fieldCasing)
	if identifierTypesErr != nil {
		return nil, identifierTypesErr
	}

	allModelTypes := []*tsdef.Object{
		modelType,
	}
	allModelTypes = append(allModelTypes, allIdentifierTypes...)
	return allModelTypes, nil
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

func getModelObjectType(r *registry.Registry, model yaml.Model, fieldCasing cfg.Casing) (*tsdef.Object, error) {
	modelType := tsdef.Object{
		Name: model.Name,
	}
	typeFields, fieldsErr := getTsFieldsForMorpheModel(r, model.Fields, model.Related, fieldCasing)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	modelType.Fields = typeFields

	objectImports, importsErr := getImportsForObjectFields(typeFields)
	if importsErr != nil {
		return nil, importsErr
	}
	modelType.Imports = objectImports
	return &modelType, nil
}

func getAllModelIdentifierObjectTypes(model yaml.Model, modelType *tsdef.Object, fieldCasing cfg.Casing) ([]*tsdef.Object, error) {
	modelIdentifiers := model.Identifiers
	allIdentifierNames := core.MapKeysSorted(modelIdentifiers)
	allIdentTypes := []*tsdef.Object{}
	for _, identifierName := range allIdentifierNames {
		identifierDef := modelIdentifiers[identifierName]

		allIdentFieldDefs, identFieldDefsErr := getModelIdentifierObjectFieldSubset(*modelType, identifierName, identifierDef, fieldCasing)
		if identFieldDefsErr != nil {
			return nil, identFieldDefsErr
		}

		identObject, identObjectErr := getModelIdentifierObjectType(modelType.Name, identifierName, allIdentFieldDefs)
		if identObjectErr != nil {
			return nil, identObjectErr
		}
		allIdentTypes = append(allIdentTypes, identObject)
	}
	return allIdentTypes, nil
}

func getModelIdentifierObjectType(modelName string, identifierName string, allIdentFieldDefs []tsdef.ObjectField) (*tsdef.Object, error) {
	identifierType := tsdef.Object{
		Name:   fmt.Sprintf("%sID%s", modelName, strcase.ToPascalCase(identifierName)),
		Fields: allIdentFieldDefs,
	}
	return &identifierType, nil
}

func getModelIdentifierObjectFieldSubset(modelType tsdef.Object, identifierName string, identifier yaml.ModelIdentifier, fieldCasing cfg.Casing) ([]tsdef.ObjectField, error) {
	identifierFieldDefs := []tsdef.ObjectField{}
	for _, fieldName := range identifier.Fields {
		if strings.HasPrefix(fieldName, "rel:") {
			relationName := strings.TrimPrefix(fieldName, "rel:")
			idFieldName := fieldCasing.Apply(relationName + "ID")
			idField, found := findObjectFieldByName(modelType.Fields, idFieldName)
			if !found {
				return nil, fmt.Errorf("identifier %s references relation '%s' but field '%s' not found on type", identifierName, relationName, idFieldName)
			}
			identifierFieldDefs = append(identifierFieldDefs, idField)
			typeFieldName := fieldCasing.Apply(relationName + "Type")
			if typeField, hasType := findObjectFieldByName(modelType.Fields, typeFieldName); hasType {
				identifierFieldDefs = append(identifierFieldDefs, typeField)
			}
			continue
		}
		identifierFieldDef := tsdef.ObjectField{}
		tsFieldName := fieldCasing.Apply(fieldName)
		for _, modelFieldDef := range modelType.Fields {
			if modelFieldDef.Name != tsFieldName {
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

func findObjectFieldByName(fields []tsdef.ObjectField, name string) (tsdef.ObjectField, bool) {
	for _, field := range fields {
		if field.Name == name {
			return tsdef.ObjectField{Name: field.Name, Type: field.Type}, true
		}
	}
	return tsdef.ObjectField{}, false
}

func getImportsForObjectFields(allFields []tsdef.ObjectField) ([]tsdef.ObjectImport, error) {
	objectImportMap := map[string]tsdef.ObjectImport{}
	for _, fieldDef := range allFields {
		allFieldImports := fieldDef.Type.GetImports()
		for _, fieldImport := range allFieldImports {
			objectImportMap[fieldImport.ModulePath] = fieldImport
		}
	}

	allModulePaths := core.MapKeysSorted(objectImportMap)

	allObjectImports := []tsdef.ObjectImport{}
	for _, modulePath := range allModulePaths {
		allObjectImports = append(allObjectImports, objectImportMap[modulePath])
	}
	return allObjectImports, nil
}
