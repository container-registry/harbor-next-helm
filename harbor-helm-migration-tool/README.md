# Harbor Helm Migration Tool

A CLI tool to migrate Harbor Helm chart values from the legacy [harbor-helm](https://github.com/goharbor/harbor-helm) format to the new [harbor-next-helm](https://github.com/goharbor/harbor-next-helm) format.

## Overview

The `harbor-migrate` command transforms your existing `values.yaml` file to be compatible with the new Harbor Helm chart. It handles:

- Field renames and structural changes
- Type conversions (e.g., port strings to integers)
- Component configuration migrations
- Detection of unsupported features with actionable warnings

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/goharbor/harbor-next-helm.git
cd harbor-next-helm/harbor-helm-migration-tool

# Build
go build -o harbor-migrate ./cmd/harbor-migrate

# Or install to GOPATH/bin
go install ./cmd/harbor-migrate
```

### Using Task

```bash
cd harbor-next-helm
task migrate:build
# Binary will be at harbor-helm-migration-tool/bin/harbor-migrate
```

## Usage

```bash
# Preview migration (dry-run mode)
harbor-migrate --dry-run legacy-values.yaml

# Migrate and save to new file
harbor-migrate legacy-values.yaml new-values.yaml

# Migrate with verbose output
harbor-migrate --verbose legacy-values.yaml new-values.yaml

# Strict mode - exit with error if any warnings are generated
harbor-migrate --strict legacy-values.yaml new-values.yaml
```

### Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Print migrated values to stdout without writing file |
| `-v, --verbose` | Enable verbose output |
| `--strict` | Exit with error if any warnings are generated |
| `--version` | Show version information |
| `-h, --help` | Show help |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (no errors) |
| 1 | Migration completed but has errors requiring manual intervention |
| 2 | Fatal error (e.g., file not found, parse error) |

## Migration Report

The tool generates a migration report showing any issues found:

```
=== Migration Report ===
Errors: 1, Warnings: 2, Info: 1

[ERROR] database.type: Internal database is not supported in harbor-next-helm.
  Action: Deploy an external PostgreSQL instance and configure database.host

[WARNING] trivy: Trivy is now provided as a subchart.
  Action: Configure Trivy options under trivy.* using the subchart format

[INFO] redis.type: Internal Redis has been migrated to Valkey subchart
  Action: Review valkey.* settings in output
```

### Warning Levels

- **ERROR**: Blocking issues that require manual intervention before the chart will work
- **WARNING**: Features that need attention but won't prevent deployment
- **INFO**: Informational messages about changes made

## Key Mappings

### Ingress

| Legacy | New |
|--------|-----|
| `expose.ingress.hosts.core` | `ingress.hosts[0].host` |
| `expose.ingress.className` | `ingress.className` |
| `expose.ingress.annotations` | `ingress.annotations` |
| `expose.tls.secret.secretName` | `ingress.tls[0].secretName` |

### Database

| Legacy | New |
|--------|-----|
| `database.external.host` | `database.host` |
| `database.external.port` | `database.port` (string to int) |
| `database.external.username` | `database.username` |
| `database.external.coreDatabase` | `database.database` |

### Redis

| Legacy | New |
|--------|-----|
| `redis.type: internal` | `valkey.enabled: true` |
| `redis.external.addr` | `externalRedis.host` + `externalRedis.port` |
| `redis.external.password` | `externalRedis.password` |

### Storage

| Legacy | New |
|--------|-----|
| `persistence.imageChartStorage.type` | `registry.storage.type` |
| `persistence.imageChartStorage.s3.*` | `registry.storage.s3.*` |
| `persistence.persistentVolumeClaim.registry.size` | `registry.persistence.size` |

### Components

For each component (core, portal, registry, jobservice, exporter):

| Legacy | New |
|--------|-----|
| `{component}.extraEnvVars` | `{component}.extraEnv` |
| `{component}.image.repository` | `{component}.image.repository` |
| `{component}.replicas` | `{component}.replicas` |

## Unsupported Features

The following features from harbor-helm are not supported in harbor-next-helm:

| Feature | Reason | Action Required |
|---------|--------|-----------------|
| `database.type: internal` | No internal database | Use external PostgreSQL |
| `nginx.*` | Ingress used directly | Configure `ingress.*` |
| `notary.*` | Deprecated | Use Sigstore/cosign |
| `chartmuseum.*` | Deprecated | Use OCI artifacts |
| `expose.type: nodePort/loadBalancer/clusterIP` | Simplified exposure | Use ingress or Gateway API |
| Swift storage | Not supported | Use s3, azure, gcs, or oss |

## Development

### Running Tests

```bash
go test -v ./pkg/migrate/...
```

### Project Structure

```
harbor-helm-migration-tool/
├── cmd/
│   └── harbor-migrate/
│       └── main.go           # CLI entry point
├── pkg/
│   └── migrate/
│       ├── legacy_types.go   # Legacy values structs
│       ├── new_types.go      # New values structs
│       ├── mappings.go       # Field mapping logic
│       ├── warnings.go       # Warning generation
│       ├── migrate.go        # Main orchestration
│       └── migrate_test.go   # Unit tests
├── go.mod
├── go.sum
└── README.md
```

## License

Apache License 2.0
