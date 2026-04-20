# validate-refs

Validates that all YAML files in a directory tree are referenced in their parent `kustomization.yaml` files. This prevents orphaned files from accumulating in the repository.

## Purpose

When working with kustomize, it's easy to create YAML files but forget to add them to the `resources` section of `kustomization.yaml`. These orphaned files:
- Won't be deployed
- Create confusion about what's actually active
- Clutter the repository

This tool scans your kustomize directories and reports any YAML files that aren't referenced.

## Building

```bash
cd infra-tools
make build
```

The binary is placed at `infra-tools/bin/validate-refs`.

## Usage

### Basic validation

```bash
./bin/validate-refs --root /path/to/components
```

### With verbose output

```bash
./bin/validate-refs --root /path/to/components --verbose
```

Shows the number of directories checked.

## Exit codes

- `0`: All YAML files are properly referenced
- `1`: Found orphaned files or encountered an error

## Example output

### Success

```
Validating YAML references in: /path/to/components

✓ All YAML files are properly referenced in kustomization files
```

### Orphaned files found

```
Validating YAML references in: /path/to/components

✗ Found 2 orphaned YAML file(s):

  Directory: monitoring/blackbox-exporter/staging/private/stone-stage-p01/probes/https/
    - api-server-probe.yaml

  Directory: ca-bundle/production/
    - extra-ca.yaml

These files should be added to their respective kustomization.yaml files
or removed if they are no longer needed.
```

## What gets checked

The tool:
- ✅ Checks all YAML files (`.yaml` and `.yml`)
- ✅ Validates against `resources`, `patches`, `patchesStrategicMerge`, `patchesJson6902`, and `components`
- ❌ Ignores `kustomization.yaml` and `kustomization.yml` files themselves
- ❌ Ignores non-YAML files
- ❌ Skips directories without kustomization files

## CI Integration

Add to your CI pipeline to prevent orphaned files from being merged:

```yaml
- name: Validate kustomization references
  run: |
    ./infra-tools/bin/validate-refs --root ./components
```

## Fixing orphaned files

When the tool reports orphaned files, you have two options:

### Option 1: Add to kustomization.yaml

```yaml
resources:
  - existing-resource.yaml
  - orphaned-file.yaml  # Add the orphaned file
```

### Option 2: Remove if unused

```bash
git rm path/to/orphaned-file.yaml
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--root` | *required* | Root directory to validate |
| `--verbose` | `false` | Show count of checked directories |
| `--version` | — | Print version and exit |
