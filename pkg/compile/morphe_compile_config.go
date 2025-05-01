package compile

import (
	"path"

	r "github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/write"
)

type MorpheCompileConfig struct {
	rcfg.MorpheLoadRegistryConfig
	cfg.MorpheModelsConfig
	cfg.MorpheEnumsConfig
	cfg.MorpheStructuresConfig
	cfg.MorpheEntitiesConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	EnumWriter write.TsEnumWriter
	EnumHooks  hook.CompileMorpheEnum

	ModelWriter write.TsObjectWriter
	ModelHooks  hook.CompileMorpheModel

	EntityWriter write.TsObjectWriter
	EntityHooks  hook.CompileMorpheEntity

	WriteObjectHooks hook.WriteTsObject
	WriteEnumHooks   hook.WriteTsEnum

	StructureWriter write.TsObjectWriter
	StructureHooks  hook.CompileMorpheStructure
}

func DefaultMorpheCompileConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) MorpheCompileConfig {
	return MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		MorpheModelsConfig:     cfg.MorpheModelsConfig{},
		MorpheEnumsConfig:      cfg.MorpheEnumsConfig{},
		MorpheStructuresConfig: cfg.MorpheStructuresConfig{},
		MorpheEntitiesConfig:   cfg.MorpheEntitiesConfig{},

		RegistryHooks: r.LoadMorpheRegistryHooks{},

		EnumWriter: &MorpheEnumFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "enums"),
		},
		EnumHooks: hook.CompileMorpheEnum{},

		ModelWriter: &MorpheObjectFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "models"),
		},
		ModelHooks: hook.CompileMorpheModel{},

		EntityWriter: &MorpheObjectFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "entities"),
		},
		EntityHooks: hook.CompileMorpheEntity{},

		WriteObjectHooks: hook.WriteTsObject{},
		WriteEnumHooks:   hook.WriteTsEnum{},

		StructureWriter: &MorpheObjectFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "structures"),
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}
}
