package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type FolderConfig struct {
    Path string
    Strategy string
    Remotes []string
    IgnoreList []string
}

var configPaths = []string{"$XDG_CONFIG_HOME/otoclone", "$HOME/.config/otoclone"}

// Load the configuration from a specific file
func LoadFile(configPath string) (map[string]FolderConfig, error) {

    configName := filepath.Base(configPath)
    configPaths := []string{filepath.Dir(configPath)}

    return loadFrom(configPaths, configName)
}

func Load() (map[string]FolderConfig, error) {

    return loadFrom(configPaths, "config")
}

// Load the configuration from the supported config locations
func loadFrom(configPaths []string, configName string) (map[string]FolderConfig, error) {
    viper.SetConfigName(configName) // name of config file (without extension)
    viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
    for _, cp := range configPaths {
        viper.AddConfigPath(cp)
    }

    err := viper.ReadInConfig() // Find and read the config file
    if err != nil {
        // TODO return error instead
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    var foldersConfig map[string]FolderConfig

    err = viper.UnmarshalKey("folders", &foldersConfig)

    if err != nil {
        // TODO return error instead
        panic(fmt.Errorf("Unable to decode into struct, %v", err))

    }

    return foldersConfig, nil
}
