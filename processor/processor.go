// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package processor

import (
	"fmt"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/rclone"
)

// Handles a FSNotify Event
func Handle(event fsnotify.FSEvent, folders map[string]config.Folder, verbose bool) (string, []error) {
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

    if contains(subject.IgnoreList, event.File) {
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
            err = rclone.Copy(subject.Path, r.Name, r.Bucket, flags)
        case "sync":
            err = rclone.Sync(subject.Path, r.Name, r.Bucket, flags)
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

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}

