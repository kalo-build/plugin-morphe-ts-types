package compile

import (
	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

func WriteAllEnumDefinitions(config MorpheCompileConfig, allEnumDefs map[string]*tsdef.Enum) (CompiledEnums, error) {
	allWrittenEnums := CompiledEnums{}

	sortedEnumNames := core.MapKeysSorted(allEnumDefs)
	for _, enumName := range sortedEnumNames {
		enumDef := allEnumDefs[enumName]
		enumDef, enumContents, writeErr := WriteEnumDefinition(config.WriteEnumHooks, config.EnumWriter, enumName, enumDef)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenEnums.AddCompiledEnum(enumDef, enumContents)
	}
	return allWrittenEnums, nil
}

func WriteEnumDefinition(hooks hook.WriteTsEnum, writer write.TsEnumWriter, mainEnumName string, enum *tsdef.Enum) (*tsdef.Enum, []byte, error) {
	writer, enum, writeStartErr := triggerWriteEnumStart(hooks, writer, enum)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeStartErr)
	}

	enumContents, writeStructErr := writer.WriteEnum(mainEnumName, enum)
	if writeStructErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeStructErr)
	}

	enum, enumContents, writeSuccessErr := triggerWriteEnumSuccess(hooks, enum, enumContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeSuccessErr)
	}
	return enum, enumContents, nil
}

func triggerWriteEnumStart(hooks hook.WriteTsEnum, writer write.TsEnumWriter, enum *tsdef.Enum) (write.TsEnumWriter, *tsdef.Enum, error) {
	if hooks.OnWriteTsEnumStart == nil {
		return writer, enum, nil
	}
	if enum == nil {
		return nil, nil, ErrNoEnum
	}
	enumClone := enum.DeepClone()

	updatedWriter, updatedEnum, startErr := hooks.OnWriteTsEnumStart(writer, &enumClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEnum, nil
}

func triggerWriteEnumSuccess(hooks hook.WriteTsEnum, enum *tsdef.Enum, enumContents []byte) (*tsdef.Enum, []byte, error) {
	if hooks.OnWriteTsEnumSuccess == nil {
		return enum, enumContents, nil
	}
	if enum == nil {
		return nil, nil, ErrNoEnum
	}
	enumClone := enum.DeepClone()
	enumContentsClone := clone.Slice(enumContents)

	updatedEnum, updatedEnumContents, successErr := hooks.OnWriteTsEnumSuccess(&enumClone, enumContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEnum, updatedEnumContents, nil
}

func triggerWriteEnumFailure(hooks hook.WriteTsEnum, writer write.TsEnumWriter, enum *tsdef.Enum, failureErr error) error {
	if hooks.OnWriteTsEnumFailure == nil {
		return failureErr
	}

	enumClone := enum.DeepClone()
	return hooks.OnWriteTsEnumFailure(writer, &enumClone, failureErr)
}
