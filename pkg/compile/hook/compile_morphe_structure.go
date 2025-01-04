package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

type CompileMorpheStructure struct {
	OnCompileMorpheStructureStart   OnCompileMorpheStructureStartHook
	OnCompileMorpheStructureSuccess OnCompileMorpheStructureSuccessHook
	OnCompileMorpheStructureFailure OnCompileMorpheStructureFailureHook
}

type OnCompileMorpheStructureStartHook = func(config cfg.MorpheStructuresConfig, structure yaml.Structure) (cfg.MorpheStructuresConfig, yaml.Structure, error)
type OnCompileMorpheStructureSuccessHook = func(structureType *tsdef.Object) (*tsdef.Object, error)
type OnCompileMorpheStructureFailureHook = func(config cfg.MorpheStructuresConfig, structure yaml.Structure, compileFailure error) error
