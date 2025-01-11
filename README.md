# plugin-morphe-ts-types

A TypeScript type definition generator plugin for [Morphe specifications](https://github.com/kaloseia/morphe).

## Overview

This plugin generates TypeScript type definitions (.d.ts files) from Morphe specifications by converting:

- Models to TypeScript interfaces
- Entities to TypeScript interfaces
- Enums to TypeScript enums
- Structures to TypeScript types

## Features

- Generates strongly-typed TypeScript definitions
- Supports all Morphe relationship types:
  - HasOne/HasMany for model relationships
  - ForOne/ForMany for entity relationships
- Handles field types:
  - All primitive types (string, number, boolean, etc.)
  - Date fields
  - UUID fields
  - Custom enum types
- Preserves model identifiers and field attributes
- Configurable output paths
- Extensible via hooks system

## Example

Given a Morphe model:

``` yaml
name: Person
fields:
  ID:
    type: AutoIncrement
    attributes:
      - mandatory
  FirstName:
    type: String
  LastName:
    type: String
  Nationality:
    type: Nationality
identifiers:
  primary: ID
  name:
    - FirstName
    - LastName
related:
  ContactInfo:
    type: HasOne
  Company:
    type: ForOne
```

The plugin generates:

``` typescript
import { Nationality } from "../enums/nationality"
import { Company } from "./company"
import { ContactInfo } from "./contact-info"

export type Person = {
	firstName: string
	id: number
	lastName: string
	nationality: Nationality
	companyID?: number
	company?: Company
	contactInfoID?: number
	contactInfo?: ContactInfo
}

export type PersonIDName = {
	firstName: string
	lastName: string
}

export type PersonIDPrimary = {
	id: number
}
```

See [testdata](https://github.com/kaloseia/plugin-morphe-ts-types/tree/main/testdata) for more examples.

## Configuration / Usage

First define your Morphe registry for all your enums, structures, models and entities.

Then install the plugin and add your paths to the plugin configuration.

Example usage:

``` go
config := compile.MorpheCompileConfig{
    MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
        RegistryEnumsDirPath:      "path/to/enums",
        RegistryStructuresDirPath: "path/to/structures",
        RegistryModelsDirPath:     "path/to/models",
        RegistryEntitiesDirPath:   "path/to/entities",
    },

    MorpheEnumsConfig: cfg.MorpheEnumsConfig{},
    EnumWriter: &compile.MorpheEnumFileWriter{
        TargetDirPath: "path/to/enums",
    },

    MorpheModelsConfig: cfg.MorpheModelsConfig{},
    ModelWriter: &compile.MorpheObjectFileWriter{
        TargetDirPath: "path/to/models",
    },

    MorpheEntitiesConfig: cfg.MorpheEntitiesConfig{},
    EntityWriter: &compile.MorpheObjectFileWriter{
        TargetDirPath: "path/to/entities",
    },

    MorpheStructuresConfig: cfg.MorpheStructuresConfig{},
    StructureWriter: &compile.MorpheObjectFileWriter{
        TargetDirPath: "path/to/structures",
    },
}

compileErr := compile.MorpheToTypescript(config)
```

## Development

1. Clone the repository
2. Install dependencies
3. Run tests: `go test ./...`

## License

MIT License
