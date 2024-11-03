package compile

import (
	"fmt"

	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsfile"
)

type MorpheEnumFileWriter struct {
	TargetDirPath string
}

func (w *MorpheEnumFileWriter) WriteEnum(enumName string, enumDefinition *tsdef.Enum) ([]byte, error) {
	allEnumLines, allLinesErr := w.getAllEnumLines(enumName, enumDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	enumFileContents, enumContentsErr := core.LinesToString(allEnumLines)
	if enumContentsErr != nil {
		return nil, enumContentsErr
	}

	return tsfile.WriteAppendTsDefinitionFile(w.TargetDirPath, enumName, enumFileContents)
}

func (w *MorpheEnumFileWriter) getAllEnumLines(enumName string, enumDefinition *tsdef.Enum) ([]string, error) {
	allEnumLines := []string{}
	if enumName != enumDefinition.Name {
		allEnumLines = append(allEnumLines, "")
	}

	allEnumLines = append(allEnumLines, fmt.Sprintf(`export enum %s {`, enumDefinition.Name))

	for _, enumEntry := range enumDefinition.Entries {
		entryName := strcase.ToPascalCase(enumEntry.Name)
		entryValue := enumEntry.Value
		enumEntryLine := fmt.Sprintf("\t%s: %v", entryName, entryValue)
		allEnumLines = append(allEnumLines, enumEntryLine)
	}

	allEnumLines = append(allEnumLines, "}")
	return allEnumLines, nil
}
