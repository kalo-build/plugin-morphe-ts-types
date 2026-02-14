package cfg

type MorpheModelsConfig struct {
	// FieldCasing specifies the casing for field names in generated model types.
	FieldCasing Casing
}

func (config MorpheModelsConfig) Validate() error {
	return nil
}
