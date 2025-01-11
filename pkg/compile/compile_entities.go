package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"

	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

func AllMorpheEntitiesToTsObjects(config MorpheCompileConfig, r *registry.Registry) (map[string][]*tsdef.Object, error) {
	allEntityTypeDefs := map[string][]*tsdef.Object{}
	for entityName, entity := range r.GetAllEntities() {
		entityTypes, entityTypesErr := MorpheEntityToTsObjects(config.EntityHooks, config.MorpheEntitiesConfig, r, entity)
		if entityTypesErr != nil {
			return nil, entityTypesErr
		}
		allEntityTypeDefs[entityName] = entityTypes
	}
	return allEntityTypeDefs, nil
}

func MorpheEntityToTsObjects(entityHooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, r *registry.Registry, entity yaml.Entity) ([]*tsdef.Object, error) {
	if r == nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, ErrNoRegistry)
	}

	config, entity, compileStartErr := triggerCompileMorpheEntityStart(entityHooks, config, entity)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileStartErr)
	}

	allEntityTypes, objectsErr := morpheEntityToTsObjectTypes(config, r, entity)
	if objectsErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, objectsErr)
	}

	allEntityTypes, compileSuccessErr := triggerCompileMorpheEntitySuccess(entityHooks, allEntityTypes)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileSuccessErr)
	}

	return allEntityTypes, nil
}

func morpheEntityToTsObjectTypes(config cfg.MorpheEntitiesConfig, r *registry.Registry, entity yaml.Entity) ([]*tsdef.Object, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := entity.Validate(r.GetAllModels(), r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	entityType, entityTypeErr := getEntityObjectType(r, entity)
	if entityTypeErr != nil {
		return nil, entityTypeErr
	}

	allIdentifierTypes, identifierTypesErr := getAllEntityIdentifierObjectTypes(entity, entityType)
	if identifierTypesErr != nil {
		return nil, identifierTypesErr
	}

	allEntityTypes := []*tsdef.Object{
		entityType,
	}
	allEntityTypes = append(allEntityTypes, allIdentifierTypes...)
	return allEntityTypes, nil
}

func getAllEntityIdentifierObjectTypes(entity yaml.Entity, entityType *tsdef.Object) ([]*tsdef.Object, error) {
	entityIdentifiers := entity.Identifiers
	allIdentifierNames := core.MapKeysSorted(entityIdentifiers)
	allIdentTypes := []*tsdef.Object{}

	for _, identifierName := range allIdentifierNames {
		identifierDef := entityIdentifiers[identifierName]

		allIdentFieldDefs, identFieldDefsErr := getEntityIdentifierObjectFieldSubset(*entityType, identifierName, identifierDef)
		if identFieldDefsErr != nil {
			return nil, identFieldDefsErr
		}

		identObject, identObjectErr := getEntityIdentifierObjectType(entityType.Name, identifierName, allIdentFieldDefs)
		if identObjectErr != nil {
			return nil, identObjectErr
		}
		allIdentTypes = append(allIdentTypes, identObject)
	}
	return allIdentTypes, nil
}

func getEntityObjectType(r *registry.Registry, entity yaml.Entity) (*tsdef.Object, error) {
	entityType := tsdef.Object{
		Name: entity.Name,
	}

	typeFields, fieldsErr := getTsFieldsForMorpheEntity(r, entity.Fields, entity.Related)
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

func getEntityIdentifierObjectType(entityName string, identifierName string, allIdentFieldDefs []tsdef.ObjectField) (*tsdef.Object, error) {
	identifierType := tsdef.Object{
		Name:   fmt.Sprintf("%sID%s", entityName, strcase.ToPascalCase(identifierName)),
		Fields: allIdentFieldDefs,
	}
	return &identifierType, nil
}

func getEntityIdentifierObjectFieldSubset(entityType tsdef.Object, identifierName string, identifier yaml.EntityIdentifier) ([]tsdef.ObjectField, error) {
	identifierFieldDefs := []tsdef.ObjectField{}
	for _, fieldName := range identifier.Fields {
		identifierFieldDef := tsdef.ObjectField{}
		for _, entityFieldDef := range entityType.Fields {
			if entityFieldDef.Name != fieldName {
				continue
			}
			identifierFieldDef = tsdef.ObjectField{
				Name: entityFieldDef.Name,
				Type: entityFieldDef.Type,
			}
		}
		if identifierFieldDef.Name == "" {
			return nil, ErrMissingMorpheIdentifierField(entityType.Name, identifierName, fieldName)
		}
		identifierFieldDefs = append(identifierFieldDefs, identifierFieldDef)
	}
	return identifierFieldDefs, nil
}

func triggerCompileMorpheEntityStart(hooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, entity yaml.Entity) (cfg.MorpheEntitiesConfig, yaml.Entity, error) {
	if hooks.OnCompileMorpheEntityStart == nil {
		return config, entity, nil
	}

	return hooks.OnCompileMorpheEntityStart(config, entity)
}

func triggerCompileMorpheEntitySuccess(hooks hook.CompileMorpheEntity, entityObjects []*tsdef.Object) ([]*tsdef.Object, error) {
	if hooks.OnCompileMorpheEntitySuccess == nil {
		return entityObjects, nil
	}

	return hooks.OnCompileMorpheEntitySuccess(entityObjects)
}

func triggerCompileMorpheEntityFailure(hooks hook.CompileMorpheEntity, config cfg.MorpheEntitiesConfig, entity yaml.Entity, failureErr error) error {
	if hooks.OnCompileMorpheEntityFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEntityFailure(config, entity, failureErr)
}
