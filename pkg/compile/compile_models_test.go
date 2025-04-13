package compile_test

import (
	"fmt"
	"testing"

	"github.com/kalo/morphe-go/pkg/registry"
	"github.com/kalo/morphe-go/pkg/yaml"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/cfg"
	"github.com/kalo/plugin-morphe-ts-types/pkg/compile/hook"
	"github.com/kalo/plugin-morphe-ts-types/pkg/tsdef"
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
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

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_EnumField() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Nationality": {
				Type: "Nationality",
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

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}

	r := registry.NewRegistry()
	r.SetEnum("Nationality", enum0)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Basic")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 3)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "AutoIncrement")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "Nationality")
	suite.Equal(tsField01.Type, tsdef.TsTypeObject{
		ModulePath: "../enums/nationality",
		Name:       "Nationality",
	})

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "UUID")
	suite.Equal(tsField02.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_EnumField_EnumNotFound() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Nationality": {
				Type: "Nationality",
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "unsupported morphe field type for typescript conversion: 'Nationality'")

	suite.Nil(allTsObjects)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_ForOne() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Basic")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "ID")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "String")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "BasicParentID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "BasicParent")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./basic-parent",
			Name:       "BasicParent",
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "ID")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_ForMany() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForMany",
			},
		},
	}
	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Basic")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "ID")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "String")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "BasicParentIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "BasicParents")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./basic-parent",
				Name:       "BasicParent",
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "ID")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasOne() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model1)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "BasicParent")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "ID")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "String")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "BasicID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "Basic")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./basic",
			Name:       "Basic",
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicParentIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "ID")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasMany() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"ID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r := registry.NewRegistry()
	r.SetModel("Basic", model0)
	r.SetModel("BasicParent", model1)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model1)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "BasicParent")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "ID")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "String")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "BasicIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "Basics")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./basic",
				Name:       "Basic",
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicParentIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "ID")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_NoModelName() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "morphe model has no name")

	suite.Nil(allTsObjects)
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "morphe model has no fields")

	suite.Nil(allTsObjects)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_NoIdentifiers() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{},
		Related:     map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "morphe model has no identifiers")

	suite.Nil(allTsObjects)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_StartHook_Successful() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelStart: func(config cfg.MorpheModelsConfig, model yaml.Model) (cfg.MorpheModelsConfig, yaml.Model, error) {
			if featureFlag != "otherName" {
				return config, model, nil
			}
			model.Name = model.Name + "CHANGED"
			delete(model.Fields, "Float")
			return config, model, nil
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]

	suite.Equal(tsObject0.Name, "BasicCHANGED")

	objectFields0 := tsObject0.Fields
	suite.Len(objectFields0, 1)

	objectFields00 := objectFields0[0]
	suite.Equal(objectFields00.Name, "UUID")
	suite.Equal(objectFields00.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicCHANGEDIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_StartHook_Failure() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelStart: func(config cfg.MorpheModelsConfig, model yaml.Model) (cfg.MorpheModelsConfig, yaml.Model, error) {
			if featureFlag != "otherName" {
				return config, model, nil
			}
			return config, model, fmt.Errorf("compile model start hook error")
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "compile model start hook error")
	suite.Nil(allTsObjects)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_SuccessHook_Successful() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelSuccess: func(allModelTypes []*tsdef.Object) ([]*tsdef.Object, error) {
			if featureFlag != "otherName" {
				return allModelTypes, nil
			}
			for _, modelObjectPtr := range allModelTypes {
				modelObjectPtr.Name = modelObjectPtr.Name + "CHANGED"
				newFields := []tsdef.ObjectField{}
				for _, modelStructField := range modelObjectPtr.Fields {
					if modelStructField.Name == "Float" {
						continue
					}
					newFields = append(newFields, modelStructField)
				}
				modelObjectPtr.Fields = newFields
			}
			return allModelTypes, nil
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]

	suite.Equal(tsObject0.Name, "BasicCHANGED")

	objectFields0 := tsObject0.Fields
	suite.Len(objectFields0, 1)

	objectFields00 := objectFields0[0]
	suite.Equal(objectFields00.Name, "UUID")
	suite.Equal(objectFields00.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimaryCHANGED")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_SuccessHook_Failure() {
	var featureFlag = "otherName"
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelSuccess: func(allModelObjects []*tsdef.Object) ([]*tsdef.Object, error) {
			if featureFlag != "otherName" {
				return allModelObjects, nil
			}
			return nil, fmt.Errorf("compile model success hook error")
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
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

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "compile model success hook error")
	suite.Nil(allTsObjects)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_FailureHook_NoModelIdentifiers() {
	modelHooks := hook.CompileMorpheModel{
		OnCompileMorpheModelFailure: func(config cfg.MorpheModelsConfig, model yaml.Model, compileFailure error) error {
			return fmt.Errorf("Model %s: %w", model.Name, compileFailure)
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{}

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{},
		Related:     map[string]yaml.ModelRelation{},
	}

	r := registry.NewRegistry()

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, model0)

	suite.NotNil(allTsObjectsErr)
	suite.ErrorContains(allTsObjectsErr, "Model Basic: morphe model has no identifiers")
	suite.Nil(allTsObjects)
}

