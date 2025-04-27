package compile

import "github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"

type CompiledEnum struct {
	Enum         *tsdef.Enum
	EnumContents []byte
}
