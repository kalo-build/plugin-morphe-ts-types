package write

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

type TsEnumWriter interface {
	WriteEnum(enumName string, enumDefinition *tsdef.Enum) ([]byte, error)
}
