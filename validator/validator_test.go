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

    folders := buildForlders()

    errs := Examine(folders)

    if (errs == nil) {
        t.Error("Expected NoSuchDirectoryError error, got nil")
    }
}

func TestExamineUnknownRemote(t *testing.T) {
    setUp()
    defer tearDown()

    os.MkdirAll(filepath.Join(tempFolder, "some", "path"), os.FileMode(0700))

    folders := buildForlders()

    errs := Examine(folders)

    if (errs == nil) {
        t.Error("Expected UnknownRemote error, got nil")
    }
}

func buildForlders() map[string]config.Folder {
    var folders map[string]config.Folder
    folders = make(map[string]config.Folder)

    folders["f1"] = config.Folder{
        Path: "tmp/some/path",
        Strategy: "copy",
        Remotes: []config.Remote{
            {
                Name: "Foo",
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
