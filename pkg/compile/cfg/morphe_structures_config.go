package cfg

type MorpheStructuresConfig struct {
	// FieldCasing specifies the casing for field names in generated structure types.
	FieldCasing Casing
}

func (config MorpheStructuresConfig) Validate() error {
	return nil
}
