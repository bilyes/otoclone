// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"otoclone/fsnotify"
	"otoclone/rclone"
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
        event, err := fsnotify.Watch(folders)

        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        process(event, ignoreList, folders, remotes)
    }
}

func process(event fsnotify.FSEvent, ignoreList []string, folders []string, remotes []string) {
    if contains(ignoreList, event.File) { return }

    folder := ""

    for _, f := range folders {
        if strings.HasPrefix(event.Folder, f) {
            folder = f
            break
        }
    }

    if folder == "" {
        fmt.Println("Error: Unwatched file or folder", event.Folder)
        return
    }

    for _, r := range remotes {
        err := rclone.Copy(folder, r, filepath.Base(folder))
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("Backed up", folder)
}

func validateRemotes(remotes []string) {
    for _, r := range remotes {
        isValid, err := rclone.RemoteIsValid(r)
        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        if !isValid {
            fmt.Print("Error: Unknown remote ", r)
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
