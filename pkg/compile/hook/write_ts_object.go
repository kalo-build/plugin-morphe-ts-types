package hook

import (
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
)

type WriteTsObject struct {
	OnWriteTsObjectStart   OnWriteTsObjectStartHook
	OnWriteTsObjectSuccess OnWriteTsObjectSuccessHook
	OnWriteTsObjectFailure OnWriteTsObjectFailureHook
}

type OnWriteTsObjectStartHook = func(writer write.TsObjectWriter, modelType *tsdef.Object) (write.TsObjectWriter, *tsdef.Object, error)
type OnWriteTsObjectSuccessHook = func(modelStruct *tsdef.Object, modelStructContents []byte) (*tsdef.Object, []byte, error)
type OnWriteTsObjectFailureHook = func(writer write.TsObjectWriter, modelStruct *tsdef.Object, failureErr error) error
