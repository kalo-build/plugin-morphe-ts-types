package hook

import (
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
)

type WriteTsType struct {
	OnWriteTsTypeStart   OnWriteTsTypeStartHook
	OnWriteTsTypeSuccess OnWriteTsTypeSuccessHook
	OnWriteTsTypeFailure OnWriteTsTypeFailureHook
}

type OnWriteTsTypeStartHook = func(writer write.TypescriptTypeWriter, modelType *tsdef.TsType) (write.TypescriptTypeWriter, *tsdef.TsType, error)
type OnWriteTsTypeSuccessHook = func(modelStruct *tsdef.TsType, modelStructContents []byte) (*tsdef.TsType, []byte, error)
type OnWriteTsTypeFailureHook = func(writer write.TypescriptTypeWriter, modelStruct *tsdef.TsType, failureErr error) error
