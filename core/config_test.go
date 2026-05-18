package core

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/gitlayzer/seapack/core/app"
	"github.com/gitlayzer/seapack/core/config"
	"github.com/stretchr/testify/require"
)

func TestGenerateConfigFromEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:    "empty environment",
			envVars: map[string]string{},
			expected: `{
				"steps": {},
				"packages": {},
				"caches": {},
				"deploy": {}
			}`,
		},

		{
			name: "kitchen sink",
			envVars: map[string]string{
				"SEAPACK_INSTALL_CMD":         "npm install",
				"SEAPACK_BUILD_CMD":           "npm run build",
				"SEAPACK_START_CMD":           "npm start",
				"SEAPACK_PACKAGES":            "node@18 python@3.9",
				"SEAPACK_BUILD_APT_PACKAGES":  "build-essential libssl-dev",
				"SEAPACK_DEPLOY_APT_PACKAGES": "libssl-dev",
			},
			expected: `{
				"steps": {
					"install": {
						"name": "install",
						"commands": [
							{ "src": ".", "dest": "." },
							"npm install"
						],
						"secrets": ["*"],
						"assets": {},
						"variables": {}
					},
					"build": {
						"name": "build",
						"commands": [
							{ "src": ".", "dest": "." },
							"npm run build"
						],
						"secrets": ["*"],
						"assets": {},
						"variables": {}
					}
				},
				"buildAptPackages": ["build-essential", "libssl-dev"],
				"packages": {
					"node": "18",
					"python": "3.9"
				},
				"caches": {},
				"deploy": {
					"startCommand": "npm start",
					"aptPackages": ["libssl-dev"]
				},
				"secrets": ["SEAPACK_BUILD_APT_PACKAGES", "SEAPACK_BUILD_CMD", "SEAPACK_DEPLOY_APT_PACKAGES",
					"SEAPACK_INSTALL_CMD", "SEAPACK_PACKAGES", "SEAPACK_START_CMD"]
			}`,
		},

		{
			name: "unversioned packages",
			envVars: map[string]string{
				"SEAPACK_PACKAGES": "jq pipx:httpie@3.2.4",
			},
			expected: `{
				"steps": {},
				"packages": {
					"jq": "latest",
					"pipx:httpie": "3.2.4"
				},
				"caches": {},
				"deploy": {},
				"secrets": ["SEAPACK_PACKAGES"]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := app.NewEnvironment(&tt.envVars)
			gotConfig := GenerateConfigFromEnvironment(env)

			serializedConfig := config.Config{}
			err := json.Unmarshal([]byte(tt.expected), &serializedConfig)
			require.NoError(t, err)

			if diff := cmp.Diff(serializedConfig, *gotConfig); diff != "" {
				t.Errorf("GenerateConfigFromEnvironment() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
