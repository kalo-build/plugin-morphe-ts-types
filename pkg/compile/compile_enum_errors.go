package compile

import (
	"errors"
	"fmt"

	"github.com/kaloseia/morphe-go/pkg/yaml"
)

var ErrNoEnumType = errors.New("no enum type provided")

func ErrUnsupportedEnumType(morpheType yaml.EnumType) error {
	return fmt.Errorf("unsupported morphe enum type '%s' provided", morpheType)
}

func ErrEnumEntryNotFound(entryName string) error {
	return fmt.Errorf("morphe enum entry '%s' not found", entryName)
}
