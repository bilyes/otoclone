// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/rclone"
)

func main() {

    var configFile string
    flag.StringVar(&configFile, "f", "", "Path to the configuration file")
    flag.Parse()

    var foldersConfig map[string]config.FolderConfig
    var err error

    if configFile != "" {
        foldersConfig, err = config.LoadFile(configFile)
    } else {
        foldersConfig, err = config.Load()
    }

    if err != nil {
        fmt.Println("Error:",  err)
        os.Exit(1)
    }

    validate(foldersConfig)

    folders := extractFolders(foldersConfig)
    foldersToWatch := strings.Join(folders, " ")


    for {
        fmt.Println("Watching", foldersToWatch)
        event, err := fsnotify.Watch(folders)

        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        process(event, foldersConfig)
    }
}

func extractFolders(foldersConfig map[string]config.FolderConfig) []string {
    fKeys := map[string]bool{}
    var folders []string

    for _, f := range foldersConfig {
        if _, value := fKeys[f.Path]; !value {
            fKeys[f.Path] = true
            folders = append(folders, f.Path)
        }
    }

    return folders
}

func validate(foldersConfig map[string]config.FolderConfig) {
    fKeys := map[string]bool{}
    var folders []string

    rKeys := map[string]bool{}
    var remotes []string

    // Extract folders and remotes
    for _, f := range foldersConfig {
        if _, value := fKeys[f.Path]; !value {
            fKeys[f.Path] = true
            folders = append(folders, f.Path)
        }
        for _, r := range f.Remotes {
            if _, value := rKeys[r]; !value {
                rKeys[r] = true
                remotes = append(remotes, r)
            }
        }
    }

    // Validate folders
    for _, f := range folders {
        if result, err := exists(f); !result {
            if err == nil {
                fmt.Println("Error: No such directory", f)
            } else {
                fmt.Println("Error:", err)

            }
            os.Exit(1)
        }
    }

    // Validate remotes
    for _, r := range remotes {
        isValid, err := rclone.RemoteIsValid(r)
        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        if !isValid {
            fmt.Print("Error: Unknown remote ", r)
            fmt.Println(". To configure it run: rclone config")
            os.Exit(1)
        }
    }
}

func process(event fsnotify.FSEvent, folders map[string]config.FolderConfig) {
    var subject config.FolderConfig

    for _, f := range folders {
        if strings.HasPrefix(event.Folder, f.Path) {
            subject = f
            break
        }
    }

    if subject.Path == "" {
        fmt.Println("Error: Unwatched file or folder", event.Folder)
        return
    }

    if contains(subject.IgnoreList, event.File) {
        fmt.Println("Ignoring", event.File)
        return
    }

    for _, r := range subject.Remotes {
        err := rclone.Copy(subject.Path, r, filepath.Base(subject.Path))
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("Backed up", subject.Path)
}

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}
