---
title: Environment Variables
description: Understanding environment variables in SeaPack
---

Some parts of the build can be configured with environment variables. These are
often prefixed with `SEAPACK_`.

## Build Configuration

| Name                           | Description                                                                                                                                                                     |
| :----------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `SEAPACK_BUILD_CMD`           | Set the command to run for the build step. This overwrites any commands that come from providers                                                                                |
| `SEAPACK_INSTALL_CMD`         | Set the command to run for the install step. This overwrites any commands that come from providers. All files are copied to the root of the project before running the command. |
| `SEAPACK_START_CMD`           | Set the command to run when the container starts                                                                                                                                |
| `SEAPACK_PACKAGES`            | Install additional Mise packages. In the format `pkg[@version]`. The version is optional; if not provided, the latest version is used. Allows list.                             |
| `SEAPACK_BUILD_APT_PACKAGES`  | Install additional Apt packages during build. Allows list.                                                                                                                      |
| `SEAPACK_DEPLOY_APT_PACKAGES` | Install additional Apt packages in the final image. Allows list.                                                                                                                |
| `SEAPACK_DISABLE_CACHES`      | Specify specific BuildKit cache keys to disable, or `*` to disable all caches. Allows list.                                                                                     |

Variables which allow a list use space-separated values. For example:

```sh
SEAPACK_PACKAGES="pipx:httpie jq@latest"
```

To configure more parts of the build, it is recommended to use a [config file](/config/file).

## Global Options

These environment variables affect the behavior of SeaPack:

| Name              | Description                                 |
| :---------------- | :------------------------------------------ |
| `FORCE_COLOR`     | Force colored output even when not in a TTY |
| `SEAPACK_VERBOSE` | Enable verbose logging (equivalent to `--verbose` flag) |
