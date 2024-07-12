package compile

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

// CompiledModelObjects maps Model.Name -> ModelObject.Name -> CompiledObject
type CompiledModelObjects map[string]map[string]CompiledObject

func (objects CompiledModelObjects) AddCompiledModelObject(modelName string, objectDef *tsdef.Object, objectContents []byte) {
	if objects[modelName] == nil {
		objects[modelName] = make(map[string]CompiledObject)
	}
	objects[modelName][objectDef.Name] = CompiledObject{
		Object:         objectDef,
		ObjectContents: objectContents,
	}
}

func (objects CompiledModelObjects) GetAllCompiledModelObjects(modelName string) map[string]CompiledObject {
	modelObjects, modelObjectsExist := objects[modelName]
	if !modelObjectsExist {
		return nil
	}
	return modelObjects
}

func (objects CompiledModelObjects) GetCompiledModelObject(modelName string, objectName string) CompiledObject {
	modelObjects, modelObjectsExist := objects[modelName]
	if !modelObjectsExist {
		return CompiledObject{}
	}
	compiledObject, compiledObjectExists := modelObjects[objectName]
	if !compiledObjectExists {
		return CompiledObject{}
	}
	return compiledObject
}
