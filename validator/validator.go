// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package validator

import (
	"fmt"
	"os"
	"otoclone/config"
	"otoclone/rclone"
	"strings"
)

var strategies = []string{"copy", "sync"}

// Validate a map of Folders.
// Check if the paths exist on the filesystem, the backup
// strategies are supported and the remotes are all configured.
func Examine(folders map[string]config.Folder) {
    pKeys := map[string]bool{}
    var paths []string

    rKeys := map[string]bool{}
    var remotes []config.Remote

    sKeys := map[string]bool{}
    var strats []string

    // Extract folders and remotes
    for _, f := range folders {
        if _, value := sKeys[f.Strategy]; !value {
            pKeys[f.Strategy] = true
            strats = append(strats, f.Strategy)
        }

        if _, value := pKeys[f.Path]; !value {
            pKeys[f.Path] = true
            paths = append(paths, f.Path)
        }

        for _, r := range f.Remotes {
            if _, value := rKeys[r.Name]; !value {
                rKeys[r.Name] = true
                remotes = append(remotes, r)
            }
        }
    }
    validateStrategies(strats)
    validatePaths(paths)
    validateRemotes(remotes)
}

func validateStrategies(strats []string) {
    for _, s := range strats {
        if !contains(strategies, s) {
            fmt.Println("Error: Unknown backup strategy", s)
            fmt.Println("Supported strategies are:", strings.Join(strategies, ", "))
            os.Exit(1)
        }
    }
}

func validateRemotes(remotes []config.Remote) {
    for _, r := range remotes {
        isValid, err := rclone.RemoteIsValid(r.Name)
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

func validatePaths(paths []string) {
    for _, p := range paths {
        if ok, err := exists(p); !ok {
            if err == nil {
                fmt.Println("Error: No such directory", p)
            } else {
                fmt.Println("Error:", err)
            }
            os.Exit(1)
        }
    }
}

// Check if a path exists on the filesystem
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}
