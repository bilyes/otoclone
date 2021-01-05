// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package validator

import (
	"fmt"
	"os"

	"otoclone/config"
	"otoclone/processor"
	"otoclone/rclone"
	"otoclone/utils"
)

var strategies = []string{"copy", "sync"}

// Validate a map of Folders.
// Check if the paths exist on the filesystem, the backup
// strategies are supported and the remotes are all configured.
func Examine(folders map[string]config.Folder) []error {
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

    var errors []error

    if err := validateStrategies(strats); err != nil {
        errors = append(errors, err)
    }
    if err := validatePaths(paths); err != nil {
        errors = append(errors, err)
    }
    if err:= validateRemotes(remotes); err != nil {
        errors = append(errors, err)
    }

    return errors
}

type UnknownRemote struct {
    Remote string
}

func (e *UnknownRemote) Error() string {
    return fmt.Sprintf("Unknown remote %s", e.Remote)
}

type NoSuchDirectoryError struct {
    Path string
}

func (e *NoSuchDirectoryError) Error() string {
    return fmt.Sprintf("No such directory %s", e.Path)
}

func validateStrategies(strats []string) error {
    for _, s := range strats {
        if !utils.ArrayContains(strategies, s) {
            return &processor.UnknownBackupStrategyError{Strategy: s}
        }
    }
    return nil
}

func validateRemotes(remotes []config.Remote) error {
    for _, r := range remotes {
        isValid, err := rclone.RemoteIsValid(r.Name)
        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        if !isValid {
            return &UnknownRemote{r.Name}
        }
    }
     return nil
}

func validatePaths(paths []string) error {
    for _, p := range paths {
        if ok, err := utils.PathExists(p); !ok {
            if err == nil {
                return &NoSuchDirectoryError{p}
            } else {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
        }
    }
    return nil
}
