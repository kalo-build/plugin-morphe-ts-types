package compile

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

type CompiledObject struct {
	Object         *tsdef.Object
	ObjectContents []byte
}
