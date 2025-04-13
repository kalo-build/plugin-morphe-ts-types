package compile

import "github.com/kalo-build/morphe-go/pkg/registry"

func MorpheToTypescript(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	allEnumDefs, compileAllEnumsErr := AllMorpheEnumsToTsEnums(config, r)
	if compileAllEnumsErr != nil {
		return compileAllEnumsErr
	}

	_, writeAllEnumsErr := WriteAllEnumDefinitions(config, allEnumDefs)
	if writeAllEnumsErr != nil {
		return writeAllEnumsErr
	}

	allModelObjectDefs, compileAllModelsErr := AllMorpheModelsToTsObjects(config, r)
	if compileAllModelsErr != nil {
		return compileAllModelsErr
	}

	_, writeAllModelsErr := WriteAllModelObjectDefinitions(config, allModelObjectDefs)
	if writeAllModelsErr != nil {
		return writeAllModelsErr
	}

	allStructureObjectDefs, compileAllStructuresErr := AllMorpheStructuresToTsObjects(config, r)
	if compileAllStructuresErr != nil {
		return compileAllStructuresErr
	}

	_, writeAllStructuresErr := WriteAllStructureObjectDefinitions(config, allStructureObjectDefs)
	if writeAllStructuresErr != nil {
		return writeAllStructuresErr
	}

	allEntityObjectDefs, compileAllEntitiesErr := AllMorpheEntitiesToTsObjects(config, r)
	if compileAllEntitiesErr != nil {
		return compileAllEntitiesErr
	}

	_, writeAllEntitiesErr := WriteAllEntityObjectDefinitions(config, allEntityObjectDefs)
	if writeAllEntitiesErr != nil {
		return writeAllEntitiesErr
	}

	return nil
}
