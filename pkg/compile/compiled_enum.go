package compile

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

type CompiledEnum struct {
	Enum         *tsdef.Enum
	EnumContents []byte
}
