package compile

import (
	"errors"
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

var ErrNoEnumType = errors.New("no enum type provided")
var ErrNoEnum = errors.New("no enum provided")

func ErrUnsupportedEnumType(morpheType yaml.EnumType) error {
	return fmt.Errorf("unsupported morphe enum type '%s' provided", morpheType)
}

func ErrEnumEntryNotFound(entryName string) error {
	return fmt.Errorf("morphe enum entry '%s' not found", entryName)
}
