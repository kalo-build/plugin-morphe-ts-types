package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/internal/testutils"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	EnumsDirPath      string
	ModelsDirPath     string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
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
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},

		MorpheEnumsConfig: cfg.MorpheEnumsConfig{},
		EnumWriter: &compile.MorpheEnumFileWriter{
			TargetDirPath: workingDirPath + "/enums",
		},

		MorpheModelsConfig: cfg.MorpheModelsConfig{},
		ModelWriter: &compile.MorpheObjectFileWriter{
			TargetDirPath: workingDirPath + "/models",
		},

		MorpheEntitiesConfig: cfg.MorpheEntitiesConfig{},
		EntityWriter: &compile.MorpheObjectFileWriter{
			TargetDirPath: workingDirPath + "/entities",
		},

		MorpheStructuresConfig: cfg.MorpheStructuresConfig{},
		StructureWriter: &compile.MorpheObjectFileWriter{
			TargetDirPath: workingDirPath + "/structures",
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

	modelPath2 := modelsDirPath + "/company.d.ts"
	gtModelPath2 := gtModelsDirPath + "/company.d.ts"
	suite.FileExists(modelPath2)
	suite.FileEquals(modelPath2, gtModelPath2)

	enumsDirPath := workingDirPath + "/enums"
	gtEnumsDirPath := suite.TestGroundTruthDirPath + "/enums"
	suite.DirExists(enumsDirPath)

	enumPath0 := enumsDirPath + "/nationality.d.ts"
	gtEnumPath0 := gtEnumsDirPath + "/nationality.d.ts"
	suite.FileExists(enumPath0)
	suite.FileEquals(enumPath0, gtEnumPath0)

	enumPath1 := enumsDirPath + "/universal-number.d.ts"
	gtEnumPath1 := gtEnumsDirPath + "/universal-number.d.ts"
	suite.FileExists(enumPath1)
	suite.FileEquals(enumPath1, gtEnumPath1)

	structuresDirPath := workingDirPath + "/structures"
	gtStructuresDirPath := suite.TestGroundTruthDirPath + "/structures"
	suite.DirExists(structuresDirPath)

	structurePath0 := structuresDirPath + "/address.d.ts"
	gtStructurePath0 := gtStructuresDirPath + "/address.d.ts"
	suite.FileExists(structurePath0)
	suite.FileEquals(structurePath0, gtStructurePath0)

	entitiesDirPath := workingDirPath + "/entities"
	gtEntitiesDirPath := suite.TestGroundTruthDirPath + "/entities"
	suite.DirExists(entitiesDirPath)

	entityPath0 := entitiesDirPath + "/company.d.ts"
	gtEntityPath0 := gtEntitiesDirPath + "/company.d.ts"
	suite.FileExists(entityPath0)
	suite.FileEquals(entityPath0, gtEntityPath0)

	entityPath1 := entitiesDirPath + "/person.d.ts"
	gtEntityPath1 := gtEntitiesDirPath + "/person.d.ts"
	suite.FileExists(entityPath1)
	suite.FileEquals(entityPath1, gtEntityPath1)
}
