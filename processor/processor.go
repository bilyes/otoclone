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

    for _, f := range folders {
        if strings.HasPrefix(event.Folder, f.Path) {
            subject = f
            break
        }
    }

    if subject.Path == "" {
        return "", []error{&UnwatchedError{event.Folder}}
    }

    if utils.ArrayContains(subject.IgnoreList, event.File) {
        return "", nil
    }

    flags := rclone.Flags{
        Verbose: verbose,
        Exclude: subject.ExcludePattern,
    }

    var errors []error = nil

    for _, r := range subject.Remotes {
        var err error
        switch subject.Strategy {
        case "copy":
            err = p.Cloner.Copy(subject.Path, r.Name, r.Bucket, flags)
        case "sync":
            err = p.Cloner.Sync(subject.Path, r.Name, r.Bucket, flags)
        default:
            err = &UnknownBackupStrategyError{subject.Strategy}
        }
        if err != nil {
            errors = append(errors, err)
        }
    }

    return subject.Path, errors
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
