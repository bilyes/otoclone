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

    "github.com/bilyes/conman"
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

type cloningTask struct {
    cloner rclone.Cloner
    source string
    destinationPath string
    destination string
    flags rclone.Flags
    strategy string
}

func (c *cloningTask) Execute() (interface{}, error) {
    switch c.strategy {
    case "copy":
        return nil, c.cloner.Copy(c.source, c.destinationPath, c.destination, c.flags)
    case "sync":
        return nil, c.cloner.Sync(c.source, c.destinationPath, c.destination, c.flags)
    default:
        return nil,  &UnknownBackupStrategyError{c.strategy}
    }
}

// Backup a list of folders
func (p *Processor) Backup(folders map[string]config.Folder, verbose bool) []error {
    cm := conman.New(100)

    for _, folder := range folders {

        flags := rclone.Flags{
            Verbose: verbose,
            Exclude: folder.ExcludePattern,
        }

        for _, remote := range folder.Remotes {
            cm.Run(
                &cloningTask{
                    cloner: p.Cloner,
                    source: folder.Path,
                    destinationPath: remote.Name,
                    destination: remote.Bucket,
                    flags: flags,
                    strategy: folder.Strategy,
                },
            )
        }
    }

    cm.Wait()

    return cm.Errors()
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
