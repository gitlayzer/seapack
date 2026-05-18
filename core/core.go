package core

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/gitlayzer/seapack/core/app"
	c "github.com/gitlayzer/seapack/core/config"
	"github.com/gitlayzer/seapack/core/generate"
	"github.com/gitlayzer/seapack/core/logger"
	"github.com/gitlayzer/seapack/core/plan"
	"github.com/gitlayzer/seapack/core/providers"
	"github.com/gitlayzer/seapack/core/resolver"
	"github.com/gitlayzer/seapack/internal/utils"
)

const (
	defaultConfigFileName = "seapack.json"
)

type GenerateBuildPlanOptions struct {
	SeaPackVersion           string
	BuildCommand             string
	StartCommand             string
	PreviousVersions         map[string]string
	ConfigFilePath           string
	ErrorMissingStartCommand bool // enabled on sealos
}

type BuildResult struct {
	SeaPackVersion    string                               `json:"seapackVersion,omitempty"`
	Plan              *plan.BuildPlan                      `json:"plan,omitempty"`
	ResolvedPackages  map[string]*resolver.ResolvedPackage `json:"resolvedPackages,omitempty"`
	Metadata          map[string]string                    `json:"metadata,omitempty"`
	DetectedProviders []string                             `json:"detectedProviders,omitempty"`
	Logs              []logger.Msg                         `json:"logs,omitempty"`
	Success           bool                                 `json:"success,omitempty"`
}

func readConfigJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	jsonBytes, err := utils.StandardizeJSON([]byte(data))
	if err != nil {
		return err
	}

	stringData := string(jsonBytes)

	if err := json.Unmarshal([]byte(stringData), v); err != nil {
		return fmt.Errorf("error reading %s as JSON: %w", path, err)
	}

	return nil
}

func GenerateBuildPlan(app *app.App, env *app.Environment, options *GenerateBuildPlanOptions) *BuildResult {
	logger := logger.NewLogger()

	config, err := GetConfig(app, env, options, logger)
	if err != nil {
		logger.LogError("%s", err.Error())
		return &BuildResult{Success: false, Logs: logger.Logs}
	}

	ctx, err := generate.NewGenerateContext(app, env, config, logger)
	if err != nil {
		logger.LogError("%s", err.Error())
		return &BuildResult{Success: false, Logs: logger.Logs}
	}

	// Set the previous versions
	if options.PreviousVersions != nil {
		for name, version := range options.PreviousVersions {
			ctx.Resolver.SetPreviousVersion(name, version)
		}
	}

	// Figure out what providers to use
	providerToUse, detectedProviderName := getProviders(ctx, config)
	ctx.Metadata.Set("providers", detectedProviderName)

	// TODO: We should indicate if we have packages specified in the config
	// so that providers can determine if they should include mise in the final image (e.g. for shell script)

	if providerToUse != nil {
		err = providerToUse.Plan(ctx)
		if err != nil {
			logger.LogError("%s", err.Error())
			return &BuildResult{Success: false, Logs: logger.Logs}
		}
	}

	// before `Generate()` any commands provided by seapack.json are *not* merged into the provider-generated
	// buildPlan. This means providers can't view any of the custom structure provided by the user via a seapack.json
	buildPlan, resolvedPackages, err := ctx.Generate()
	if err != nil {
		logger.LogError("%s", err.Error())
		return &BuildResult{Success: false, Logs: logger.Logs}
	}

	if providerToUse != nil {
		providerToUse.CleansePlan(buildPlan)
	}

	if !ValidatePlan(buildPlan, app, logger, &ValidatePlanOptions{
		ErrorMissingStartCommand: options.ErrorMissingStartCommand,
		ProviderToUse:            providerToUse,
	}) {
		return &BuildResult{Success: false, Logs: logger.Logs}
	}

	buildResult := &BuildResult{
		SeaPackVersion:    options.SeaPackVersion,
		Plan:              buildPlan,
		ResolvedPackages:  resolvedPackages,
		Metadata:          ctx.Metadata.Properties,
		DetectedProviders: []string{detectedProviderName},
		Logs:              logger.Logs,
		Success:           true,
	}

	return buildResult
}

// GetConfig merges the options, environment, and file config into a single config
func GetConfig(app *app.App, env *app.Environment, options *GenerateBuildPlanOptions, logger *logger.Logger) (*c.Config, error) {
	optionsConfig := GenerateConfigFromOptions(options)

	envConfig := GenerateConfigFromEnvironment(env)

	fileConfig, err := GenerateConfigFromFile(app, env, options, logger)
	if err != nil {
		return nil, err
	}

	mergedConfig := c.Merge(optionsConfig, envConfig, fileConfig)

	return mergedConfig, nil
}

func GenerateConfigFromFile(app *app.App, env *app.Environment, options *GenerateBuildPlanOptions, logger *logger.Logger) (*c.Config, error) {
	config := c.EmptyConfig()

	configFileName := defaultConfigFileName
	if options.ConfigFilePath != "" {
		configFileName = options.ConfigFilePath
	}

	if envConfigFileName, _ := env.GetConfigVariable("CONFIG_FILE"); envConfigFileName != "" {
		configFileName = envConfigFileName
	}

	// always assume config file path is relative to the app source directory
	// https://github.com/gitlayzer/seapack/pull/226
	absConfigFileName := filepath.Join(app.Source, configFileName)

	if _, err := os.Stat(absConfigFileName); err != nil && os.IsNotExist(err) {
		// if a specific path was specified, we should indicate that it was not found and hard fail
		if configFileName != defaultConfigFileName {
			return nil, fmt.Errorf("config file %q not found", absConfigFileName)
		}

		return config, nil
	}

	// if a JSON file was provided, we should hard fail if we cannot parse it
	if err := readConfigJSON(absConfigFileName, config); err != nil {
		logger.LogWarn("Failed to read config file `%s`\nUse the following schema to validate your config file: %s\n", configFileName, c.SchemaUrl)
		return nil, err
	}

	logger.LogInfo("Using config file `%s`", configFileName)
	logger.LogWarn("The config file format is not yet finalized and subject to change.")

	return config, nil
}

func GenerateConfigFromEnvironment(env *app.Environment) *c.Config {
	config := c.EmptyConfig()

	if env == nil {
		return config
	}

	if installCmdVar, _ := env.GetConfigVariable("INSTALL_CMD"); installCmdVar != "" {
		installStep := config.GetOrCreateStep("install")
		installStep.Commands = []plan.Command{
			plan.NewCopyCommand("."),
			plan.NewExecShellCommand(installCmdVar, plan.ExecOptions{CustomName: installCmdVar}),
		}
	}

	if buildCmdVar, _ := env.GetConfigVariable("BUILD_CMD"); buildCmdVar != "" {
		buildStep := config.GetOrCreateStep("build")
		buildStep.Commands = []plan.Command{
			plan.NewCopyCommand("."),
			plan.NewExecShellCommand(buildCmdVar, plan.ExecOptions{CustomName: buildCmdVar}),
		}
	}

	if startCmdVar, _ := env.GetConfigVariable("START_CMD"); startCmdVar != "" {
		config.Deploy.StartCmd = startCmdVar
	}

	if packages, _ := env.GetConfigVariableList("PACKAGES"); len(packages) > 0 {
		config.Packages = utils.ParsePackageWithVersion(packages)
	}

	if aptPackages, _ := env.GetConfigVariableList("BUILD_APT_PACKAGES"); len(aptPackages) > 0 {
		config.BuildAptPackages = aptPackages
	}

	if aptPackages, _ := env.GetConfigVariableList("DEPLOY_APT_PACKAGES"); len(aptPackages) > 0 {
		config.Deploy.AptPackages = aptPackages
	}

	config.Secrets = append(config.Secrets, slices.Sorted(maps.Keys(env.Variables))...)

	return config
}

// generates a config from the CLI options
func GenerateConfigFromOptions(options *GenerateBuildPlanOptions) *c.Config {
	config := c.EmptyConfig()

	if options == nil {
		return config
	}

	if options.BuildCommand != "" {
		buildStep := config.GetOrCreateStep("build")
		buildStep.Commands = []plan.Command{
			plan.NewCopyCommand("."),
			plan.NewExecShellCommand(options.BuildCommand, plan.ExecOptions{CustomName: options.BuildCommand}),
		}
	}

	if options.StartCommand != "" {
		config.Deploy.StartCmd = options.StartCommand
	}

	return config
}

func getProviders(ctx *generate.GenerateContext, config *c.Config) (providers.Provider, string) {
	allProviders := providers.GetLanguageProviders()

	var providerToUse providers.Provider
	var detectedProvider string

	// Even if there are providers manually specified, we want to detect to see what type of app this is
	for _, provider := range allProviders {
		matched, err := provider.Detect(ctx)
		if err != nil {
			log.Warnf("Failed to detect provider `%s`: %s", provider.Name(), err.Error())
			continue
		}

		if matched {
			detectedProvider = provider.Name()

			// If there are no providers manually specified in the config,
			if config.Provider == nil {
				if err := provider.Initialize(ctx); err != nil {
					ctx.Logger.LogWarn("Failed to initialize provider `%s`: %s", provider.Name(), err.Error())
					continue
				}

				ctx.Logger.LogInfo("Detected %s", utils.CapitalizeFirst(provider.Name()))

				providerToUse = provider
			}

			break
		}
	}

	if config.Provider != nil {
		provider := providers.GetProvider(*config.Provider)

		if provider == nil {
			ctx.Logger.LogWarn("Provider `%s` not found", *config.Provider)
			return providerToUse, detectedProvider
		}

		if err := provider.Initialize(ctx); err != nil {
			ctx.Logger.LogWarn("Failed to initialize provider `%s`: %s", *config.Provider, err.Error())
			return providerToUse, detectedProvider
		}

		ctx.Logger.LogInfo("Using provider %s from config", utils.CapitalizeFirst(*config.Provider))
		providerToUse = provider
	}

	return providerToUse, detectedProvider
}
