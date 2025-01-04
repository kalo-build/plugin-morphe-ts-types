package compile

import (
	r "github.com/kaloseia/morphe-go/pkg/registry"
	rcfg "github.com/kaloseia/morphe-go/pkg/registry/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/write"
)

type MorpheCompileConfig struct {
	rcfg.MorpheLoadRegistryConfig
	cfg.MorpheModelsConfig
	cfg.MorpheEnumsConfig
	cfg.MorpheStructuresConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	EnumWriter write.TsEnumWriter
	EnumHooks  hook.CompileMorpheEnum

	ModelWriter write.TsObjectWriter
	ModelHooks  hook.CompileMorpheModel

	WriteObjectHooks hook.WriteTsObject
	WriteEnumHooks   hook.WriteTsEnum

	StructureWriter write.TsObjectWriter
	StructureHooks  hook.CompileMorpheStructure
}
