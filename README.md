# Morphe TypeScript Types Plugin

A plugin for the Kalo CLI that converts Morphe registry definitions to TypeScript type definitions.

## Table of Contents

- [Usage](#usage)
  - [Configuration](#configuration)
  - [Output Structure](#output-structure)
- [Error Codes](#error-codes)
- [Development](#development)
  - [Building as WASM (WASI)](#building-as-wasm-wasi)
- [Overview](#overview)
- [Features](#features)
- [Example](#example)
- [Use as a Dependency](#use-as-a-dependency)
- [License](#license)

## Usage

This plugin is intended to be called by the Kalo CLI. It converts Morphe registry YAML files to TypeScript type definitions.

### Configuration

The plugin accepts a JSON config string with the following parameters:

```json
{
  "inputPath": "/path/to/morphe/registry",
  "outputPath": "/path/to/output/directory",
  "verbose": true,
  "config": {
    // Plugin configuration overrides (currently none)
  }
}
```

#### Parameters

- `inputPath` (required): Path to the Morphe registry directory.
- `outputPath` (required): Path where TypeScript type definitions will be generated.
- `verbose` (optional): Enable verbose logging for debugging. If not provided, defaults to 'false'.
- `config` (optional): Additional configuration options. If not provided, defaults apply.

### Output Structure

The plugin generates TypeScript definitions with the following structure:

```
outputPath/
  ├── enums/
  │   └── [enum-files].ts
  ├── models/
  │   └── [model-files].ts
  ├── structures/
  │   └── [structure-files].ts
  └── entities/
      └── [entity-files].ts
```

## Error Codes

| Code | Description |
|------|-------------|
| 1 | Compilation failed |
| 3 | Missing config |
| 4 | Invalid config JSON |
| 12 | Input path is required |
| 13 | Output path is required |

## Development

This plugin is designed to work in WASM (WASI) format when called by the Kalo CLI.

### Building as WASM (WASI)

The utility scripts `./scripts/build.bat` and `./scripts/build.sh` can be used to generate a build under `/dist`.

To build manually, run this command from the project root:

```bash
GOOS=wasip1 GOARCH=wasm go build -o ./dist/morphe-ts-types-v1.0.0.wasm ./cmd/plugin/main.go
```

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

## Use as a Dependency

This plugin can be installed as a typical Go dependency in other derivative plugin implementations:

```bash
go get github.com/kalo-build/plugin-morphe-ts-types
```

You can then use it as a base language compiler to accelerate development and customize its behavior by implementing hook functions and assigning them to the configuration:

```go
import (
    "github.com/kalo-build/plugin-morphe-ts-types/pkg/compile"
    "github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/cfg"
    "github.com/kalo-build/plugin-morphe-ts-types/pkg/compile/hook"
    "github.com/kalo-build/morphe-go/pkg/yaml"
    rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

// Configure full paths
config := compile.DefaultMorpheCompileConfig(
    "path/to/morphe/registry",
    "path/to/output/directory",
)

// Or customize the configuration manually
config = compile.MorpheCompileConfig{
    MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
        RegistryEnumsDirPath:      "path/to/enums",
        RegistryModelsDirPath:     "path/to/models",
        RegistryStructuresDirPath: "path/to/structures",
        RegistryEntitiesDirPath:   "path/to/entities",
    },
    // Other configuration options...
}

// Add custom hooks to modify behavior
config.EnumHooks.OnCompileMorpheEnumStart = func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
    // Modify the enum or config before compilation
    if enum.Name == "Status" {
        // Rename the enum
        enum.Name = enum.Name + "Type"
        // Remove specific entries
        delete(enum.Entries, "Inactive")
    }
    return config, enum, nil
}

// Inside your plugin's compile implementation, call the ts types plugin's main compile function with your customized config
err := compile.MorpheToTypescript(config)
```

> **Note:** This integration pattern is experimental and may change or be removed in the near future.

## License

MIT License
