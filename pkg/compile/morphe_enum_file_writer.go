package compile

import (
	"fmt"
	ti "time"

	"github.com/kalo/go-util/core"
	"github.com/kalo/go-util/strcase"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsfile"
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

	for enumIdx, enumEntry := range enumDefinition.Entries {
		entryName := strcase.ToPascalCase(enumEntry.Name)
		entryValue := w.formatEnumValue(enumEntry.Value)
		enumEntryLine := fmt.Sprintf("\t%s = %v", entryName, entryValue)
		if enumIdx != len(enumDefinition.Entries)-1 {
			enumEntryLine += ","
		}
		allEnumLines = append(allEnumLines, enumEntryLine)
	}

	allEnumLines = append(allEnumLines, "}")
	return allEnumLines, nil
}

func (w *MorpheEnumFileWriter) formatEnumValue(value any) string {
	switch typedValue := value.(type) {
	case string:
		return fmt.Sprintf("'%v'", typedValue)
	case ti.Time:
		formattedValue := ""
		if typedValue.Hour() == 0 && typedValue.Minute() == 0 && typedValue.Second() == 0 && typedValue.Nanosecond() == 0 {
			formattedValue = typedValue.Format("2006-01-02")
		} else {
			formattedValue = typedValue.Format(ti.RFC3339)
		}
		return fmt.Sprintf("'%v'", formattedValue)
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
