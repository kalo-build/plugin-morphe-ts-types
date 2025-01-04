package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kaloseia/go-util/assertfile"
	rcfg "github.com/kaloseia/morphe-go/pkg/registry/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/internal/testutils"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	EnumsDirPath    string
	ModelsDirPath   string
	EntitiesDirPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToTypescript() {
	workingDirPath := suite.TestDirPath + "/working"
	suite.Nil(os.Mkdir(workingDirPath, 0644))
	defer os.RemoveAll(workingDirPath)

	config := compile.MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:    suite.EnumsDirPath,
			RegistryModelsDirPath:   suite.ModelsDirPath,
			RegistryEntitiesDirPath: suite.EntitiesDirPath,
		},

		MorpheEnumsConfig: cfg.MorpheEnumsConfig{},
		EnumWriter: &compile.MorpheEnumFileWriter{
			TargetDirPath: workingDirPath + "/enums",
		},

		MorpheModelsConfig: cfg.MorpheModelsConfig{},
		ModelWriter: &compile.MorpheObjectFileWriter{
			TargetDirPath: workingDirPath + "/models",
		},
	}

	compileErr := compile.MorpheToTypescript(config)

	suite.NoError(compileErr)

	modelsDirPath := workingDirPath + "/models"
	gtModelsDirPath := suite.TestGroundTruthDirPath + "/models"
	suite.DirExists(modelsDirPath)

	modelPath0 := modelsDirPath + "/contact-info.d.ts"
	gtModelPath0 := gtModelsDirPath + "/contact-info.d.ts"
	suite.FileExists(modelPath0)
	suite.FileEquals(modelPath0, gtModelPath0)

	modelPath1 := modelsDirPath + "/person.d.ts"
	gtModelPath1 := gtModelsDirPath + "/person.d.ts"
	suite.FileExists(modelPath1)
	suite.FileEquals(modelPath1, gtModelPath1)

	enumsDirPath := workingDirPath + "/enums"
	gtEnumsDirPath := suite.TestGroundTruthDirPath + "/enums"
	suite.DirExists(enumsDirPath)

	enumPath0 := enumsDirPath + "/nationality.d.ts"
	gtEnumPath0 := gtEnumsDirPath + "/nationality.d.ts"
	suite.FileExists(enumPath0)
	suite.FileEquals(enumPath0, gtEnumPath0)
}
