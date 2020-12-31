package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Remote struct {
    Name string
    Bucket string
}

type Folder struct {
    Path string
    Strategy string
    Remotes []Remote
    IgnoreList []string
}

var configPaths = []string{"$XDG_CONFIG_HOME/otoclone", "$HOME/.config/otoclone"}

// Load the configuration from a specific file
func LoadFile(configPath string) (map[string]Folder, error) {

    configName := filepath.Base(configPath)
    configPaths := []string{filepath.Dir(configPath)}

    return loadFrom(configPaths, configName)
}

func Load() (map[string]Folder, error) {

    return loadFrom(configPaths, "config")
}

// Load the configuration from the supported config locations
func loadFrom(configPaths []string, configName string) (map[string]Folder, error) {

    viper.SetConfigName(configName) // name of config file (without extension)
    viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

    for _, cp := range configPaths {
        viper.AddConfigPath(cp)
    }

    var folders map[string]Folder

    err := viper.ReadInConfig()
    if err != nil {
        return folders, fmt.Errorf("Fatal error config file: %s \n", err)
    }

    err = viper.UnmarshalKey("folders", &folders)

    if err != nil {
        return folders, fmt.Errorf("Unable to decode into struct, %v", err)
    }

    return folders, nil
}
