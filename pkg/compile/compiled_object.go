package compile

import "github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"

type CompiledObject struct {
	Object         *tsdef.Object
	ObjectContents []byte
}
