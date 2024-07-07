package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToTypescriptTypes(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	_, compileAllErr := AllMorpheModelsToTsObjects(config, r)
	if compileAllErr != nil {
		return compileAllErr
	}
	return nil
}
