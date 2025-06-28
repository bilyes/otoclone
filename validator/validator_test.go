package validator

import (
	"os"
	"path/filepath"
	"testing"

	"otoclone/config"
	"otoclone/mocks"
)

var tempFolder = "tmp"

func TestExamine(t *testing.T) {
	setUp()
	defer tearDown()

	os.MkdirAll(filepath.Join(tempFolder, "some", "path"), os.FileMode(0700))

	folders := buildForlders("tmp/some/path", "copy", "Foo")

	mockCloner := &mocks.Cloner{}
	mockCloner.On("RemoteIsValid", "Foo").Return(true, nil).Once()
	testVal := &Validator{Cloner: mockCloner}

	errs := testVal.Examine(folders)

	mockCloner.AssertExpectations(t)

	if errs != nil {
		t.Errorf("Expected no errors, got %v", errs)
	}
}

func TestExamineNoDirectory(t *testing.T) {
	setUp()
	defer tearDown()

	folders := buildForlders("tmp/some/path", "copy", "Foo")

	mockCloner := &mocks.Cloner{}
	mockCloner.On("RemoteIsValid", "Foo").Return(true, nil).Once()
	testVal := &Validator{Cloner: mockCloner}

	errs := testVal.Examine(folders)

	mockCloner.AssertExpectations(t)

	if errs == nil {
		t.Error("Expected NoSuchDirectoryError error, got nil")
	}
}

func TestExamineUnknownBackupStrategy(t *testing.T) {
	setUp()
	defer tearDown()

	folders := buildForlders("tmp/some/path", "none", "Foo")

	mockCloner := &mocks.Cloner{}
	mockCloner.On("RemoteIsValid", "Foo").Return(true, nil).Once()
	testVal := &Validator{Cloner: mockCloner}

	errs := testVal.Examine(folders)

	mockCloner.AssertExpectations(t)

	if errs == nil {
		t.Error("Expected UnknownBackupStrategy error, got nil")
	}
}

func TestExamineUnknownRemote(t *testing.T) {
	setUp()
	defer tearDown()

	os.MkdirAll(filepath.Join(tempFolder, "some", "path"), os.FileMode(0700))

	folders := buildForlders("tmp/some/path", "copy", "Foo")

	mockCloner := &mocks.Cloner{}
	mockCloner.On("RemoteIsValid", "Foo").Return(false, nil).Once()
	testVal := &Validator{Cloner: mockCloner}

	errs := testVal.Examine(folders)

	mockCloner.AssertExpectations(t)

	if errs == nil {
		t.Error("Expected UnknownRemote error, got nil")
	}
}

func buildForlders(path string, strat string, remote string) map[string]config.Folder {
	folders := make(map[string]config.Folder)

	folders["f1"] = config.Folder{
		Path:     path,
		Strategy: strat,
		Remotes: []config.Remote{
			{
				Name:   remote,
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
