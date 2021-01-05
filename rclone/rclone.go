// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package rclone

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Flags struct {
    Verbose bool
    Exclude string
}

var cmd string = "rclone"
var remotes []string

// Copies the content of a folder into a remote bucket
func Copy(folder string, remote string, bucket string, flags Flags) error {
    return transfer("copy", folder, remote, bucket, flags)
}

// Syncs the content of a source folder and a remote bucket
func Sync(folder string, remote string, bucket string, flags Flags) error {
    return transfer("sync", folder, remote, bucket, flags)
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

func transfer(strategy string, folder string, remote string, bucket string, flags Flags) error {
    args := []string{}
    if flags.Verbose {
        args = append(args, "-v")
    }
    if len(flags.Exclude) > 0 {
        args = append(args, fmt.Sprintf("--exclude=%s", flags.Exclude))
    }
    args = append(args, strategy, folder, fmt.Sprintf("%s:%s", remote, bucket))

    copy := exec.Command(cmd, args...)

    copy.Stdout = os.Stdout
    copy.Stderr = os.Stderr

    if err := copy.Run(); err != nil {
        return err
    }
    return nil
}

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}
