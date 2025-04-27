package hook

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

type CompileMorpheEnum struct {
	OnCompileMorpheEnumStart   OnCompileMorpheEnumStartHook
	OnCompileMorpheEnumSuccess OnCompileMorpheEnumSuccessHook
	OnCompileMorpheEnumFailure OnCompileMorpheEnumFailureHook
}

type OnCompileMorpheEnumStartHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error)
type OnCompileMorpheEnumSuccessHook = func(enumType *tsdef.Enum) (*tsdef.Enum, error)
type OnCompileMorpheEnumFailureHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum, compileFailure error) error
