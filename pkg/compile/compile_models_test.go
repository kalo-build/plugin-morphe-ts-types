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
	suite.Equal(tsField00.Name, "autoIncrement")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "boolean")
	suite.Equal(tsField01.Type, tsdef.TsTypeBoolean)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "date")
	suite.Equal(tsField02.Type, tsdef.TsTypeDate)

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "float")
	suite.Equal(tsField03.Type, tsdef.TsTypeNumber)

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "integer")
	suite.Equal(tsField04.Type, tsdef.TsTypeNumber)

	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "protected")
	suite.Equal(tsField05.Type, tsdef.TsTypeString)

	tsField06 := tsFields0[6]
	suite.Equal(tsField06.Name, "sealed")
	suite.Equal(tsField06.Type, tsdef.TsTypeString)

	tsField07 := tsFields0[7]
	suite.Equal(tsField07.Name, "string")
	suite.Equal(tsField07.Type, tsdef.TsTypeString)

	tsField08 := tsFields0[8]
	suite.Equal(tsField08.Name, "time")
	suite.Equal(tsField08.Type, tsdef.TsTypeDate)

	tsField09 := tsFields0[9]
	suite.Equal(tsField09.Name, "uuid")
	suite.Equal(tsField09.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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
	suite.Equal(tsField00.Name, "autoIncrement")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "nationality")
	suite.Equal(tsField01.Type, tsdef.TsTypeObject{
		ModulePath: "../enums/nationality",
		Name:       "Nationality",
	})

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "uuid")
	suite.Equal(tsField02.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "string")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "basicParentID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "basicParent")
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
	suite.Equal(tsField10.Name, "id")
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
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "string")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "basicParentIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "basicParents")
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
	suite.Equal(tsField10.Name, "id")
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
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "string")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "basicID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "basic")
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
	suite.Equal(tsField10.Name, "id")
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
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "string")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "basicIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "basics")
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
	suite.Equal(tsField10.Name, "id")
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
	suite.Equal(objectFields00.Name, "uuid")
	suite.Equal(objectFields00.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicCHANGEDIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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
					if modelStructField.Name == "float" {
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
	suite.Equal(objectFields00.Name, "uuid")
	suite.Equal(objectFields00.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "BasicIDPrimaryCHANGED")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_ForOnePoly() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Title": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
	}

	articleModel := yaml.Model{
		Name: "Article",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Title": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
	}

	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Content": {
				Type: yaml.ModelFieldTypeString,
			},
			"CreatedAt": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)
	r.SetModel("Comment", commentModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, commentModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Comment")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 6)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "content")
	suite.Equal(tsField00.Type, tsdef.TsTypeString)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "createdAt")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "id")
	suite.Equal(tsField02.Type, tsdef.TsTypeNumber)

	// Polymorphic ID field
	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "commentableID")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	// Polymorphic type field
	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "commentableType")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	// Polymorphic object field (union type)
	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "commentable")
	suite.Equal(tsField05.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeUnion{
			Types: []tsdef.TsType{
				tsdef.TsTypeObject{
					ModulePath: "./post",
					Name:       "Post",
				},
				tsdef.TsTypeObject{
					ModulePath: "./article",
					Name:       "Article",
				},
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "CommentIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_ForManyPoly() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
	}

	companyModel := yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
	}

	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
			"Color": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Person", personModel)
	r.SetModel("Company", companyModel)
	r.SetModel("Tag", tagModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, tagModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Tag")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 6)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "color")
	suite.Equal(tsField00.Type, tsdef.TsTypeString)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "id")
	suite.Equal(tsField01.Type, tsdef.TsTypeNumber)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "name")
	suite.Equal(tsField02.Type, tsdef.TsTypeString)

	// Polymorphic ID field
	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "taggableIDs")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeString,
		},
	})

	// Polymorphic type field - not present for ForManyPoly?
	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "taggableType")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	// Polymorphic object field (union type array)
	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "taggables")
	suite.Equal(tsField05.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeUnion{
				Types: []tsdef.TsType{
					tsdef.TsTypeObject{
						ModulePath: "./person",
						Name:       "Person",
					},
					tsdef.TsTypeObject{
						ModulePath: "./company",
						Name:       "Company",
					},
				},
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "TagIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasOnePoly() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Content": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"FirstName": {
				Type: yaml.ModelFieldTypeString,
			},
			"LastName": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Comment": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Person", personModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, personModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Person")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 5)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "firstName")
	suite.Equal(tsField00.Type, tsdef.TsTypeString)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "id")
	suite.Equal(tsField01.Type, tsdef.TsTypeNumber)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "lastName")
	suite.Equal(tsField02.Type, tsdef.TsTypeString)

	// Has* polymorphic generates regular ID and struct fields
	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "commentID")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "comment")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./comment",
			Name:       "Comment",
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "PersonIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasManyPoly() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"FirstName": {
				Type: yaml.ModelFieldTypeString,
			},
			"LastName": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
				Aliased: "Tag",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Tag", tagModel)
	r.SetModel("Person", personModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, personModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Person")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 5)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "firstName")
	suite.Equal(tsField00.Type, tsdef.TsTypeString)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "id")
	suite.Equal(tsField01.Type, tsdef.TsTypeNumber)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "lastName")
	suite.Equal(tsField02.Type, tsdef.TsTypeString)

	// Has* polymorphic generates regular ID array and struct array fields
	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "tagIDs")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "tags")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./tag",
				Name:       "Tag",
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "PersonIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasOnePoly_WithoutAliased() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Content": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Title": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Comment": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				// No Aliased property - should use relationship name "Comment" as target model
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, postModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Post")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "title")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	// Has* polymorphic generates regular ID and struct fields using relationship name
	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "commentID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "comment")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./comment",
			Name:       "Comment",
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "PostIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_HasManyPoly_WithoutAliased() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Post"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Title": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Tag": {
				Type:    "HasManyPoly",
				Through: "Taggable",
				// No Aliased property - should use relationship name "Tag" as target model
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Tag", tagModel)
	r.SetModel("Post", postModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, postModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Post")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 4)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "title")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	// Has* polymorphic generates regular ID array and struct array fields using relationship name
	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "tagIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "tags")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./tag",
				Name:       "Tag",
			},
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "PostIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToTsObjects_Related_ForOne_Aliased() {
	modelHooks := hook.CompileMorpheModel{}

	modelsConfig := cfg.MorpheModelsConfig{}

	contactModel := yaml.Model{
		Name: "Contact",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Email": {
				Type: yaml.ModelFieldTypeString,
			},
			"Phone": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
	}

	personModel := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
			"PersonalContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Contact", contactModel)
	r.SetModel("Person", personModel)

	allTsObjects, allTsObjectsErr := compile.MorpheModelToTsObjects(modelHooks, modelsConfig, r, personModel)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Person")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 6)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "name")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	// PersonalContact relationship
	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "personalContactID")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "personalContact")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./contact",
			Name:       "Contact",
		},
	})

	// WorkContact relationship
	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "workContactID")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeNumber,
	})

	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "workContact")
	suite.Equal(tsField05.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeObject{
			ModulePath: "./contact",
			Name:       "Contact",
		},
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "PersonIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}
