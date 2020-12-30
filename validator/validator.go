// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package validator

import (
	"fmt"
	"os"
	"otoclone/config"
	"otoclone/rclone"
)

// Validate a map of Folders
func Examine(folders map[string]config.Folder) {
    pKeys := map[string]bool{}
    var paths []string

    rKeys := map[string]bool{}
    var remotes []string

    // Extract folders and remotes
    for _, f := range folders {
        if _, value := pKeys[f.Path]; !value {
            pKeys[f.Path] = true
            paths = append(paths, f.Path)
        }
        for _, r := range f.Remotes {
            if _, value := rKeys[r]; !value {
                rKeys[r] = true
                remotes = append(remotes, r)
            }
        }
    }

    validatePaths(paths)
    validateRemotes(remotes)
}

func validateRemotes(remotes []string) {
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

