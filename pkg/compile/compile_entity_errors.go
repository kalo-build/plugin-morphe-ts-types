package compile

import (
	"errors"
	"fmt"
)

var ErrNoEntityObjects = errors.New("no entity objects provided")
var ErrNoEntityObject = errors.New("no entity object provided")

var ErrInvalidEntityFieldPath = func(fieldType string) error {
	return fmt.Errorf("invalid entity field type path: %s", fieldType)
}

var ErrRootModelNotFound = func(modelName string) error {
	return fmt.Errorf("root model not found: %s", modelName)
}

var ErrRelatedModelNotFound = func(relatedName, fieldType string) error {
	return fmt.Errorf("related model not found: %s in path %s", relatedName, fieldType)
}

var ErrFailedToGetRelatedModel = func(relatedName, fieldType string) error {
	return fmt.Errorf("failed to get related model: %s in path %s", relatedName, fieldType)
}

var ErrTerminalFieldNotFound = func(fieldName, fieldType string) error {
	return fmt.Errorf("terminal field not found: %s in path %s", fieldName, fieldType)
}
