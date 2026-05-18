package cpp

import "github.com/gitlayzer/seapack/core/generate"

type buildSystem interface {
	Install(ctx *generate.GenerateContext, pkgs *generate.MiseStepBuilder)
	Build(build *generate.CommandStepBuilder)
}
