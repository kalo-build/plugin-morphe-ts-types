package tsfile

import (
	"os"
	"path/filepath"

	"github.com/kaloseia/go-util/strcase"
)

func WriteAppendTsObjectFile(dirPath string, objectName string, objectFileContents string) ([]byte, error) {
	objectFileName := strcase.ToKebabCaseLower(objectName)
	objectFilePath := filepath.Join(dirPath, objectFileName+".d.ts")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return []byte(objectFileContents), appendToFile(objectFilePath, objectFileContents)
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
