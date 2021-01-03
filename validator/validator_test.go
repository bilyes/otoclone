package validator

import (
	"os"
	"otoclone/config"
	"path/filepath"

	"testing"
)

var tempFolder = "tmp"

func TestExamineNoDirectory(t *testing.T) {
    setUp()
    defer tearDown()

    folders := buildForlders("tmp/some/path", "copy", "Foo")

    errs := Examine(folders)

    if (errs == nil) {
        t.Error("Expected NoSuchDirectoryError error, got nil")
    }
}

func TestExamineUnknownBackupStrategy(t *testing.T) {
    setUp()
    defer tearDown()

    folders := buildForlders("tmp/some/path", "none", "Foo")

    errs := Examine(folders)

    if (errs == nil) {
        t.Error("Expected UnknownBackupStrategy error, got nil")
    }
}

func TestExamineUnknownRemote(t *testing.T) {
    setUp()
    defer tearDown()

    os.MkdirAll(filepath.Join(tempFolder, "some", "path"), os.FileMode(0700))

    folders := buildForlders("tmp/some/path", "copy", "Foo")

    errs := Examine(folders)

    if (errs == nil) {
        t.Error("Expected UnknownRemote error, got nil")
    }
}

func buildForlders(path string, strat string, remote string) map[string]config.Folder {
    folders := make(map[string]config.Folder)

    folders["f1"] = config.Folder{
        Path: path,
        Strategy: strat,
        Remotes: []config.Remote{
            {
                Name: remote,
                Bucket: "bar",
            },
        },
    }
    return folders
}

func setUp() {
    os.Mkdir(tempFolder, os.FileMode(0700))
}

func tearDown() {
    os.RemoveAll(tempFolder)
}
