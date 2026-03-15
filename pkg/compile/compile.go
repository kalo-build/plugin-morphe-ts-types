package compile

import (
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
)

func MorpheToTypescript(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	var barrelCategories []BarrelFileCategory

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

		if config.GenerateBarrelFiles {
			barrelCategories = append(barrelCategories, BarrelFileCategory{
				DirPath:         config.EnumWriter.(*MorpheEnumFileWriter).TargetDirPath,
				DefinitionNames: core.MapKeysSorted(allEnumDefs),
			})
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

		if config.GenerateBarrelFiles {
			barrelCategories = append(barrelCategories, BarrelFileCategory{
				DirPath:         config.ModelWriter.(*MorpheObjectFileWriter).TargetDirPath,
				DefinitionNames: core.MapKeysSorted(allModelObjectDefs),
			})
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

		if config.GenerateBarrelFiles {
			barrelCategories = append(barrelCategories, BarrelFileCategory{
				DirPath:         config.StructureWriter.(*MorpheObjectFileWriter).TargetDirPath,
				DefinitionNames: core.MapKeysSorted(allStructureObjectDefs),
			})
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

		if config.GenerateBarrelFiles {
			barrelCategories = append(barrelCategories, BarrelFileCategory{
				DirPath:         config.EntityWriter.(*MorpheObjectFileWriter).TargetDirPath,
				DefinitionNames: core.MapKeysSorted(allEntityObjectDefs),
			})
		}
	}

	if config.GenerateBarrelFiles && len(barrelCategories) > 0 {
		if err := WriteBarrelFiles(barrelCategories); err != nil {
			return err
		}
	}

	return nil
}
