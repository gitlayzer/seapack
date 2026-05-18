package testing

import (
	"testing"

	"github.com/gitlayzer/seapack/core/app"
	"github.com/gitlayzer/seapack/core/config"
	"github.com/gitlayzer/seapack/core/generate"
	"github.com/gitlayzer/seapack/core/logger"
)

// CreateGenerateContext creates a new GenerateContext for testing purposes
func CreateGenerateContext(t *testing.T, path string) *generate.GenerateContext {
	t.Helper() // This marks the function as a test helper, which improves test output

	userApp, err := app.NewApp(path)
	if err != nil {
		t.Fatalf("error creating app: %v", err)
	}

	env := app.NewEnvironment(nil)

	config := config.EmptyConfig()

	ctx, err := generate.NewGenerateContext(userApp, env, config, logger.NewLogger())
	if err != nil {
		t.Fatalf("error creating generate context: %v", err)
	}

	return ctx
}
