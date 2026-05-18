---
title: Packages and Version Resolution
description: Understanding how SeaPack resolves package versions using Mise
---

SeaPack providers will analyze the app and determine _fuzzy_ versions of
executables to install. Versions like `3.13`, or `22`. The version resolution
step will resolve those fuzzy versions into the latest valid version that exists.

[Mise](https://mise.jdx.dev/) is used for the package resolution using the `mise
latest package@version` command. Mise is also used for (most) package
installations in the builds as well. However, this is not a requirement of
SeaPack and alternative installation methods are possible when a provider needs
language-specific behavior.
SeaPack enables Mise paranoid mode for stricter security validation.

For more information on how SeaPack utilizes Mise and our philosophy on tool defaults, see the [Mise Configuration](/config/mise) guide.

## Previous and default versions

One important aspect of SeaPack is that updating the default version of
executables in providers (e.g. Node 20 -> 22) should not change the version
installed for apps that have already been building successfully with
SeaPack. This is mainly useful on platforms that use SeaPack to build user
applications (e.g. Sealos).

To support this, you can pass in a `--previous pkg@name ...` flag when
generating the build plan. The typical flow will go like this

- User builds for the first time with SeaPack. The default Node version is used (20).
- The platform saves the resolved versions of packages used
- SeaPack updates the default version of Node to 22
- User submits a new build. The platform passes a `--previous` flag
- SeaPack will use the previous versions instead of using the new default version

This means that SeaPack can freely update default versions of packages without
having to worry about breaking existing apps that rely on the previous defaults.

Passing in a previous version will only be used in place of the default. If a
more specific version of a package is requested (e.g. through a package.json
engines field or env var), then we will always use that.

This is done on Sealos automatically.
