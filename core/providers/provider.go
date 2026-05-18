package providers

import (
	"github.com/gitlayzer/seapack/core/generate"
	"github.com/gitlayzer/seapack/core/plan"
	"github.com/gitlayzer/seapack/core/providers/deno"
	"github.com/gitlayzer/seapack/core/providers/golang"
	"github.com/gitlayzer/seapack/core/providers/java"
	"github.com/gitlayzer/seapack/core/providers/node"
	"github.com/gitlayzer/seapack/core/providers/python"
)

type Provider interface {
	Name() string
	Detect(ctx *generate.GenerateContext) (bool, error)
	Initialize(ctx *generate.GenerateContext) error
	Plan(ctx *generate.GenerateContext) error
	CleansePlan(buildPlan *plan.BuildPlan)
	StartCommandHelp() string
}

func GetLanguageProviders() []Provider {
	// Order is important here. The first provider that returns true from Detect() will be used.
	return []Provider{
		&node.NodeProvider{},
		&python.PythonProvider{},
		&golang.GoProvider{},
		&java.JavaProvider{},
		&deno.DenoProvider{},
	}
}

func GetProvider(name string) Provider {
	for _, provider := range GetLanguageProviders() {
		if provider.Name() == name {
			return provider
		}
	}

	return nil
}
