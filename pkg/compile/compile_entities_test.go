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
			"Name": {
				Type: "User.Name",
			},
			"Street": {
				Type: "User.Address.Street",
			},
			"HouseNr": {
				Type: "User.Address.HouseNr",
			},
			"ZipCode": {
				Type: "User.Address.ZipCode",
			},
			"City": {
				Type: "User.Address.City",
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
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Address": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("User", userModel)

	addressModel := yaml.Model{
		Name: "Address",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Street": {
				Type: yaml.ModelFieldTypeString,
			},
			"HouseNr": {
				Type: yaml.ModelFieldTypeString,
			},
			"ZipCode": {
				Type: yaml.ModelFieldTypeString,
			},
			"City": {
				Type: yaml.ModelFieldTypeString,
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
	r.SetModel("Address", addressModel)

	tsObject, tsObjectErr := compile.MorpheEntityToTsObject(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
	suite.NotNil(tsObject)

	suite.Equal(tsObject.Name, "User")

	tsFields := tsObject.Fields
	suite.Len(tsFields, 6)

	tsField0 := tsFields[0]
	suite.Equal(tsField0.Name, "City")
	suite.Equal(tsField0.Type, tsdef.TsTypeString)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "HouseNr")
	suite.Equal(tsField1.Type, tsdef.TsTypeString)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "Name")
	suite.Equal(tsField2.Type, tsdef.TsTypeString)

	tsField3 := tsFields[3]
	suite.Equal(tsField3.Name, "Street")
	suite.Equal(tsField3.Type, tsdef.TsTypeString)

	tsField4 := tsFields[4]
	suite.Equal(tsField4.Name, "UUID")
	suite.Equal(tsField4.Type, tsdef.TsTypeString)

	tsField5 := tsFields[5]
	suite.Equal(tsField5.Name, "ZipCode")
	suite.Equal(tsField5.Type, tsdef.TsTypeString)
}
