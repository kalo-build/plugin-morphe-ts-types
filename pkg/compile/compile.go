package compile

import "github.com/kalo-build/morphe-go/pkg/registry"

func MorpheToTypescript(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	hasEnums := r.HasEnums()
	if hasEnums {
		allEnumDefs, compileAllEnumsErr := AllMorpheEnumsToTsEnums(config, r)
		if compileAllEnumsErr != nil {
			return compileAllEnumsErr
		}

		_, writeAllEnumsErr := WriteAllEnumDefinitions(config, allEnumDefs)
		if writeAllEnumsErr != nil {
			return writeAllEnumsErr
		}
	}

	hasModels := r.HasModels()
	if hasModels {
		allModelObjectDefs, compileAllModelsErr := AllMorpheModelsToTsObjects(config, r)
		if compileAllModelsErr != nil {
			return compileAllModelsErr
		}

		_, writeAllModelsErr := WriteAllModelObjectDefinitions(config, allModelObjectDefs)
		if writeAllModelsErr != nil {
			return writeAllModelsErr
		}
	}

	hasStructures := r.HasStructures()
	if hasStructures {
		allStructureObjectDefs, compileAllStructuresErr := AllMorpheStructuresToTsObjects(config, r)
		if compileAllStructuresErr != nil {
			return compileAllStructuresErr
		}

		_, writeAllStructuresErr := WriteAllStructureObjectDefinitions(config, allStructureObjectDefs)
		if writeAllStructuresErr != nil {
			return writeAllStructuresErr
		}
	}

	hasEntities := r.HasEntities()
	if hasEntities {
		allEntityObjectDefs, compileAllEntitiesErr := AllMorpheEntitiesToTsObjects(config, r)
		if compileAllEntitiesErr != nil {
			return compileAllEntitiesErr
		}

		_, writeAllEntitiesErr := WriteAllEntityObjectDefinitions(config, allEntityObjectDefs)
		if writeAllEntitiesErr != nil {
			return writeAllEntitiesErr
		}
	}

	return nil
}
