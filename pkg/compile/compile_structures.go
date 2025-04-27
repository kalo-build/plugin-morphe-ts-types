package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

func AllMorpheStructuresToTsObjects(config MorpheCompileConfig, r *registry.Registry) (map[string]*tsdef.Object, error) {
	allStructureTypeDefs := map[string]*tsdef.Object{}
	for structureName, structure := range r.GetAllStructures() {
		structureType, structureErr := MorpheStructureToTsObject(config.StructureHooks, config.MorpheStructuresConfig, r, structure)
		if structureErr != nil {
			return nil, structureErr
		}
		allStructureTypeDefs[structureName] = structureType
	}
	return allStructureTypeDefs, nil
}

func MorpheStructureToTsObject(structureHooks hook.CompileMorpheStructure, config cfg.MorpheStructuresConfig, r *registry.Registry, structure yaml.Structure) (*tsdef.Object, error) {
	if r == nil {
		return nil, triggerCompileMorpheStructureFailure(structureHooks, config, structure, ErrNoRegistry)
	}
	config, structure, compileStartErr := triggerCompileMorpheStructureStart(structureHooks, config, structure)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheStructureFailure(structureHooks, config, structure, compileStartErr)
	}

	structureType, objectErr := morpheStructureToTsObjectType(config, r, structure)
	if objectErr != nil {
		return nil, triggerCompileMorpheStructureFailure(structureHooks, config, structure, objectErr)
	}

	structureType, compileSuccessErr := triggerCompileMorpheStructureSuccess(structureHooks, structureType)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheStructureFailure(structureHooks, config, structure, compileSuccessErr)
	}
	return structureType, nil
}

func morpheStructureToTsObjectType(config cfg.MorpheStructuresConfig, r *registry.Registry, structure yaml.Structure) (*tsdef.Object, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := structure.Validate(r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	structureType := tsdef.Object{
		Name: structure.Name,
	}

	typeFields, fieldsErr := getTsFieldsForMorpheStructure(r, structure.Fields)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	structureType.Fields = typeFields

	objectImports, importsErr := getImportsForObjectFields(typeFields)
	if importsErr != nil {
		return nil, importsErr
	}
	structureType.Imports = objectImports

	return &structureType, nil
}

func triggerCompileMorpheStructureStart(hooks hook.CompileMorpheStructure, config cfg.MorpheStructuresConfig, structure yaml.Structure) (cfg.MorpheStructuresConfig, yaml.Structure, error) {
	if hooks.OnCompileMorpheStructureStart == nil {
		return config, structure, nil
	}
	return hooks.OnCompileMorpheStructureStart(config, structure)
}

func triggerCompileMorpheStructureSuccess(hooks hook.CompileMorpheStructure, structureType *tsdef.Object) (*tsdef.Object, error) {
	if hooks.OnCompileMorpheStructureSuccess == nil {
		return structureType, nil
	}
	return hooks.OnCompileMorpheStructureSuccess(structureType)
}

func triggerCompileMorpheStructureFailure(hooks hook.CompileMorpheStructure, config cfg.MorpheStructuresConfig, structure yaml.Structure, compileErr error) error {
	if hooks.OnCompileMorpheStructureFailure == nil {
		return compileErr
	}
	return hooks.OnCompileMorpheStructureFailure(config, structure, compileErr)
}
