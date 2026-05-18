package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gitlayzer/seapack/core/app"
	"github.com/gitlayzer/seapack/core/logger"
	"github.com/stretchr/testify/require"
)

// generate plans for the supported SeaPack MVP providers
func TestGenerateBuildPlanForExamples(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	examplesDir := filepath.Join(filepath.Dir(wd), "examples")
	tests := map[string]string{
		"node-npm":    "node",
		"python-pip":  "python",
		"go-mod":      "golang",
		"java-gradle": "java",
		"deno-2":      "deno",
	}

	for example, expectedProvider := range tests {
		t.Run(example, func(t *testing.T) {
			examplePath := filepath.Join(examplesDir, example)

			userApp, err := app.NewApp(examplePath)
			require.NoError(t, err)

			env := app.NewEnvironment(nil)
			buildResult := GenerateBuildPlan(userApp, env, &GenerateBuildPlanOptions{})

			require.True(t, buildResult.Success, "logs: %v", buildResult.Logs)
			require.Equal(t, []string{expectedProvider}, buildResult.DetectedProviders)

			usesSeaPackBase := strings.HasPrefix(buildResult.Plan.Deploy.Base.Image, "ghcr.io/gitlayzer/seapack-runtime:")
			for _, step := range buildResult.Plan.Steps {
				for _, input := range step.Inputs {
					if input.Image == "" {
						continue
					}
					if strings.HasPrefix(input.Image, "ghcr.io/gitlayzer/seapack-") {
						usesSeaPackBase = true
						isBuilder := strings.HasPrefix(input.Image, "ghcr.io/gitlayzer/seapack-builder:")
						isRuntime := strings.HasPrefix(input.Image, "ghcr.io/gitlayzer/seapack-runtime:")
						require.True(t, isBuilder || isRuntime, "unexpected SeaPack image %q", input.Image)
					}
				}
			}
			require.True(t, usesSeaPackBase)
		})
	}
}

func TestGenerateConfigFromFile_NotFound(t *testing.T) {
	// Use an existing example app directory so relative paths resolve
	appPath := "../examples/config-file"
	userApp, err := app.NewApp(appPath)
	require.NoError(t, err)

	env := app.NewEnvironment(nil)
	l := logger.NewLogger()

	options := &GenerateBuildPlanOptions{ConfigFilePath: "does-not-exist.seapack.json"}
	cfg, genErr := GenerateConfigFromFile(userApp, env, options, l)

	require.Error(t, genErr, "expected an error when explicit config file does not exist")
	require.Nil(t, cfg, "config should be nil on error")
}

func TestGenerateConfigFromFile_Malformed(t *testing.T) {
	appPath := "../examples/config-file"
	userApp, err := app.NewApp(appPath)
	require.NoError(t, err)

	env := app.NewEnvironment(nil)
	l := logger.NewLogger()

	options := &GenerateBuildPlanOptions{ConfigFilePath: "seapack.malformed.json"}
	cfg, genErr := GenerateConfigFromFile(userApp, env, options, l)

	require.Error(t, genErr, "expected an error for malformed JSON config file")
	require.Nil(t, cfg, "config should be nil on error")
}

func TestGenerateBuildPlan_DockerignoreMetadata(t *testing.T) {
	sourcePath := "../examples/node-npm"
	appPath := filepath.Join("..", "tmp", "core-test-node-npm-dockerignore")
	require.NoError(t, copyDir(sourcePath, appPath))
	require.NoError(t, os.WriteFile(filepath.Join(appPath, ".dockerignore"), []byte("node_modules\n"), 0644))

	userApp, err := app.NewApp(appPath)
	require.NoError(t, err)

	env := app.NewEnvironment(nil)
	buildResult := GenerateBuildPlan(userApp, env, &GenerateBuildPlanOptions{})

	require.True(t, buildResult.Success)
	require.NotNil(t, buildResult.Metadata)
	require.Equal(t, "true", buildResult.Metadata["dockerIgnore"])
}

func copyDir(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}
