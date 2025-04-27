package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsfile"
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

	return tsfile.WriteAppendTsDefinitionFile(w.TargetDirPath, mainObjectName, objectFileContents)
}

func (w *MorpheObjectFileWriter) getAllObjectLines(mainObjectName string, objectDefinition *tsdef.Object) ([]string, error) {
	allObjectLines := []string{}

	importLines, importsErr := w.getAllObjectImportLines(objectDefinition)
	if importsErr != nil {
		return nil, importsErr
	}

	if len(importLines) > 0 {
		allObjectLines = append(allObjectLines, importLines...)
		allObjectLines = append(allObjectLines, "")
	}

	if mainObjectName != objectDefinition.Name {
		allObjectLines = append(allObjectLines, "")
	}

	allObjectLines = append(allObjectLines, fmt.Sprintf(`export type %s = {`, objectDefinition.Name))

	for _, objectField := range objectDefinition.Fields {
		fieldName := strcase.ToCamelCase(objectField.Name)
		fieldTypeSyntax := objectField.Type.GetSyntax()
		if objectField.Type.IsOptional() {
			structFieldLine := fmt.Sprintf("\t%s?: %s", fieldName, fieldTypeSyntax)
			allObjectLines = append(allObjectLines, structFieldLine)
			continue
		}
		structFieldLine := fmt.Sprintf("\t%s: %s", fieldName, fieldTypeSyntax)
		allObjectLines = append(allObjectLines, structFieldLine)
	}

	allObjectLines = append(allObjectLines, "}")
	return allObjectLines, nil
}

func (w *MorpheObjectFileWriter) getAllObjectImportLines(objectDefinition *tsdef.Object) ([]string, error) {
	if len(objectDefinition.Imports) == 0 {
		return nil, nil
	}

	filteredImportsMap := map[string]tsdef.ObjectImport{}
	for _, objectImport := range objectDefinition.Imports {
		filteredImportsMap[objectImport.ModulePath] = objectImport
	}

	allImportLines := []string{}

	filteredImports := core.MapKeysSorted(filteredImportsMap)
	for _, objectImportPath := range filteredImports {
		objectImport := filteredImportsMap[objectImportPath]
		if len(objectImport.ModuleNames) <= 3 {
			importNames := strings.Join(objectImport.ModuleNames, ", ")
			allImportLines = append(allImportLines, `import { `+importNames+` } from "`+objectImportPath+`"`)
			continue
		}
		allImportLines = append(allImportLines, `import { `)
		for _, importName := range objectImport.ModuleNames {
			allImportLines = append(allImportLines, importName+`,`)
		}
		allImportLines = append(allImportLines, `} from "`+objectImportPath+`"`)
	}

	return allImportLines, nil
}
