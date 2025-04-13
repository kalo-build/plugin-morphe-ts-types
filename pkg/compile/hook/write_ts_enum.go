package hook

import (
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/write"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
)

type WriteTsEnum struct {
	OnWriteTsEnumStart   OnWriteTsEnumStartHook
	OnWriteTsEnumSuccess OnWriteTsEnumSuccessHook
	OnWriteTsEnumFailure OnWriteTsEnumFailureHook
}

type OnWriteTsEnumStartHook = func(writer write.TsEnumWriter, enumType *tsdef.Enum) (write.TsEnumWriter, *tsdef.Enum, error)
type OnWriteTsEnumSuccessHook = func(enumType *tsdef.Enum, enumContents []byte) (*tsdef.Enum, []byte, error)
type OnWriteTsEnumFailureHook = func(writer write.TsEnumWriter, enumType *tsdef.Enum, failureErr error) error
