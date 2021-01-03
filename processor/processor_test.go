package processor

import (
	"testing"

	"otoclone/config"
	"otoclone/fsnotify"
)

var tempFolder = "tmp"

func TestHandleIgnoreList(t *testing.T) {
    folders := buildForlders([]string{"ignoreme.txt"}, "copy")

    event := fsnotify.FSEvent{
        Folder: "tmp/some/path",
        Event: "happened",
        File: "ignoreme.txt",
    }

    p, errs := Handle(event, folders, false)

    if (errs != nil) {
        t.Errorf("No errors expected, got %v", errs)
    }

    if (p != "") {
        t.Errorf("Expected empty string got %s", p)
    }
}

func TestHandleUnknownBackupStrategy(t *testing.T) {
    folders := buildForlders([]string{}, "sleep")

    event := fsnotify.FSEvent{
        Folder: "tmp/some/path",
        Event: "happened",
        File: "some-file",
    }

    _, errs := Handle(event, folders, false)

    if (errs == nil) {
        t.Error("Expected UnknownBackupStrategy error, got nil")
    }
}

func TestHandleUnwatched(t *testing.T) {
    folders := buildForlders([]string{}, "copy")

    event := fsnotify.FSEvent{
        Folder: "other/than/tmp/some/path",
        Event: "happened",
        File: "some-file",
    }

    _, errs := Handle(event, folders, false)

    if (errs == nil) {
        t.Error("Expected Unwatched error, got nil")
    }
}

func buildForlders(ignoreList []string, strat string) map[string]config.Folder {
    folders := make(map[string]config.Folder)

    folders["f1"] = config.Folder{
        Path: "tmp/some/path",
        Strategy: strat,
        Remotes: []config.Remote{
            {
                Name: "Foo",
                Bucket: "bar",
            },
        },
        IgnoreList: ignoreList,
    }
    return folders
}
