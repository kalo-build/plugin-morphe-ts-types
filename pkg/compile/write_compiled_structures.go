package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

type CompiledStructureObjects map[string]*CompiledStructureObject

type CompiledStructureObject struct {
	Definition *tsdef.Object
	Contents   []byte
}

func WriteAllStructureObjectDefinitions(config MorpheCompileConfig, allStructureObjectDefs map[string]*tsdef.Object) (CompiledStructureObjects, error) {
	allWrittenStructures := CompiledStructureObjects{}

	objectKeys := core.MapKeysSorted(allStructureObjectDefs)
	for _, objectKey := range objectKeys {
		structureObject := allStructureObjectDefs[objectKey]
		if clearErr := config.StructureWriter.ClearFile(structureObject.Name); clearErr != nil {
			return nil, clearErr
		}
	}

	for _, objectKey := range objectKeys {
		structureObject := allStructureObjectDefs[objectKey]
		structureObject, structureObjectContents, writeErr := WriteStructureObjectDefinition(config.WriteObjectHooks, config.StructureWriter, structureObject)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenStructures[structureObject.Name] = &CompiledStructureObject{
			Definition: structureObject,
			Contents:   structureObjectContents,
		}
	}
	return allWrittenStructures, nil
}

func WriteStructureObjectDefinition(hooks hook.WriteTsObject, writer write.TsObjectWriter, structureObject *tsdef.Object) (*tsdef.Object, []byte, error) {
	writer, structureObject, writeStartErr := triggerWriteStructureObjectStart(hooks, writer, structureObject)
	if writeStartErr != nil {
		return nil, nil, triggerWriteStructureObjectFailure(hooks, writer, structureObject, writeStartErr)
	}

	structureObjectContents, writeStructErr := writer.WriteObject(structureObject.Name, structureObject)
	if writeStructErr != nil {
		return nil, nil, triggerWriteStructureObjectFailure(hooks, writer, structureObject, writeStructErr)
	}

	structureObject, structureObjectContents, writeSuccessErr := triggerWriteStructureObjectSuccess(hooks, structureObject, structureObjectContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteStructureObjectFailure(hooks, writer, structureObject, writeSuccessErr)
	}
	return structureObject, structureObjectContents, nil
}

func triggerWriteStructureObjectStart(hooks hook.WriteTsObject, writer write.TsObjectWriter, structureObject *tsdef.Object) (write.TsObjectWriter, *tsdef.Object, error) {
	if hooks.OnWriteTsObjectStart == nil {
		return writer, structureObject, nil
	}
	if structureObject == nil {
		return nil, nil, ErrNoStructureObject
	}
	structureObjectClone := structureObject.DeepClone()

	updatedWriter, updatedStructureObject, startErr := hooks.OnWriteTsObjectStart(writer, &structureObjectClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedStructureObject, nil
}

func triggerWriteStructureObjectSuccess(hooks hook.WriteTsObject, structureObject *tsdef.Object, structureObjectContents []byte) (*tsdef.Object, []byte, error) {
	if hooks.OnWriteTsObjectSuccess == nil {
		return structureObject, structureObjectContents, nil
	}
	if structureObject == nil {
		return nil, nil, ErrNoStructureObject
	}
	structureObjectClone := structureObject.DeepClone()
	structureObjectContentsClone := clone.Slice(structureObjectContents)

	updatedStructureObject, updatedStructureObjectContents, successErr := hooks.OnWriteTsObjectSuccess(&structureObjectClone, structureObjectContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedStructureObject, updatedStructureObjectContents, nil
}

func triggerWriteStructureObjectFailure(hooks hook.WriteTsObject, writer write.TsObjectWriter, structureObject *tsdef.Object, failureErr error) error {
	if hooks.OnWriteTsObjectFailure == nil {
		return failureErr
	}

	structureObjectClone := structureObject.DeepClone()
	return hooks.OnWriteTsObjectFailure(writer, &structureObjectClone, failureErr)
}
