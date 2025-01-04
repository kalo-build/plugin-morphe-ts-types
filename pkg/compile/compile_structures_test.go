package compile_test

import (
	"fmt"
	"testing"

	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/stretchr/testify/suite"
)

type CompileStructuresTestSuite struct {
	suite.Suite
}

func TestCompileStructuresTestSuite(t *testing.T) {
	suite.Run(t, new(CompileStructuresTestSuite))
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToTsObject() {
	structureHooks := hook.CompileMorpheStructure{}
	structuresConfig := cfg.MorpheStructuresConfig{}

	structure0 := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
			"HouseNr": {
				Type: yaml.StructureFieldTypeString,
			},
			"ZipCode": {
				Type: yaml.StructureFieldTypeString,
			},
			"City": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	tsObject, tsObjectErr := compile.MorpheStructureToTsObject(structureHooks, structuresConfig, r, structure0)

	suite.Nil(tsObjectErr)
	suite.NotNil(tsObject)

	suite.Equal(tsObject.Name, "Address")

	tsFields := tsObject.Fields
	suite.Len(tsFields, 4)

	tsField0 := tsFields[0]
	suite.Equal(tsField0.Name, "City")
	suite.Equal(tsField0.Type, tsdef.TsTypeString)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "HouseNr")
	suite.Equal(tsField1.Type, tsdef.TsTypeString)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "Street")
	suite.Equal(tsField2.Type, tsdef.TsTypeString)

	tsField3 := tsFields[3]
	suite.Equal(tsField3.Name, "ZipCode")
	suite.Equal(tsField3.Type, tsdef.TsTypeString)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToTsObject_StartHook_Successful() {
	var featureFlag = "otherName"
	structureHooks := hook.CompileMorpheStructure{
		OnCompileMorpheStructureStart: func(config cfg.MorpheStructuresConfig, structure yaml.Structure) (cfg.MorpheStructuresConfig, yaml.Structure, error) {
			if featureFlag != "otherName" {
				return config, structure, nil
			}
			structure.Name = structure.Name + "CHANGED"
			delete(structure.Fields, "City")
			return config, structure, nil
		},
	}

	structuresConfig := cfg.MorpheStructuresConfig{}

	structure0 := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
			"HouseNr": {
				Type: yaml.StructureFieldTypeString,
			},
			"ZipCode": {
				Type: yaml.StructureFieldTypeString,
			},
			"City": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	tsObject, tsObjectErr := compile.MorpheStructureToTsObject(structureHooks, structuresConfig, r, structure0)

	suite.Nil(tsObjectErr)
	suite.NotNil(tsObject)

	suite.Equal(tsObject.Name, "AddressCHANGED")

	tsFields := tsObject.Fields
	suite.Len(tsFields, 3)

	tsField0 := tsFields[0]
	suite.Equal(tsField0.Name, "HouseNr")
	suite.Equal(tsField0.Type, tsdef.TsTypeString)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "Street")
	suite.Equal(tsField1.Type, tsdef.TsTypeString)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "ZipCode")
	suite.Equal(tsField2.Type, tsdef.TsTypeString)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToTsObject_StartHook_Failure() {
	var featureFlag = "otherName"
	structureHooks := hook.CompileMorpheStructure{
		OnCompileMorpheStructureStart: func(config cfg.MorpheStructuresConfig, structure yaml.Structure) (cfg.MorpheStructuresConfig, yaml.Structure, error) {
			if featureFlag != "otherName" {
				return config, structure, nil
			}
			return config, structure, fmt.Errorf("compile structure start hook error")
		},
	}

	structuresConfig := cfg.MorpheStructuresConfig{}

	structure0 := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	tsObject, tsObjectErr := compile.MorpheStructureToTsObject(structureHooks, structuresConfig, r, structure0)

	suite.ErrorContains(tsObjectErr, "compile structure start hook error")
	suite.Nil(tsObject)
}
