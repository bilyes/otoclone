// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package config

import (
	"fmt"
	"os"
	"otoclone/utils"
	"path/filepath"

	"github.com/spf13/viper"
)

type Remote struct {
	Name   string
	Bucket string
}

type Folder struct {
	Path           string
	Strategy       string
	Remotes        []Remote
	IgnoreList     []string
	ExcludePattern string
}

var configPaths = []string{"$XDG_CONFIG_HOME/otoclone", "$HOME/.config/otoclone"}
var configName = "config"

// Load the configuration from a specific file
func LoadFile(configPath string) (map[string]Folder, error) {
	configName := filepath.Base(configPath)
	configPaths := []string{filepath.Dir(configPath)}

	return loadFrom(configPaths, configName)
}

func Load() (map[string]Folder, error) {
	return loadFrom(configPaths, configName)
}

// Write a Folder to the default config file
func Write(folder Folder) error {
	fs, err := Load()

	if err != nil {
		// Return the error if it's not a viper.ConfigFileNotFoundError
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}

		if err := createConfigFile(); err != nil {
			return err
		}
	}

	return writeFolder(folder, fs)
}

// Write a Folder to a custom config file
func WriteTo(folder Folder, configFile string) error {
	fs, err := LoadFile(configFile)

	if err != nil {
		return err
	}

	return writeFolder(folder, fs)
}

// Remove a folder from the configuration
func Remove(folderName string) error {
	fs, err := Load()

	if err != nil {
		return err
	}

	return removeFolder(folderName, fs)
}

// Remove a folder from a custom config file
func RemoveFrom(folderName string, configFile string) error {
	fs, err := LoadFile(configFile)

	if err != nil {
		return err
	}

	return removeFolder(folderName, fs)
}

func Location() (string, error) {
	p, err := selectConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, configName+".yml"), nil
}

func removeFolder(folderName string, folders map[string]Folder) error {
	delete(folders, folderName)

	viper.Set("folders", folders)

	return viper.WriteConfig()
}

func writeFolder(folder Folder, folders map[string]Folder) error {
	if folders == nil {
		folders = make(map[string]Folder)
	}

	folders[filepath.Base(folder.Path)] = folder

	viper.Set("folders", folders)

	return viper.WriteConfig()
}

func createConfigFile() error {
	configFilePath, err := Location()

	if err != nil {
		return err
	}

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

	return nil
}

func selectConfigPath() (string, error) {
	for _, pa := range configPaths {
		p := os.ExpandEnv(pa)
		ok, err := utils.PathExists(p)

		if err != nil {
			return "", err
		}

		if ok {
			return p, nil
		}
	}

	// TODO Return a proper error type and display instructions
	// explaining how to fix
	return "", fmt.Errorf("No valid path found for the configuration")
}

func setup(configPaths []string, configName string) {
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name

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
		return folders, err
	}

	err = viper.UnmarshalKey("folders", &folders)

	if err != nil {
		return folders, fmt.Errorf("Unable to decode into struct, %v", err)
	}

	return folders, nil
}
