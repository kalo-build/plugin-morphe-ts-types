package compile

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

// CompiledEntityObjects maps Entity.Name -> EntityObject.Name -> CompiledObject
type CompiledEntityObjects map[string]map[string]CompiledObject

func (objects CompiledEntityObjects) AddCompiledEntityObject(entityName string, objectDef *tsdef.Object, objectContents []byte) {
	if objects[entityName] == nil {
		objects[entityName] = make(map[string]CompiledObject)
	}
	objects[entityName][objectDef.Name] = CompiledObject{
		Object:         objectDef,
		ObjectContents: objectContents,
	}
}

func (objects CompiledEntityObjects) GetAllCompiledEntityObjects(entityName string) map[string]CompiledObject {
	entityObjects, entityObjectsExist := objects[entityName]
	if !entityObjectsExist {
		return nil
	}
	return entityObjects
}

func (objects CompiledEntityObjects) GetCompiledEntityObject(entityName string, objectName string) CompiledObject {
	entityObjects, entityObjectsExist := objects[entityName]
	if !entityObjectsExist {
		return CompiledObject{}
	}
	compiledObject, compiledObjectExists := entityObjects[objectName]
	if !compiledObjectExists {
		return CompiledObject{}
	}
	return compiledObject
}
