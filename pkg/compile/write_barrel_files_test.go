package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile"
	"github.com/stretchr/testify/suite"
)

type BarrelFilesTestSuite struct {
	suite.Suite
	outputDir string
}

func TestBarrelFilesTestSuite(t *testing.T) {
	suite.Run(t, new(BarrelFilesTestSuite))
}

func (suite *BarrelFilesTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "barrel-test-*")
	suite.Require().NoError(err)
	suite.outputDir = dir
}

func (suite *BarrelFilesTestSuite) TearDownTest() {
	os.RemoveAll(suite.outputDir)
}

func (suite *BarrelFilesTestSuite) TestWriteBarrelFiles_SingleCategory() {
	modelsDir := filepath.Join(suite.outputDir, "models")

	err := compile.WriteBarrelFiles([]compile.BarrelFileCategory{
		{
			DirPath:         modelsDir,
			DefinitionNames: []string{"User", "Organization"},
		},
	})

	suite.Require().NoError(err)

	content, readErr := os.ReadFile(filepath.Join(modelsDir, "index.d.ts"))
	suite.Require().NoError(readErr)

	expected := "export * from './organization'\nexport * from './user'\n"
	suite.Equal(expected, string(content))
}

func (suite *BarrelFilesTestSuite) TestWriteBarrelFiles_MultipleCategories() {
	enumsDir := filepath.Join(suite.outputDir, "enums")
	modelsDir := filepath.Join(suite.outputDir, "models")
	entitiesDir := filepath.Join(suite.outputDir, "entities")

	err := compile.WriteBarrelFiles([]compile.BarrelFileCategory{
		{DirPath: enumsDir, DefinitionNames: []string{"InvoiceStatus", "TaskPriority"}},
		{DirPath: modelsDir, DefinitionNames: []string{"Invoice"}},
		{DirPath: entitiesDir, DefinitionNames: []string{"Invoice", "Task"}},
	})

	suite.Require().NoError(err)

	enumContent, _ := os.ReadFile(filepath.Join(enumsDir, "index.d.ts"))
	suite.Contains(string(enumContent), "export * from './invoice-status'")
	suite.Contains(string(enumContent), "export * from './task-priority'")

	modelContent, _ := os.ReadFile(filepath.Join(modelsDir, "index.d.ts"))
	suite.Contains(string(modelContent), "export * from './invoice'")

	entityContent, _ := os.ReadFile(filepath.Join(entitiesDir, "index.d.ts"))
	suite.Contains(string(entityContent), "export * from './invoice'")
	suite.Contains(string(entityContent), "export * from './task'")
}

func (suite *BarrelFilesTestSuite) TestWriteBarrelFiles_EmptyCategory_Skipped() {
	modelsDir := filepath.Join(suite.outputDir, "models")

	err := compile.WriteBarrelFiles([]compile.BarrelFileCategory{
		{DirPath: modelsDir, DefinitionNames: []string{}},
	})

	suite.Require().NoError(err)

	_, statErr := os.Stat(filepath.Join(modelsDir, "index.d.ts"))
	suite.True(os.IsNotExist(statErr), "should not create barrel file for empty category")
}

func (suite *BarrelFilesTestSuite) TestWriteBarrelFiles_SortedAlphabetically() {
	modelsDir := filepath.Join(suite.outputDir, "models")

	err := compile.WriteBarrelFiles([]compile.BarrelFileCategory{
		{DirPath: modelsDir, DefinitionNames: []string{"Zebra", "Apple", "Mango"}},
	})

	suite.Require().NoError(err)

	content, _ := os.ReadFile(filepath.Join(modelsDir, "index.d.ts"))
	expected := "export * from './apple'\nexport * from './mango'\nexport * from './zebra'\n"
	suite.Equal(expected, string(content))
}

func (suite *BarrelFilesTestSuite) TestWriteBarrelFiles_KebabCaseNames() {
	modelsDir := filepath.Join(suite.outputDir, "models")

	err := compile.WriteBarrelFiles([]compile.BarrelFileCategory{
		{DirPath: modelsDir, DefinitionNames: []string{"InvoiceLineItem", "AdminSetting"}},
	})

	suite.Require().NoError(err)

	content, _ := os.ReadFile(filepath.Join(modelsDir, "index.d.ts"))
	suite.Contains(string(content), "export * from './invoice-line-item'")
	suite.Contains(string(content), "export * from './admin-setting'")
}
