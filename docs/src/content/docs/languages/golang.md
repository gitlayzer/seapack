---
title: Go
description: Building Go applications with SeaPack
---

SeaPack builds and deploys Go applications with static binary compilation.

## Detection

Your project will be detected as a Go application if any of these conditions are
met:

- A `go.mod` file exists in the root directory
- A `go.work` file exists in the root directory (Go workspaces)
- A `main.go` file exists in the root directory

## Versions

The Go version is determined in the following order:

- Any mise-supported version file (`mise.toml`, `.tool-versions`, etc)
- Read from the `go.mod` file
- Set via the `SEAPACK_GO_VERSION` environment variable
- Defaults to `1.23`

## Configuration

SeaPack builds your Go application as a static binary by default. The build
process:

- Installs Go dependencies
- Builds your application with optimized flags (`-ldflags="-w -s"`)
- Names the output binary `out`

SeaPack determines the main package to build in the following order:

1. The module specified by the `SEAPACK_GO_WORKSPACE_MODULE` environment variable (for workspaces)
2. The package specified by the `SEAPACK_GO_BIN` environment variable
3. The root directory if it contains Go files
4. The first subdirectory in the `cmd/` directory
5. For workspaces: the first module containing a `main.go` file
6. The `main.go` file in the root directory

### Config Variables

| Variable                       | Description                                  | Example  |
| ------------------------------ | -------------------------------------------- | -------- |
| `SEAPACK_GO_VERSION`          | Override the Go version                      | `1.22`   |
| `SEAPACK_GO_BIN`              | Specify which command in cmd/ to build       | `server` |
| `SEAPACK_GO_WORKSPACE_MODULE` | Specify which workspace module to build      | `api`    |
| `CGO_ENABLED`                  | Enable CGO for non-static binary compilation | `1`      |

### Go Workspaces

SeaPack supports Go workspaces (introduced in Go 1.18) for multi-module projects:

- Detects projects with a `go.work` file at the root
- Automatically discovers and copies all module dependencies
- Builds the first module with a `main.go` file by default
- Use `SEAPACK_GO_WORKSPACE_MODULE` to specify which module to build

Example workspace structure:

```
├── go.work
├── api/
│   ├── go.mod
│   └── main.go
└── shared/
    ├── go.mod
    └── lib.go
```

To build a specific module:

```bash
SEAPACK_GO_WORKSPACE_MODULE=api
```

### CGO Support

By default, SeaPack builds static binaries with `CGO_ENABLED=0`. If you need
CGO support:

- Set the `CGO_ENABLED` environment variable to `1`
- SeaPack will include the necessary build dependencies (gcc, g++, libc6-dev)
- The runtime image will include libc6 for dynamic linking

## BuildKit Caching

The Go provider will cache `~/.cache/go-build` under the cache key `go-build`.