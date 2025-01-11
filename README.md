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
name: Company
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
related:
  Employees:
    type: HasMany
    model: Person
```

The plugin generates:

``` typescript
export interface Company {
  ID: number;
  Name: string;
  EmployeeIDs?: number[];
  Employees?: Person[];
}
```

## Configuration

Reference implementation from compile_test.go:
startLine: 71
endLine: 78

## Development

1. Clone the repository
2. Install dependencies
3. Run tests: `go test ./...`

## License

MIT License
