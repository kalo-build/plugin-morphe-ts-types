package compile

import (
	"errors"
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

var ErrNoRegistry = errors.New("registry not initialized")
var ErrNoMorpheModelName = errors.New("morphe model has no name")
var ErrNoMorpheModelFields = errors.New("morphe model has no fields")
var ErrNoMorpheModelIdentifiers = errors.New("morphe model has no identifiers")

func ErrUnsupportedMorpheFieldType[TType yaml.ModelFieldType | yaml.StructureFieldType | yaml.ModelFieldPath](unsupportedType TType) error {
	return fmt.Errorf("unsupported morphe field type for typescript conversion: '%s'", unsupportedType)
}

func ErrMissingMorpheIdentifierField(modelName string, identifierName string, fieldName string) error {
	return fmt.Errorf("morphe model '%s' has no field '%s' referenced in identifiers ('%s')", modelName, identifierName, fieldName)
}
