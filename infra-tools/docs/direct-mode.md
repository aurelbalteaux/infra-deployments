# Direct Mode for render-diff

The `--detection-mode=direct` flag enables render-diff to work with simpler repository structures that don't use ArgoCD ApplicationSets for component discovery.

## When to use direct mode

Use direct mode when your repository:
- Has a simple components directory structure (e.g., `components/foo/staging/`, `components/foo/production/`)
- Does not use ArgoCD ApplicationSets
- Has kustomization files directly in component directories

## Usage

```bash
./bin/render-diff \
  --detection-mode=direct \
  --components-dir=components \
  --repo-root=/path/to/repo
```

## Flags specific to direct mode

| Flag | Default | Description |
|------|---------|-------------|
| `--detection-mode` | `appset` | Set to `direct` to enable direct mode |
| `--components-dir` | `components` | Path to components directory relative to repo root |

All other render-diff flags work the same in direct mode (`--color`, `--open`, `--output-dir`, etc.).

## How it works

1. For each changed file, walks up the directory tree to find the closest `kustomization.yaml`
2. Infers the environment (development/staging/production) from the directory path
3. Extracts cluster information if present (e.g., `production/private/kflux-ocp-p01`)
4. Runs the standard render-diff engine to compute diffs

## Example directory structure

```
components/
├── monitoring/
│   └── blackbox-exporter/
│       ├── base/
│       │   └── kustomization.yaml
│       ├── staging/
│       │   └── kustomization.yaml
│       └── production/
│           ├── base/
│           │   └── kustomization.yaml
│           └── private/
│               └── kflux-ocp-p01/
│                   └── kustomization.yaml
└── ca-bundle/
    ├── base/
    │   └── kustomization.yaml
    ├── development/
    │   └── kustomization.yaml
    └── production/
        └── kustomization.yaml
```

When a file in `components/monitoring/blackbox-exporter/production/private/kflux-ocp-p01/` changes, direct mode:
- Finds the kustomization in that directory
- Detects environment = `production`
- Detects cluster = `kflux-ocp-p01`
- Runs diff for that component

## CI Integration

Direct mode works with all CI output modes:

```bash
./bin/render-diff \
  --detection-mode=direct \
  --components-dir=components \
  --output-mode=ci-summary,ci-comment,ci-artifact-dir \
  --output-dir=../render-diff-output
```

## Comparison with appset mode

| Feature | appset mode (default) | direct mode |
|---------|----------------------|-------------|
| Component discovery | Via ArgoCD ApplicationSets | Via directory traversal |
| Repository structure | Requires `argo-cd-apps/overlays/` | Simple `components/` directory |
| Environment detection | From overlay names | From directory names in path |
| Cluster detection | From ApplicationSet generators | From directory names in path |
