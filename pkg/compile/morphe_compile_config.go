package compile

import (
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
