package compile_test

import (
	"fmt"
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-ts-types/pkg/tsdef"
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
	suite.Equal(tsField0.Name, "city")
	suite.Equal(tsField0.Type, tsdef.TsTypeString)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "houseNr")
	suite.Equal(tsField1.Type, tsdef.TsTypeString)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "street")
	suite.Equal(tsField2.Type, tsdef.TsTypeString)

	tsField3 := tsFields[3]
	suite.Equal(tsField3.Name, "zipCode")
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
	suite.Equal(tsField0.Name, "houseNr")
	suite.Equal(tsField0.Type, tsdef.TsTypeString)

	tsField1 := tsFields[1]
	suite.Equal(tsField1.Name, "street")
	suite.Equal(tsField1.Type, tsdef.TsTypeString)

	tsField2 := tsFields[2]
	suite.Equal(tsField2.Name, "zipCode")
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

func (suite *CompileStructuresTestSuite) TestMorpheStructureToTsObject_StructureComposition() {
	structureHooks := hook.CompileMorpheStructure{}
	structuresConfig := cfg.MorpheStructuresConfig{}

	lineItemStructure := yaml.Structure{
		Name: "InvoiceLineItem",
		Fields: map[string]yaml.StructureField{
			"Description": {Type: yaml.StructureFieldTypeString},
			"Quantity":   {Type: yaml.StructureFieldTypeInteger},
			"UnitAmount": {Type: yaml.StructureFieldTypeInteger},
		},
	}

	invoiceStructure := yaml.Structure{
		Name: "Invoice",
		Fields: map[string]yaml.StructureField{
			"ID": {Type: yaml.StructureFieldTypeString},
			"LineItem": {
				Type:       yaml.StructureFieldType("InvoiceLineItem"),
				Attributes: []string{"optional"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetStructure("InvoiceLineItem", lineItemStructure)
	r.SetStructure("Invoice", invoiceStructure)

	tsObject, tsObjectErr := compile.MorpheStructureToTsObject(structureHooks, structuresConfig, r, invoiceStructure)

	suite.Nil(tsObjectErr)
	suite.NotNil(tsObject)
	suite.Equal(tsObject.Name, "Invoice")

	suite.Len(tsObject.Fields, 2)

	idField := tsObject.Fields[0]
	suite.Equal(idField.Name, "id")
	suite.Equal(idField.Type, tsdef.TsTypeString)

	lineItemField := tsObject.Fields[1]
	suite.Equal(lineItemField.Name, "lineItem")
	suite.True(lineItemField.Type.IsOptional())
	refType, ok := lineItemField.Type.(tsdef.TsTypeOptional)
	suite.True(ok)
	objType, ok := refType.ValueType.(tsdef.TsTypeObject)
	suite.True(ok)
	suite.Equal(objType.Name, "InvoiceLineItem")
	suite.Equal(objType.ModulePath, "./invoice-line-item")
}
