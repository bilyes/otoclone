package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tempFolder = "tmp"

func TestArrayContainsTrue(t *testing.T) {
    r := ArrayContains([]string{"foo", "bar"}, "foo")

    assert.Equal(t, r, true, "ArrayContains should have returned true")
}

func TestArrayContainsFalse(t *testing.T) {
    r := ArrayContains([]string{"foo", "bar"}, "intruder")

    assert.Equal(t, r, false, "ArrayContains should have returned false")
}

func TestPathExistsFalse(t *testing.T) {
    setUp()
    defer tearDown()

    r, err := PathExists("not-found.txt")

    assert.Equal(t, r, false, "PathExists should have returned false")
    assert.Equal(t, err, nil, "PathExists should not have returned an error")
}

func TestPathExistsTrue(t *testing.T) {
    setUp()
    defer tearDown()

    os.MkdirAll(filepath.Join(tempFolder, "some", "path"), os.FileMode(0700))

    r, err := PathExists("tmp/some/path")

    assert.Equal(t, r, true, "PathExists should have returned true")
    assert.Equal(t, err, nil, "PathExists should not have returned an error")
}

func setUp() {
    os.Mkdir(tempFolder, os.FileMode(0700))
}

func tearDown() {
    os.RemoveAll(tempFolder)
}
