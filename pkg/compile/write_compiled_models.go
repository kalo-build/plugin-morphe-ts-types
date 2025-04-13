package compile

import (
	"github.com/kalo/clone"
	"github.com/kalo/go-util/core"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
)

func WriteAllModelObjectDefinitions(config MorpheCompileConfig, allModelObjectDefs map[string][]*tsdef.Object) (CompiledModelObjects, error) {
	allWrittenModels := CompiledModelObjects{}

	sortedModelNames := core.MapKeysSorted(allModelObjectDefs)
	for _, modelName := range sortedModelNames {
		modelObjects := allModelObjectDefs[modelName]
		for _, subModelObject := range modelObjects {
			subModelObject, subModelObjectContents, writeErr := WriteModelObjectDefinition(config.WriteObjectHooks, config.ModelWriter, modelName, subModelObject)
			if writeErr != nil {
				return nil, writeErr
			}
			allWrittenModels.AddCompiledModelObject(modelName, subModelObject, subModelObjectContents)
		}
	}
	return allWrittenModels, nil
}

func WriteModelObjectDefinition(hooks hook.WriteTsObject, writer write.TsObjectWriter, mainObjectName string, modelObject *tsdef.Object) (*tsdef.Object, []byte, error) {
	writer, modelObject, writeStartErr := triggerWriteModelObjectStart(hooks, writer, modelObject)
	if writeStartErr != nil {
		return nil, nil, triggerWriteModelObjectFailure(hooks, writer, modelObject, writeStartErr)
	}

	modelObjectContents, writeStructErr := writer.WriteObject(mainObjectName, modelObject)
	if writeStructErr != nil {
		return nil, nil, triggerWriteModelObjectFailure(hooks, writer, modelObject, writeStructErr)
	}

	modelObject, modelObjectContents, writeSuccessErr := triggerWriteModelObjectSuccess(hooks, modelObject, modelObjectContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteModelObjectFailure(hooks, writer, modelObject, writeSuccessErr)
	}
	return modelObject, modelObjectContents, nil
}

func triggerWriteModelObjectStart(hooks hook.WriteTsObject, writer write.TsObjectWriter, modelObject *tsdef.Object) (write.TsObjectWriter, *tsdef.Object, error) {
	if hooks.OnWriteTsObjectStart == nil {
		return writer, modelObject, nil
	}
	if modelObject == nil {
		return nil, nil, ErrNoModelObject
	}
	modelObjectClone := modelObject.DeepClone()

	updatedWriter, updatedModelObject, startErr := hooks.OnWriteTsObjectStart(writer, &modelObjectClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedModelObject, nil
}

func triggerWriteModelObjectSuccess(hooks hook.WriteTsObject, modelObject *tsdef.Object, modelObjectContents []byte) (*tsdef.Object, []byte, error) {
	if hooks.OnWriteTsObjectSuccess == nil {
		return modelObject, modelObjectContents, nil
	}
	if modelObject == nil {
		return nil, nil, ErrNoModelObject
	}
	modelObjectClone := modelObject.DeepClone()
	modelObjectContentsClone := clone.Slice(modelObjectContents)

	updatedModelObject, updatedModelObjectContents, successErr := hooks.OnWriteTsObjectSuccess(&modelObjectClone, modelObjectContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedModelObject, updatedModelObjectContents, nil
}

func triggerWriteModelObjectFailure(hooks hook.WriteTsObject, writer write.TsObjectWriter, modelObject *tsdef.Object, failureErr error) error {
	if hooks.OnWriteTsObjectFailure == nil {
		return failureErr
	}

	modelObjectClone := modelObject.DeepClone()
	return hooks.OnWriteTsObjectFailure(writer, &modelObjectClone, failureErr)
}
