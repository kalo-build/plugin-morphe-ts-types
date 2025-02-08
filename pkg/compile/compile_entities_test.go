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

type CompileEntitiesTestSuite struct {
	suite.Suite
}

func TestCompileEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEntitiesTestSuite))
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects() {
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
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
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

	entity1 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetEntity("BasicParent", entity1)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "User")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 8)

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
	suite.Equal(tsField05.Name, "String")
	suite.Equal(tsField05.Type, tsdef.TsTypeString)

	tsField06 := tsFields0[6]
	suite.Equal(tsField06.Name, "Time")
	suite.Equal(tsField06.Type, tsdef.TsTypeDate)

	tsField07 := tsFields0[7]
	suite.Equal(tsField07.Name, "UUID")
	suite.Equal(tsField07.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "UserIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_NoEntityName() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "",
		Fields: map[string]yaml.EntityField{
			"AutoIncrement": {
				Type: "User.AutoIncrement",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "entity has no name")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_NoFields() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name:   "Basic",
		Fields: map[string]yaml.EntityField{},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "morphe entity Basic has no fields")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_NoIdentifiers() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{},
		Related:     map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "entity 'Basic' has no identifiers")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_EnumField() {
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
			"Nationality": {
				Type: "User.Nationality",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
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
			"Nationality": {
				Type: "Nationality",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}
	r.SetEnum("Nationality", enum0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "User")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 2)

	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "Nationality")
	suite.Equal(tsField00.Type, tsdef.TsTypeObject{
		ModulePath: "../enums/nationality",
		Name:       "Nationality",
	})

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "UUID")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "UserIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "UUID")
	suite.Equal(tsField10.Type, tsdef.TsTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_EnumField_EnumNotFound() {
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
			"Nationality": {
				Type: "User.Nationality",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
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
			"Nationality": {
				Type: "Nationality",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "morphe entity 'User' field 'Nationality' has unknown non-primitive type 'Nationality'")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_ForOne() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}

	r := registry.NewRegistry()

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetEntity("BasicParent", entity1)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_ForMany() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForMany",
			},
		},
	}

	r := registry.NewRegistry()

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForMany",
			},
		},
	}
	r.SetModel("Basic", model0)

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetEntity("BasicParent", entity1)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_HasOne() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}

	r := registry.NewRegistry()

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetEntity("Basic", entity1)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_HasMany() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}

	r := registry.NewRegistry()

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

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
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "ForOne",
			},
		},
	}
	r.SetEntity("Basic", entity1)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_StartHook_Successful() {
	var hookCalled = false
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheEntitiesConfig, entity yaml.Entity) (cfg.MorpheEntitiesConfig, yaml.Entity, error) {
			hookCalled = true
			return config, entity, nil
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.True(hookCalled)
	suite.Nil(tsObjectErr)
	suite.Len(allTsObjects, 2)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_StartHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheEntitiesConfig, entity yaml.Entity) (cfg.MorpheEntitiesConfig, yaml.Entity, error) {
			if featureFlag != "otherName" {
				return config, entity, nil
			}
			return config, entity, fmt.Errorf("compile entity start hook error")
		},
	}

	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "compile entity start hook error")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_SuccessHook_Successful() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(allEntityObjects []*tsdef.Object) ([]*tsdef.Object, error) {
			if featureFlag != "otherName" {
				return allEntityObjects, nil
			}
			for _, entityObject := range allEntityObjects {
				entityObject.Name = entityObject.Name + "CHANGED"
			}
			return allEntityObjects, nil
		},
	}

	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "Basic.UUID",
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

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
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.Nil(tsObjectErr)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_SuccessHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(allEntityObjects []*tsdef.Object) ([]*tsdef.Object, error) {
			if featureFlag != "otherName" {
				return allEntityObjects, nil
			}
			return nil, fmt.Errorf("compile entity success hook error")
		},
	}

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
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
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
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "compile entity success hook error")
	suite.Nil(allTsObjects)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_FailureHook_UnknownRootModel() {
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityFailure: func(config cfg.MorpheEntitiesConfig, entity yaml.Entity, compileFailure error) error {
			return fmt.Errorf("Entity %s: %w", entity.Name, compileFailure)
		},
	}

	entitiesConfig := cfg.MorpheEntitiesConfig{}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "NonExistentModel.UUID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	allTsObjects, tsObjectErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, entity0)

	suite.NotNil(tsObjectErr)
	suite.ErrorContains(tsObjectErr, "Entity User: morphe entity User field UUID references unknown root model: NonExistentModel")
	suite.Nil(allTsObjects)
}
