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
	suite.Equal(tsField05.Name, "string")
	suite.Equal(tsField05.Type, tsdef.TsTypeString)

	tsField06 := tsFields0[6]
	suite.Equal(tsField06.Name, "time")
	suite.Equal(tsField06.Type, tsdef.TsTypeDate)

	tsField07 := tsFields0[7]
	suite.Equal(tsField07.Name, "uuid")
	suite.Equal(tsField07.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "UserIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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
	suite.Equal(tsField00.Name, "nationality")
	suite.Equal(tsField00.Type, tsdef.TsTypeObject{
		ModulePath: "../enums/nationality",
		Name:       "Nationality",
	})

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "uuid")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "UserIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "uuid")
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_ForOnePoly() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Comment model with polymorphic ForOnePoly relationship
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
				Type: yaml.ModelFieldTypeTime,
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
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
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
		Related: map[string]yaml.ModelRelation{},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Content": {
				Type: "Comment.Content",
			},
			"CreatedAt": {
				Type: "Comment.CreatedAt",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"Name": {
				Type: "Person.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	companyEntity := yaml.Entity{
		Name: "Company",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Company.ID",
			},
			"Name": {
				Type: "Company.Name",
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
	r.SetModel("Comment", commentModel)
	r.SetModel("Person", personModel)
	r.SetModel("Company", companyModel)
	r.SetEntity("Comment", commentEntity)
	r.SetEntity("Person", personEntity)
	r.SetEntity("Company", companyEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, commentEntity)

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
	suite.Equal(tsField01.Type, tsdef.TsTypeDate)

	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "id")
	suite.Equal(tsField02.Type, tsdef.TsTypeNumber)

	// ForOnePoly generates ID, Type, and union fields
	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "commentableID")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "commentableType")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "commentable")
	suite.Equal(tsField05.Type, tsdef.TsTypeOptional{
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
	})

	tsObject1 := allTsObjects[1]
	suite.Equal(tsObject1.Name, "CommentIDPrimary")

	tsFields1 := tsObject1.Fields
	suite.Len(tsFields1, 1)

	tsField10 := tsFields1[0]
	suite.Equal(tsField10.Name, "id")
	suite.Equal(tsField10.Type, tsdef.TsTypeNumber)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_ForManyPoly() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Tag model with polymorphic ForManyPoly relationship
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
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
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
		Related: map[string]yaml.ModelRelation{},
	}

	// Tag entity
	tagEntity := yaml.Entity{
		Name: "Tag",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Tag.ID",
			},
			"Name": {
				Type: "Tag.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"Name": {
				Type: "Person.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	companyEntity := yaml.Entity{
		Name: "Company",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Company.ID",
			},
			"Name": {
				Type: "Company.Name",
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
	r.SetModel("Tag", tagModel)
	r.SetModel("Person", personModel)
	r.SetModel("Company", companyModel)
	r.SetEntity("Tag", tagEntity)
	r.SetEntity("Person", personEntity)
	r.SetEntity("Company", companyEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, tagEntity)

	suite.Nil(allTsObjectsErr)
	suite.Len(allTsObjects, 2)

	tsObject0 := allTsObjects[0]
	suite.Equal(tsObject0.Name, "Tag")

	tsFields0 := tsObject0.Fields
	suite.Len(tsFields0, 5)

	// Regular fields
	tsField00 := tsFields0[0]
	suite.Equal(tsField00.Name, "id")
	suite.Equal(tsField00.Type, tsdef.TsTypeNumber)

	tsField01 := tsFields0[1]
	suite.Equal(tsField01.Name, "name")
	suite.Equal(tsField01.Type, tsdef.TsTypeString)

	// ForManyPoly generates IDs, Type, and union array fields
	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "taggableIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeString,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "taggableType")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeString,
	})

	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "taggables")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_HasOnePoly() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Comment model with polymorphic inverse
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

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Content": {
				Type: "Comment.Content",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"FirstName": {
				Type: "Person.FirstName",
			},
			"LastName": {
				Type: "Person.LastName",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
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
	r.SetEntity("Comment", commentEntity)
	r.SetEntity("Person", personEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, personEntity)

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

	// HasOnePoly generates regular ID and object fields using relationship name
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_HasManyPoly() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Tag model with polymorphic inverse
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

	// Tag entity
	tagEntity := yaml.Entity{
		Name: "Tag",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Tag.ID",
			},
			"Name": {
				Type: "Tag.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"FirstName": {
				Type: "Person.FirstName",
			},
			"LastName": {
				Type: "Person.LastName",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
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
	r.SetEntity("Tag", tagEntity)
	r.SetEntity("Person", personEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, personEntity)

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

	// HasManyPoly generates regular ID array and object array fields using relationship name
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_ForOne_Aliased() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Person model with multiple aliased relationships to Contact
	contactModel := yaml.Model{
		Name: "Contact",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Email": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
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
			"PersonalContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	// Contact entity
	contactEntity := yaml.Entity{
		Name: "Contact",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Contact.ID",
			},
			"Email": {
				Type: "Contact.Email",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"Name": {
				Type: "Person.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"PersonalContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Contact", contactModel)
	r.SetModel("Person", personModel)
	r.SetEntity("Contact", contactEntity)
	r.SetEntity("Person", personEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, personEntity)

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

	// PersonalContact relationship (aliased to Contact)
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

	// WorkContact relationship (also aliased to Contact)
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

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToTsObjects_Related_HasMany_Aliased() {
	entityHooks := hook.CompileMorpheEntity{}
	entitiesConfig := cfg.MorpheEntitiesConfig{}

	// Project model
	projectModel := yaml.Model{
		Name: "Project",
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
		Related: map[string]yaml.ModelRelation{},
	}

	// Person model with multiple aliased HasMany relationships to Project
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
			"PersonalProject": {
				Type:    "HasMany",
				Aliased: "Project",
			},
			"WorkProject": {
				Type:    "HasMany",
				Aliased: "Project",
			},
		},
	}

	// Project entity
	projectEntity := yaml.Entity{
		Name: "Project",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Project.ID",
			},
			"Title": {
				Type: "Project.Title",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	personEntity := yaml.Entity{
		Name: "Person",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Person.ID",
			},
			"Name": {
				Type: "Person.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"PersonalProject": {
				Type:    "HasMany",
				Aliased: "Project",
			},
			"WorkProject": {
				Type:    "HasMany",
				Aliased: "Project",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Project", projectModel)
	r.SetModel("Person", personModel)
	r.SetEntity("Project", projectEntity)
	r.SetEntity("Person", personEntity)

	allTsObjects, allTsObjectsErr := compile.MorpheEntityToTsObjects(entityHooks, entitiesConfig, r, personEntity)

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

	// PersonalProject relationship (aliased to Project)
	tsField02 := tsFields0[2]
	suite.Equal(tsField02.Name, "personalProjectIDs")
	suite.Equal(tsField02.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField03 := tsFields0[3]
	suite.Equal(tsField03.Name, "personalProjects")
	suite.Equal(tsField03.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./project",
				Name:       "Project",
			},
		},
	})

	// WorkProject relationship (also aliased to Project)
	tsField04 := tsFields0[4]
	suite.Equal(tsField04.Name, "workProjectIDs")
	suite.Equal(tsField04.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeNumber,
		},
	})

	tsField05 := tsFields0[5]
	suite.Equal(tsField05.Name, "workProjects")
	suite.Equal(tsField05.Type, tsdef.TsTypeOptional{
		ValueType: tsdef.TsTypeArray{
			ValueType: tsdef.TsTypeObject{
				ModulePath: "./project",
				Name:       "Project",
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
