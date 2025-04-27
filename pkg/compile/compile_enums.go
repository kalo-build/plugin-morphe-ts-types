package compile

import (
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

func AllMorpheEnumsToTsEnums(config MorpheCompileConfig, r *registry.Registry) (map[string]*tsdef.Enum, error) {
	allEnumTypeDefs := map[string]*tsdef.Enum{}
	for enumName, enum := range r.GetAllEnums() {
		enumType, enumErr := MorpheEnumToTsEnum(config.EnumHooks, config.MorpheEnumsConfig, enum)
		if enumErr != nil {
			return nil, enumErr
		}
		allEnumTypeDefs[enumName] = enumType
	}
	return allEnumTypeDefs, nil
}

func MorpheEnumToTsEnum(enumHooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum) (*tsdef.Enum, error) {
	config, enum, compileStartErr := triggerCompileMorpheEnumStart(enumHooks, config, enum)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, compileStartErr)
	}
	enumType, structsErr := morpheEnumToTsEnumType(config, enum)
	if structsErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, structsErr)
	}

	enumType, compileSuccessErr := triggerCompileMorpheEnumSuccess(enumHooks, enumType)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, compileSuccessErr)
	}
	return enumType, nil
}

func morpheEnumToTsEnumType(config cfg.MorpheEnumsConfig, enum yaml.Enum) (*tsdef.Enum, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := enum.Validate()
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	enumType, enumTypeErr := getTypescriptEnum(enum)
	if enumTypeErr != nil {
		return nil, enumTypeErr
	}

	return enumType, nil
}

func triggerCompileMorpheEnumStart(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
	if hooks.OnCompileMorpheEnumStart == nil {
		return config, enum, nil
	}

	updatedConfig, updatedEnum, startErr := hooks.OnCompileMorpheEnumStart(config, enum)
	if startErr != nil {
		return cfg.MorpheEnumsConfig{}, yaml.Enum{}, startErr
	}

	return updatedConfig, updatedEnum, nil
}

func triggerCompileMorpheEnumSuccess(hooks hook.CompileMorpheEnum, enumType *tsdef.Enum) (*tsdef.Enum, error) {
	if hooks.OnCompileMorpheEnumSuccess == nil {
		return enumType, nil
	}
	if enumType == nil {
		return nil, ErrNoEnumType
	}
	enumTypeClone := enumType.DeepClone()

	enumType, successErr := hooks.OnCompileMorpheEnumSuccess(&enumTypeClone)
	if successErr != nil {
		return nil, successErr
	}
	return enumType, nil
}

func triggerCompileMorpheEnumFailure(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum, failureErr error) error {
	if hooks.OnCompileMorpheEnumFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEnumFailure(config, enum.DeepClone(), failureErr)
}

func getTypescriptEnum(enum yaml.Enum) (*tsdef.Enum, error) {
	enumType := tsdef.Enum{
		Name: enum.Name,
	}
	tsEnumType, tsEnumTypeErr := morpheEnumTypeToTsEnumType(enum.Type)
	if tsEnumTypeErr != nil {
		return nil, tsEnumTypeErr
	}
	enumType.Type = tsEnumType

	tsEntries, entriesErr := getTsEntriesForMorpheEnum(enum.Entries)
	if entriesErr != nil {
		return nil, entriesErr
	}
	enumType.Entries = tsEntries
	return &enumType, nil
}

func morpheEnumTypeToTsEnumType(morpheType yaml.EnumType) (tsdef.TsType, error) {
	switch morpheType {
	case yaml.EnumTypeInteger:
		return tsdef.TsTypeNumber, nil
	case yaml.EnumTypeFloat:
		return tsdef.TsTypeNumber, nil
	case yaml.EnumTypeString:
		return tsdef.TsTypeString, nil
	default:
		return nil, ErrUnsupportedEnumType(morpheType)
	}
}

func getTsEntriesForMorpheEnum(entries map[string]any) ([]tsdef.EnumEntry, error) {
	tsEntries := []tsdef.EnumEntry{}
	entryNames := core.MapKeysSorted(entries)

	for _, entryName := range entryNames {
		entryValue, entryExists := entries[entryName]
		if !entryExists {
			return nil, ErrEnumEntryNotFound(entryName)
		}
		tsEntries = append(tsEntries, tsdef.EnumEntry{
			Name:  entryName,
			Value: entryValue,
		})
	}
	return tsEntries, nil
}
