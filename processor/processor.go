// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/rclone"
)

// Handles a FSNotify Event
func Handle(event fsnotify.FSEvent, folders map[string]config.Folder) (string, error) {
    var subject config.Folder

    for _, f := range folders {
        if strings.HasPrefix(event.Folder, f.Path) {
            subject = f
            break
        }
    }

    if subject.Path == "" {
        return "", &UnwatchedError{event.Folder}
    }

    if contains(subject.IgnoreList, event.File) {
        return "", nil
    }

    for _, r := range subject.Remotes {
        err := rclone.Copy(subject.Path, r, filepath.Base(subject.Path))
        if err != nil {
            return "", err
        }
    }

    return subject.Path, nil
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

