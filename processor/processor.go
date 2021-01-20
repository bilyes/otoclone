// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package processor

import (
	"fmt"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/rclone"
	"otoclone/utils"
)

type Processor struct {
    Cloner rclone.Cloner
}

// Handles a FSNotify Event
func (p *Processor) Handle(event fsnotify.FSEvent, folders map[string]config.Folder, verbose bool) (string, []error) {
    var subject config.Folder
    var fKey string

    for k, f := range folders {
        if strings.HasPrefix(event.Folder, f.Path) {
            subject = f
            fKey = k
            break
        }
    }

    if subject.Path == "" {
        return "", []error{&UnwatchedError{event.Folder}}
    }

    if utils.ArrayContains(subject.IgnoreList, event.File) {
        return "", nil
    }

    flds := make(map[string]config.Folder)
    flds[fKey] = subject

    return subject.Path, p.Backup(flds, verbose)
}

// Backup a list of folders
func (p *Processor) Backup(folders map[string]config.Folder, verbose bool) []error {
    var errors []error = nil

    for _, folder := range folders {

        flags := rclone.Flags{
            Verbose: verbose,
            Exclude: folder.ExcludePattern,
        }

        for _, r := range folder.Remotes {
            var err error
            switch folder.Strategy {
            case "copy":
                err = p.Cloner.Copy(folder.Path, r.Name, r.Bucket, flags)
            case "sync":
                err = p.Cloner.Sync(folder.Path, r.Name, r.Bucket, flags)
            default:
                err = &UnknownBackupStrategyError{folder.Strategy}
            }
            if err != nil {
                errors = append(errors, err)
            }
        }
    }

    return errors
}

type UnknownBackupStrategyError struct {
    Strategy string
}

func (e *UnknownBackupStrategyError) Error() string {
    return fmt.Sprintf("Unsupported backup strategy %s", e.Strategy)
}

type UnwatchedError struct {
    Subject string
}

func (e *UnwatchedError) Error() string {
    return fmt.Sprintf("Unwatched file or directory %s", e.Subject)
}
