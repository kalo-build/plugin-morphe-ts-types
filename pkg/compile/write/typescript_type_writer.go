package write

import "github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"

type TypescriptTypeWriter interface {
	WriteStruct(*tsdef.TsType) ([]byte, error)
}
