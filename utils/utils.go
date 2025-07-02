// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package utils

import (
	"os"
)

// Check if a path exists on the filesystem
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
