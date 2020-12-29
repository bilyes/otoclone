// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

    if len(os.Args) < 3 {
        fmt.Println("Missing argument. Usage: otoclone FOLDERS REMOTES")
        os.Exit(1)
    }

    folders := strings.Split(os.Args[1], ",")
    validateFolders(folders)

    remotes := strings.Split(os.Args[2], ",")
    validateRemotes(remotes)

    foldersToWatch := strings.Join(folders, " ")

    // TODO read from config file
    ignoreList := []string{"4913"}

    for {
        fmt.Println("Watching", foldersToWatch)
        watcher := exec.Command("inotifywait", "-r", "-e", "modify,create,delete,move", foldersToWatch)
        stdout, err := watcher.Output()

        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        process(parseINotifyEvent(stdout), ignoreList, folders, remotes)
    }
}

func parseINotifyEvent(output []byte) []string {
    event := strings.Split(string(output), " ")
    event[2] = strings.TrimSuffix(event[2], "\n")

    return event
}

func process(event []string, ignoreList []string, folders []string, remotes []string) {
    if contains(ignoreList, event[2]) { return }

    folder := ""

    for _, f := range folders {
        if strings.HasPrefix(event[0], f) {
            folder = f
            break
        }
    }

    if folder == "" {
        fmt.Println("Error: Unwatched file or folder", event[0])
        return
    }

    for _, r := range remotes {
        remoteBucket := r + ":" + filepath.Base(folder)

        //fmt.Println("running rclone copy -v", folder, remoteBucket)
        copy := exec.Command("rclone", "copy", "-v", folder, remoteBucket)

        copy.Stdout = os.Stdout
        copy.Stderr = os.Stderr

        if err := copy.Run(); err != nil {
            log.Panic(err)
        }
    }

    fmt.Println("Backed up", folder)
}

func validateRemotes(remotes []string) {
    listRemotes := exec.Command("rclone", "listremotes")
    stdout, err := listRemotes.Output()

    if err != nil {
        fmt.Println("Error:",  err)
        os.Exit(1)
    }

    configuredRemotes := strings.Split(string(stdout), ":\n")

    for _, remote := range remotes {
        if !contains(configuredRemotes, remote) {
            fmt.Print("Error: Unknown remote ", remote)
            fmt.Println(". To configure it run: rclone config")
            os.Exit(1)
        }
    }
}

func contains(arr []string, str string) bool {
    for _, i := range arr {
        if i == str { return true }
    }
    return false
}

func validateFolders(folders []string) {
    for _, folder := range folders {
        if result, err := exists(folder); !result {
            if err == nil {
                fmt.Println("Error: No such directory", folder)
            } else {
                fmt.Println("Error:", err)

            }
            os.Exit(1)
        }
    }
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}
