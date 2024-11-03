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

type CompileEnumsTestSuite struct {
	suite.Suite
}

func TestCompileEnumsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEnumsTestSuite))
}

func (suite *CompileEnumsTestSuite) SetupTest() {
}

func (suite *CompileEnumsTestSuite) TearDownTest() {
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_String() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(tsEnumErr)

	suite.Equal(tsEnum.Name, "Color")
	suite.Equal(tsEnum.Type, tsdef.TsTypeString)

	tsEntries0 := tsEnum.Entries
	suite.Len(tsEntries0, 3)

	tsEntry00 := tsEntries0[0]
	suite.Equal(tsEntry00.Name, "Blue")
	suite.Equal(tsEntry00.Value, "rgb(0,0,255)")

	tsEntry01 := tsEntries0[1]
	suite.Equal(tsEntry01.Name, "Green")
	suite.Equal(tsEntry01.Value, "rgb(0,255,0)")

	tsEntry02 := tsEntries0[2]
	suite.Equal(tsEntry02.Name, "Red")
	suite.Equal(tsEntry02.Value, "rgb(255,0,0)")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_Float() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeFloat,
		Entries: map[string]any{
			"Pi":    3.141,
			"Euler": 2.718,
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(tsEnumErr)

	suite.Equal(tsEnum.Name, "Analytics")
	suite.Equal(tsEnum.Type, tsdef.TsTypeNumber)

	tsEntries0 := tsEnum.Entries
	suite.Len(tsEntries0, 2)

	tsEntry00 := tsEntries0[0]
	suite.Equal(tsEntry00.Name, "Euler")
	suite.Equal(tsEntry00.Value, 2.718)

	tsEntry01 := tsEntries0[1]
	suite.Equal(tsEntry01.Name, "Pi")
	suite.Equal(tsEntry01.Value, 3.141)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_Integer() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"AnswerToLife":  42,
			"FineStructure": 317,
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(tsEnumErr)

	suite.Equal(tsEnum.Name, "Analytics")
	suite.Equal(tsEnum.Type, tsdef.TsTypeNumber)

	tsEntries0 := tsEnum.Entries
	suite.Len(tsEntries0, 2)

	tsEntry00 := tsEntries0[0]
	suite.Equal(tsEntry00.Name, "AnswerToLife")
	suite.Equal(tsEntry00.Value, 42)

	tsEntry01 := tsEntries0[1]
	suite.Equal(tsEntry01.Name, "FineStructure")
	suite.Equal(tsEntry01.Value, 317)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_NoName() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(tsEnumErr, yaml.ErrNoMorpheEnumName)

	suite.Nil(tsEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_NoType() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: "",
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(tsEnumErr, yaml.ErrNoMorpheEnumType)

	suite.Nil(tsEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_NoEntries() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name:    "Color",
		Type:    yaml.EnumTypeString,
		Entries: map[string]any{},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(tsEnumErr, yaml.ErrNoMorpheEnumEntries)

	suite.Nil(tsEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToTsObjects_EntryTypeMismatch() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	tsEnum, tsEnumErr := compile.MorpheEnumToTsEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorContains(tsEnumErr, "enum entry 'Blue' value 'rgb(0,0,255)' with type 'string' does not match the enum type of 'Integer'")

	suite.Nil(tsEnum)
}
