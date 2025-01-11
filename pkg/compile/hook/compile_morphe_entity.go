package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

type CompileMorpheEntity struct {
	OnCompileMorpheEntityStart   OnCompileMorpheEntityStartHook
	OnCompileMorpheEntitySuccess OnCompileMorpheEntitySuccessHook
	OnCompileMorpheEntityFailure OnCompileMorpheEntityFailureHook
}

type OnCompileMorpheEntityStartHook = func(config cfg.MorpheEntitiesConfig, entity yaml.Entity) (cfg.MorpheEntitiesConfig, yaml.Entity, error)
type OnCompileMorpheEntitySuccessHook = func(entityObjects []*tsdef.Object) ([]*tsdef.Object, error)
type OnCompileMorpheEntityFailureHook = func(config cfg.MorpheEntitiesConfig, entity yaml.Entity, compileFailure error) error
