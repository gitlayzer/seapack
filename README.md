# SeaPack

[![CI](https://github.com/gitlayzer/seapack/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/gitlayzer/seapack/actions/workflows/unit-tests.yml)
[![Run Tests](https://github.com/gitlayzer/seapack/actions/workflows/integration-tests.yml/badge.svg)](https://github.com/gitlayzer/seapack/actions/workflows/integration-tests.yml)

SeaPack is a Sealos-oriented tool for building runnable container images from
source code with minimal configuration. It focuses on the core languages used
by Sealos app workloads: Node.js, Python, Go, Java, and Deno.

## Getting Started

```bash
# Install seapack on Linux amd64/arm64
curl -sSL https://raw.githubusercontent.com/gitlayzer/seapack/refs/heads/main/install.sh | sh

# start BuildKit container & let seapack know about it
docker run --rm --privileged -d --name buildkit ghcr.io/gitlayzer/seapack-buildkit:latest
export BUILDKIT_HOST='docker-container://buildkit'

# create a Next.js app
npm create next-app@latest my-app
cd my-app

# build and run the app!
seapack build .
docker run -p 3000:3000 -it my-app
```

SeaPack automatically detects the project type, generates an optimized build
plan, and builds a container image with the unified SeaPack base images.

**Note:** The above steps are for running SeaPack locally to experiment and
test. In a Sealos build pipeline, SeaPack is intended to run as the app image
builder behind the deployment flow.

## Documentation

SeaPack documentation is written for Sealos operators and developers who need a
small, predictable source-to-image builder for Sealos workloads.

## Contributing

SeaPack is open source and open to contributions. See the
[CONTRIBUTING.md](CONTRIBUTING.md) file for more information.
