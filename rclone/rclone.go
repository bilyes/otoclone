// Author: ilyess bachiri
// Copyright (c) 2020-present ilyess bachiri

package rclone

import (
	"os"
	"os/exec"
	"strings"
)

var cmd string = "rclone"
var remotes []string

// Copies the content of a folder into a remote bucket
func Copy(folder string, remote string, bucket string) error {

    copy := exec.Command(cmd, "copy", "-v", folder, remote + ":" + bucket)

    copy.Stdout = os.Stdout
    copy.Stderr = os.Stderr

    if err := copy.Run(); err != nil {
        return err
    }

    return nil
}

// Checks if a given remote has been configured
func RemoteIsValid(remote string) (bool, error) {
    if len(remotes) == 0 {
        listRemotes := exec.Command(cmd, "listremotes")
        stdout, err := listRemotes.Output()
        if err != nil {
            return false, err
        }

        remotes = strings.Split(string(stdout), ":\n")
    }


    if !contains(remotes, remote) {
        return false, nil
    }

    return true, nil
}

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}
