// Author: ilyess bachiri
// Copyright (c) 2020-present ilyess bachiri

package fsnotify

import (
	"os/exec"
	"strings"
)


type FSEvent struct {
    Folder string
    Event string
    File string
}

// Sets a filesystem watcher on a given set of folders
func Watch(folders []string) (FSEvent, error) {
    foldersToWatch := strings.Join(folders, " ")
    watcher := exec.Command("inotifywait", "-r", "-e", "modify,create,delete,move", foldersToWatch)
    stdout, err := watcher.Output()

    if err != nil {
        return FSEvent{}, err
    }

    return parseEvent(stdout), nil
}

func parseEvent(output []byte) FSEvent {
    event := strings.Split(string(output), " ")
    event[2] = strings.TrimSuffix(event[2], "\n")

    return FSEvent{
        Folder: event[0],
        Event: event[1],
        File: event[2],
    }
}
