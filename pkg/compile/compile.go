package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToTypescriptObjects(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	allModelObjectDefs, compileAllErr := AllMorpheModelsToTsObjects(config, r)
	if compileAllErr != nil {
		return compileAllErr
	}

	_, writeAllErr := WriteAllModelObjectDefinitions(config, allModelObjectDefs)
	if writeAllErr != nil {
		return writeAllErr
	}
	return nil
}
