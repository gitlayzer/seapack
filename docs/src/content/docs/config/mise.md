---
title: Mise Configuration
description: How to customize your image using Mise configuration
---

SeaPack is built on top of [Mise](https://mise.jdx.dev/). You can use the various mise configuration options
to customize the SeaPack-generated image. For instance, you can set environment variables, pin language
versions, and add additional utilities like `jq` to your image all through the mise configuration toml.

## Philosophy

* We use the latest mise version. There is automated tooling setup to ensure
  the mise version on the latest SeaPack version is no more than a couple weeks
  out of date.
* SeaPack assumes the default Mise configuration options unless a Sealos build
  requirement needs a stricter default.
* SeaPack generates global mise configuration based on analyizing the application
  source code. However, this global configuration is set in `/etc/mise/config.toml`
  so it can easily be overwritten in your application.

## Default Settings

SeaPack sets the following mise settings by default in the generated
`/etc/mise/config.toml`. These can be overridden in your own `mise.toml`.

| Setting | Value | Reason |
|---------|-------|--------|
| `paranoid` | `true` | Enforces HTTPS and stricter security validation |
| `trusted_config_paths` | `["/app"]` | Trusts app config files to avoid warnings during build |
| `idiomatic_version_file_enable_tools` | *(language list)* | Auto-reads version files like `.node-version`, `.python-version`, etc. |
| `install_before` | `"14d"` | Only resolves tool versions released more than 14 days ago, avoiding newly-released versions that may be broken |
| `node.verify` | `false` | Skips asset signature verification for Node, since recently released versions may not yet have a public key |

## Customization

Use a mise configuration file or environment variables to customize mise in the SeaPack-generated container.
Configuration files are generally a better idea.

### Configuration Files

SeaPack automatically detects mise configuration files and passes them
into the build. This includes:

- **Config files**: `mise.toml`, `.mise.toml`, `mise/config.toml`,
  `.mise/config.toml`, `.config/mise.toml`, `.config/mise/config.toml`,
  `.tool-versions`
- **Environment-specific configs**: `mise.*.toml`, `.mise.*.toml`,
  `.config/mise/conf.d/*.toml`
- **Idiomatic version files**: `.ruby-version`, `.python-version`,
  `.python-versions`, `.node-version`, `.nvmrc`, `.go-version`,
  `.java-version`, `.sdkmanrc`, `.bun-version`, `.yvmrc`
- **Lock files**: `mise.lock` files co-located with any detected
  `*.toml` config

### Example: Additional Tools

To add tools that your build needs, add a `mise.toml` to your repository:

```toml
[tools]
node = "22"
jq = "latest"
```
