package compile_test

import (
	"testing"

	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/stretchr/testify/suite"
)

type CompileModelsTestSuite struct {
	suite.Suite
}

func TestCompileModelsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileModelsTestSuite))
}

func (suite *CompileModelsTestSuite) SetupTest() {
}

func (suite *CompileModelsTestSuite) TearDownTest() {
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Boolean": {
				Type: yaml.ModelFieldTypeBoolean,
			},
			"Date": {
				Type: yaml.ModelFieldTypeDate,
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
			},
			"Integer": {
				Type: yaml.ModelFieldTypeInteger,
			},
			"Protected": {
				Type: yaml.ModelFieldTypeProtected,
			},
			"Sealed": {
				Type: yaml.ModelFieldTypeSealed,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	allTsObjectTypes, allTsObjectTypesErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, model0)

	suite.Nil(allTsObjectTypesErr)
	suite.Len(allTsObjectTypes, 2)

	tsObject0 := allTsObjectTypes[0]
	suite.Equal(tsObject0.Name, "Basic")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 10)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "AutoIncrement")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "Boolean")
	suite.Equal(tsField01.Type, tsdef.TsTypeBoolean)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "Date")
	suite.Equal(tsField02.Type, tsdef.TsTypeDate)

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "Float")
	suite.Equal(tsField03.Type, tsdef.TsTypeNumber)

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "Integer")
	suite.Equal(tsField04.Type, tsdef.TsTypeNumber)

	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "Protected")
	suite.Equal(tsField05.Type, tsdef.TsTypeString)

	tsField06 := tsFields0[6]
	suite.Equal(tsField06.Name, "Sealed")
	suite.Equal(tsField06.Type, tsdef.TsTypeString)

	tsField07 := tsFields0[7]
	suite.Equal(tsField07.Name, "String")
	suite.Equal(tsField07.Type, tsdef.TsTypeString)

	tsField08 := tsFields0[8]
	suite.Equal(tsField08.Name, "Time")
	suite.Equal(tsField08.Type, tsdef.TsTypeDate)

	tsField09 := tsFields0[9]
	suite.Equal(tsField09.Name, "UUID")
	suite.Equal(tsField09.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjectTypes[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_NoFields() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name:   "Basic",
		Fields: map[string]yaml.ModelField{},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	allTsObjectTypes, allTsObjectTypesErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, model0)

	suite.NotNil(allTsObjectTypesErr)
	suite.ErrorContains(allTsObjectTypesErr, "morphe model has no fields")

	suite.Nil(allTsObjectTypes)
}

