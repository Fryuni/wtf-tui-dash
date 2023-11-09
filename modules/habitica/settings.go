package habitica

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
	"os"
)

const (
	defaultFocusable = true
	defaultTitle     = "Habitica"
)

// Settings defines the configuration properties for this module
type Settings struct {
	common *cfg.Common

	userId   string `help:"User ID for Habitica API."`
	apiToken string `help:"Your Habitica API token."`
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		apiToken: ymlConfig.UString("apiToken", os.Getenv("WTF_HABITICA_TOKEN")),
		userId:   ymlConfig.UString("userId"),
	}

	cfg.ModuleSecret(name, globalConfig, &settings.apiToken).Load()

	return &settings
}
