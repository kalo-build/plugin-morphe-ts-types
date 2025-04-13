package compile

import (
	"github.com/kalo/clone"
	"github.com/kalo/go-util/core"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
)

func WriteAllEntityObjectDefinitions(config MorpheCompileConfig, allEntityObjectDefs map[string][]*tsdef.Object) (CompiledEntityObjects, error) {
	allWrittenEntities := CompiledEntityObjects{}

	sortedEntityNames := core.MapKeysSorted(allEntityObjectDefs)
	for _, entityName := range sortedEntityNames {
		entityObjects := allEntityObjectDefs[entityName]
		for _, subEntityObject := range entityObjects {
			subEntityObject, subEntityObjectContents, writeErr := WriteEntityObjectDefinition(config.WriteObjectHooks, config.EntityWriter, entityName, subEntityObject)
			if writeErr != nil {
				return nil, writeErr
			}
			allWrittenEntities.AddCompiledEntityObject(entityName, subEntityObject, subEntityObjectContents)
		}
	}
	return allWrittenEntities, nil
}

func WriteEntityObjectDefinition(hooks hook.WriteTsObject, writer write.TsObjectWriter, mainObjectName string, entityObject *tsdef.Object) (*tsdef.Object, []byte, error) {
	writer, entityObject, writeStartErr := triggerWriteEntityObjectStart(hooks, writer, entityObject)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEntityObjectFailure(hooks, writer, entityObject, writeStartErr)
	}

	entityObjectContents, writeStructErr := writer.WriteObject(mainObjectName, entityObject)
	if writeStructErr != nil {
		return nil, nil, triggerWriteEntityObjectFailure(hooks, writer, entityObject, writeStructErr)
	}

	entityObject, entityObjectContents, writeSuccessErr := triggerWriteEntityObjectSuccess(hooks, entityObject, entityObjectContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEntityObjectFailure(hooks, writer, entityObject, writeSuccessErr)
	}
	return entityObject, entityObjectContents, nil
}

func triggerWriteEntityObjectStart(hooks hook.WriteTsObject, writer write.TsObjectWriter, entityObject *tsdef.Object) (write.TsObjectWriter, *tsdef.Object, error) {
	if hooks.OnWriteTsObjectStart == nil {
		return writer, entityObject, nil
	}
	if entityObject == nil {
		return nil, nil, ErrNoEntityObject
	}
	entityObjectClone := entityObject.DeepClone()

	updatedWriter, updatedEntityObject, startErr := hooks.OnWriteTsObjectStart(writer, &entityObjectClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEntityObject, nil
}

func triggerWriteEntityObjectSuccess(hooks hook.WriteTsObject, entityObject *tsdef.Object, entityObjectContents []byte) (*tsdef.Object, []byte, error) {
	if hooks.OnWriteTsObjectSuccess == nil {
		return entityObject, entityObjectContents, nil
	}
	if entityObject == nil {
		return nil, nil, ErrNoEntityObject
	}
	entityObjectClone := entityObject.DeepClone()
	entityObjectContentsClone := clone.Slice(entityObjectContents)

	updatedEntityObject, updatedEntityObjectContents, successErr := hooks.OnWriteTsObjectSuccess(&entityObjectClone, entityObjectContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEntityObject, updatedEntityObjectContents, nil
}

func triggerWriteEntityObjectFailure(hooks hook.WriteTsObject, writer write.TsObjectWriter, entityObject *tsdef.Object, failureErr error) error {
	if hooks.OnWriteTsObjectFailure == nil {
		return failureErr
	}

	entityObjectClone := entityObject.DeepClone()
	return hooks.OnWriteTsObjectFailure(writer, &entityObjectClone, failureErr)
}
