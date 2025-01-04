package tsfile

import (
	"os"
	"path/filepath"

	"github.com/kaloseia/go-util/strcase"
)

func WriteAppendTsDefinitionFile(dirPath string, definitionName string, definitionFileContents string) ([]byte, error) {
	definitionFileName := strcase.ToKebabCaseLower(definitionName)
	definitionFilePath := filepath.Join(dirPath, definitionFileName+".d.ts")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return []byte(definitionFileContents), appendToFile(definitionFilePath, definitionFileContents)
}

func appendToFile(filePath string, content string) error {
	fileHandle, handleErr := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if handleErr != nil {
		return handleErr
	}
	defer fileHandle.Close()

	_, writeErr := fileHandle.WriteString(content)
	return writeErr
}
