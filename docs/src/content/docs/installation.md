---
title: Installation
description: How to install SeaPack
---

SeaPack is available as a CLI tool. The latest release is available [on
GitHub](https://github.com/gitlayzer/seapack/releases).

The BuildKit frontend is available as a Docker Hub image:
`ghcr.io/gitlayzer/seapack-frontend`.

## Mise

We love mise, and you can install SeaPack using mise:

```sh
mise use github:gitlayzer/seapack@latest
```

## Curl

Download SeaPack from GitHub releases and install automatically.

```sh
curl -sSL https://raw.githubusercontent.com/gitlayzer/seapack/refs/heads/main/install.sh | sh
```

You can also customize the version, destination, and other config options:

```sh
curl -sSL https://raw.githubusercontent.com/gitlayzer/seapack/refs/heads/main/install.sh | SEAPACK_VERSION=0.2.3 sh -s -- --bin-dir ~/.local/bin
```

## GitHub Releases

Go to the [latest release](https://github.com/gitlayzer/seapack/releases) and
download the `seapack` binary for your platform.

## From Source

```sh
git clone https://github.com/gitlayzer/seapack.git
cd seapack
go build -o seapack ./cmd/...

./seapack --help
```

## Supported Platforms

Linux and MacOS are supported.

Windows builds are generated but not officially supported. That being said, PRs are welcome to fix any Windows-specific bugs.

## Help

Need help? Check out our [Help page](/help) for support options.
