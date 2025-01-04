package compile

import (
	"strings"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/typemap"

	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

func MorpheEntityToTsObject(entityHooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, r *registry.Registry, entity yaml.Entity) (*tsdef.Object, error) {
	if r == nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, ErrNoRegistry)
	}

	config, entity, compileStartErr := triggerCompileMorpheEntityStart(entityHooks, config, entity)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileStartErr)
	}

	entityObject, objectErr := morpheEntityToTsObjectType(config, r, entity)
	if objectErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, objectErr)
	}

	entityObject, compileSuccessErr := triggerCompileMorpheEntitySuccess(entityHooks, entityObject)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileSuccessErr)
	}

	return entityObject, nil
}

func morpheEntityToTsObjectType(config cfg.MorpheEntitiesConfig, r *registry.Registry, entity yaml.Entity) (*tsdef.Object, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	// validateMorpheErr := entity.Validate()
	// if validateMorpheErr != nil {
	// 	return nil, validateMorpheErr
	// }

	entityType, entityTypeErr := getEntityObjectType(r, entity)
	if entityTypeErr != nil {
		return nil, entityTypeErr
	}

	return entityType, nil
}

func getEntityObjectType(r *registry.Registry, entity yaml.Entity) (*tsdef.Object, error) {
	entityType := tsdef.Object{
		Name: entity.Name,
	}

	typeFields, fieldsErr := getTsFieldsForMorpheEntity(r, entity.Fields)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	entityType.Fields = typeFields

	objectImports, importsErr := getImportsForObjectFields(typeFields)
	if importsErr != nil {
		return nil, importsErr
	}
	entityType.Imports = objectImports

	return &entityType, nil
}

func getTsFieldsForMorpheEntity(r *registry.Registry, entityFields map[string]yaml.EntityField) ([]tsdef.ObjectField, error) {
	allFieldNames := core.MapKeysSorted(entityFields)
	allTypeFields := make([]tsdef.ObjectField, 0, len(entityFields))

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
		allTypeFields = append(allTypeFields, typeField)
	}

	return allTypeFields, nil
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

	tsFieldType, typeSupported := typemap.MorpheModelFieldToTsField[terminalField.Type]
	if !typeSupported {
		return nil, ErrUnsupportedMorpheFieldType(terminalField.Type)
	}
	return tsFieldType, nil
}

func triggerCompileMorpheEntityStart(hooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, entity yaml.Entity) (cfg.MorpheEntitiesConfig, yaml.Entity, error) {
	if hooks.OnCompileMorpheEntityStart == nil {
		return config, entity, nil
	}

	return hooks.OnCompileMorpheEntityStart(config, entity)
}

func triggerCompileMorpheEntitySuccess(hooks hook.CompileMorpheEntity, entityObject *tsdef.Object) (*tsdef.Object, error) {
	if hooks.OnCompileMorpheEntitySuccess == nil {
		return entityObject, nil
	}

	return hooks.OnCompileMorpheEntitySuccess(entityObject)
}

func triggerCompileMorpheEntityFailure(hooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, entity yaml.Entity, failureErr error) error {
	if hooks.OnCompileMorpheEntityFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEntityFailure(config, entity, failureErr)
}
