package write

import "github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"

type TsEnumWriter interface {
	WriteEnum(enumName string, enumDefinition *tsdef.Enum) ([]byte, error)
	ClearFile(enumName string) error
}
