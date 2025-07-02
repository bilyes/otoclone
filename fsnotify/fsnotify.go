// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package fsnotify

import (
	"fmt"
	"os/exec"
	"strings"
)

type FSEvent struct {
	Folder string
	Event  string
	File   string
}

// Sets a filesystem watcher on a given set of folders
func Watch(folders []string) (FSEvent, error) {
	args := []string{
		"-r",
		"-e",
		"modify,create,delete,move",
	}
	args = append(args, folders...)

	watcher := exec.Command("inotifywait", args...)
	stdout, err := watcher.Output()

	if err != nil {
		return FSEvent{}, err
	}

	event, err := parseEvent(stdout)
	if err != nil {
		return FSEvent{}, err
	}

	return event, nil
}

func parseEvent(output []byte) (FSEvent, error) {
	event := strings.Split(string(output), " ")
	if len(event) < 3 {
		return FSEvent{}, fmt.Errorf("Unable to parse event: %s", string(output))
	}
	event[2] = strings.TrimSuffix(event[2], "\n")

	return FSEvent{
		Folder: event[0],
		Event:  event[1],
		File:   event[2],
	}, nil
}
