package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsfile"
)

type MorpheObjectFileWriter struct {
	TargetDirPath string
}

func (w *MorpheObjectFileWriter) WriteObject(mainObjectName string, objectDefinition *tsdef.Object) ([]byte, error) {
	allObjectLines, allLinesErr := w.getAllObjectLines(mainObjectName, objectDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	objectFileContents, objectContentsErr := core.LinesToString(allObjectLines)
	if objectContentsErr != nil {
		return nil, objectContentsErr
	}

	return tsfile.WriteAppendTsObjectFile(w.TargetDirPath, mainObjectName, objectFileContents)
}

func (w *MorpheObjectFileWriter) getAllObjectLines(mainObjectName string, objectDefinition *tsdef.Object) ([]string, error) {
	allObjectLines := []string{}
	if mainObjectName != objectDefinition.Name {
		allObjectLines = append(allObjectLines, "")
	}

	allObjectLines = append(allObjectLines, fmt.Sprintf(`export type %s = {`, objectDefinition.Name))

	for _, objectField := range objectDefinition.Fields {
		fieldName := strcase.ToCamelCase(objectField.Name)
		fieldTypeSyntax := objectField.Type.GetSyntax()
		structFieldLine := fmt.Sprintf("\t%s: %s", fieldName, fieldTypeSyntax)
		allObjectLines = append(allObjectLines, structFieldLine)
	}

	allObjectLines = append(allObjectLines, "}")
	return allObjectLines, nil
}
