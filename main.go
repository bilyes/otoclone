// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/processor"
	"otoclone/validator"
)

func main() {
    configFile, verbose := parseFlags()

    folders := loadConfig(configFile)
    validator.Examine(folders)

    paths := extractPaths(folders)
    foldersToWatch := strings.Join(paths, " ")

    for {
        fmt.Println("Watching", foldersToWatch)
        event, err := fsnotify.Watch(paths)

        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        path, err := processor.Handle(event, folders, verbose)

        if path == "" {
            fmt.Println("Ignored", event.File)
        } else {
            fmt.Println("Backed up", path)
        }
    }
}

func parseFlags() (string, bool) {
    var configFile string
    var verbose bool
    flag.StringVar(&configFile, "f", "", "Path to the configuration file")
    flag.BoolVar(&verbose, "v", false, "Increase verbosity")
    flag.Parse()

    return configFile, verbose
}

func loadConfig(configFile string) map[string]config.Folder {

    var folders map[string]config.Folder
    var err error

    if configFile != "" {
        folders, err = config.LoadFile(configFile)
    } else {
        folders, err = config.Load()
    }

    if err != nil {
        fmt.Println("Error:",  err)
        os.Exit(1)
    }
    return folders
}

func extractPaths(folders map[string]config.Folder) []string {
    fKeys := map[string]bool{}
    var paths []string

    for _, f := range folders {
        if _, value := fKeys[f.Path]; !value {
            fKeys[f.Path] = true
            paths = append(paths, f.Path)
        }
    }
    return paths
}

