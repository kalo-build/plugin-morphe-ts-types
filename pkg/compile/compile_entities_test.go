package compile_test

import (
	"testing"

	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-ts-types/pkg/tsdef"
	"github.com/stretchr/testify/suite"
)

type CompileEntitiesTestSuite struct {
	suite.Suite
}

func TestCompileEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEntitiesTestSuite))
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObject() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"AutoIncrement": {
				Type: "User.Child.AutoIncrement",
			},
			"Boolean": {
				Type: "User.Child.Boolean",
			},
			"Date": {
				Type: "User.Child.Date",
			},
			"Float": {
				Type: "User.Child.Float",
			},
			"Integer": {
				Type: "User.Child.Integer",
			},
			"String": {
				Type: "User.Child.String",
			},
			"Time": {
				Type: "User.Child.Time",
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Company": {
				Type: "ForOne",
			},
		},
	}

	r := registry.NewRegistry()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Child": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("User", userModel)

	childModel := yaml.Model{
		Name: "Child",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
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
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"User": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Child", childModel)

	tsObject, tsObjectErr := compile.MorpheEntityToTsObject(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
	suite.NotNil(tsObject)

	suite.Equal(tsObject.Name, "User")

	tsFields := tsObject.Fields
	suite.Len(tsFields, 8)

	tsField0 := tsFields[0]
	suite.Equal(tsField0.Name, "AutoIncrement")
	suite.Equal(tsField0.Type, tsdef.TsTypeNumber)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "Boolean")
	suite.Equal(tsField1.Type, tsdef.TsTypeBoolean)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "Date")
	suite.Equal(tsField2.Type, tsdef.TsTypeDate)

	tsField3 := tsFields[3]
	suite.Equal(tsField3.Name, "Float")
	suite.Equal(tsField3.Type, tsdef.TsTypeNumber)

	tsField4 := tsFields[4]
	suite.Equal(tsField4.Name, "Integer")
	suite.Equal(tsField4.Type, tsdef.TsTypeNumber)

	tsField5 := tsFields[5]
	suite.Equal(tsField5.Name, "String")
	suite.Equal(tsField5.Type, tsdef.TsTypeString)

	tsField6 := tsFields[6]
	suite.Equal(tsField6.Name, "Time")
	suite.Equal(tsField6.Type, tsdef.TsTypeDate)

	tsField7 := tsFields[7]
	suite.Equal(tsField7.Name, "UUID")
	suite.Equal(tsField7.Type, tsdef.TsTypeString)
}
