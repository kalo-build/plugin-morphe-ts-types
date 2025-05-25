package write

import "github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"

type TsObjectWriter interface {
	WriteObject(mainObjectName string, objectDefinition *tsdef.Object) ([]byte, error)
	ClearFile(mainObjectName string) error
}
