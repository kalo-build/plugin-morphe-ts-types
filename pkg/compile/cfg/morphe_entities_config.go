package cfg

type MorpheEntitiesConfig struct {
	// FieldCasing specifies the casing for field names in generated entity types.
	FieldCasing Casing
}

func (config MorpheEntitiesConfig) Validate() error {
	return nil
}
