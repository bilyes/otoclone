// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package processor

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/rclone"

	"github.com/bilyes/conman"
)

type Processor struct {
	Concurrency int64
	Cloner      rclone.Cloner
}

// Handles a FSNotify Event
func (p *Processor) Handle(
	ctx context.Context,
	event fsnotify.FSEvent,
	folders map[string]config.Folder,
	verbose bool,
) (string, []error, error) {
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
		return "", []error{&UnwatchedError{event.Folder}}, nil
	}

	if slices.Contains(subject.IgnoreList, event.File) {
		return "", nil, nil
	}

	flds := make(map[string]config.Folder)
	flds[fKey] = subject

	errors, err := p.Backup(ctx, flds, verbose)
	return subject.Path, errors, err
}

type cloningTask struct {
	cloner          rclone.Cloner
	source          string
	destinationPath string
	destination     string
	flags           rclone.Flags
	strategy        string
}

func (c *cloningTask) Execute(ctx context.Context) (any, error) {
	switch c.strategy {
	case "copy":
		return nil, c.cloner.Copy(c.source, c.destinationPath, c.destination, c.flags)
	case "sync":
		return nil, c.cloner.Sync(c.source, c.destinationPath, c.destination, c.flags)
	default:
		return nil, &UnknownBackupStrategyError{c.strategy}
	}
}

// Backup a list of folders
func (p *Processor) Backup(ctx context.Context, folders map[string]config.Folder, verbose bool) ([]error, error) {
	folderCount := len(folders)
	if folderCount == 0 {
		return nil, nil
	}
	concurrency := p.Concurrency
	if concurrency <= 0 {
		concurrency = int64(folderCount)
	}
	if concurrency < 2 {
		// conman requires at least 2 concurrency
		concurrency = 2
	}

	cm, err := conman.New[any](concurrency)
	if err != nil {
		return nil, err
	}

	for _, folder := range folders {

		flags := rclone.Flags{
			Verbose: verbose,
			Exclude: folder.ExcludePattern,
		}

		for _, remote := range folder.Remotes {
			cm.Run(
				ctx,
				&cloningTask{
					cloner:          p.Cloner,
					source:          folder.Path,
					destinationPath: remote.Name,
					destination:     remote.Bucket,
					flags:           flags,
					strategy:        folder.Strategy,
				},
			)
		}
	}

	if err := cm.Wait(ctx); err != nil {
		return nil, err
	}

	return cm.Errors(), nil
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
