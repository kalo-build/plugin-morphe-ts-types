package tsdef

var (
	TsTypeString = TsTypePrimitive{
		Syntax: "string",
	}
	TsTypeNumber = TsTypePrimitive{
		Syntax: "number",
	}
	TsTypeBoolean = TsTypePrimitive{
		Syntax: "boolean",
	}
	TsTypeDate = TsTypeObject{
		Name: "Date",
	}
)
