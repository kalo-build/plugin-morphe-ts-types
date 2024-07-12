package write

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

type TsObjectWriter interface {
	WriteObject(mainObjectName string, objectDefinition *tsdef.Object) ([]byte, error)
}
