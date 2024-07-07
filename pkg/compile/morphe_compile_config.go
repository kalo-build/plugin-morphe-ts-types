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

	RegistryHooks r.LoadMorpheRegistryHooks

	ModelWriter write.TypescriptTypeWriter
	ModelHooks  hook.CompileMorpheModel

	WriteTypeHooks hook.WriteTsType
}
