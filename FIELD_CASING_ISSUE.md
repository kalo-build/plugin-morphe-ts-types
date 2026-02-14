# Field Casing Issue - Investigate Before Commit

## Status
**Blocked** - `fieldCasing` configuration added but not confirmed working end-to-end via WASM plugin.

## What Was Done
1. Added `cfg.Casing` type (`pkg/compile/cfg/casing.go`) with `Apply()` method supporting `camel`, `snake`, `pascal` conventions
2. Added `FieldCasing cfg.Casing` to `MorpheCompileConfig` (top-level)
3. Threaded `fieldCasing` through all compile functions:
   - `compile_model_fields.go` - all field generation + relationship fields
   - `compile_structure_fields.go` - structure field generation
   - `compile_entity_fields.go` - all entity field generation + relationship fields
4. Wired `fieldCasing` parsing in `cmd/plugin/main.go` from `compileConfig.Config["fieldCasing"]`
5. Exposed `fieldCasing` in `plugin.yaml` configSchema (enum: camel/snake/pascal, default: camel)
6. All unit tests pass with `go test -count=1 ./...`

## What Didn't Work
When the WASM plugin was rebuilt, installed, and `kalo compile` was run in the client directory with `fieldCasing: "snake"` configured in `client/kalo.yaml`, the generated TypeScript types still showed camelCase output.

## Possible Root Causes

### Build / Install Issue
- WASM build may not have picked up the changes (stale cache, missing `replace` directive for local deps, etc.)
- The installed `.wasm` file in `.kalo/plugins/` may be stale
- The `kalo.lock` may be pinning an old version

### Logical / Flow Bug
- Code review did not find issues - the wiring chain looks correct:
  - `main.go` parses `fieldCasing` → `morpheConfig.FieldCasing`
  - `AllMorpheModelsToTsObjects` copies `config.FieldCasing` → `modelsConfig.FieldCasing`
  - `AllMorpheStructuresToTsObjects` copies `config.FieldCasing` → `structuresConfig.FieldCasing`
  - `AllMorpheEntitiesToTsObjects` copies `config.FieldCasing` → `entitiesConfig.FieldCasing`
  - Each section passes `fieldCasing` to field generation functions
- Possible issue: CLI may not be passing `config` block to the plugin at all (check how `kalo-cli` marshals config to WASM plugins)

## Investigation Steps
1. Add debug logging in `main.go` to print the raw `compileConfig.Config` map - verify `fieldCasing` is present
2. Verify the WASM binary is fresh: check file timestamp of `.kalo/plugins/kalo-build-plugin-morphe-ts-types-v1.0.0.wasm`
3. Check if the CLI actually passes the `config` section from `kalo.yaml` to the plugin invocation
4. Try running the WASM plugin directly with a hand-crafted JSON config to isolate CLI vs plugin issue

## Note on Casing Strategy
After investigation, the decision was made to keep TS types in **camelCase** (default) and only use `fieldCasing: "snake"` for the Zod schemas plugin. The TS types represent the application-facing types developers write against. A deserialization bridge (see `MORPHEMAP_BRAINSTORM.md` in kalo-plugin-registry) would handle the snake_case API → camelCase app type conversion.

The `fieldCasing` capability is still valuable for projects that prefer snake_case TypeScript (e.g., matching a Python backend), so the feature should still work correctly.
