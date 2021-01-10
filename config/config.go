package config

import (
	"fmt"
	"os"
	"otoclone/utils"
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
    ExcludePattern string
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

// Write a Folder to the config file
func Write(folder Folder) error {
    configPath := "/tmp/play"
    configName := "test"

    configFilePath := filepath.Join(configPath, configName + ".yml")

    ok, err := utils.PathExists(configFilePath)
    if err != nil {
        return err
    }
    if !ok {
        cnf, err := os.Create(configFilePath)
        if err != nil {
            return err
        }
        defer cnf.Close()
    }

    setup([]string{configPath},configName)

    fs := make(map[string]Folder)
    fs[filepath.Base(folder.Path)] = folder

    viper.Set("folders", fs)

    return viper.WriteConfig()
}

func setup(configPaths []string, configName string) {
    viper.SetConfigName(configName) // name of config file (without extension)
    viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

    for _, cp := range configPaths {
        viper.AddConfigPath(cp)
    }
}

// Load the configuration from the supported config locations
func loadFrom(configPaths []string, configName string) (map[string]Folder, error) {
    setup(configPaths, configName)

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
