// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"otoclone/config"
	"otoclone/processor"
	"otoclone/rclone"
)

var folder string

func init() {
    backupCmd.Flags().StringVarP(&folder, "folder", "f", "", "The folder to backup")
    rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
    Use: "backup",
    Short: "Run a one time backup",
    Long: `Execute a backup of one or all of the folders defined in the configuration and exit`,
    Run: backup,
}

func backup(cmd *cobra.Command, args []string) {
    folders := loadFolders()

    proc := &processor.Processor{ Cloner: &rclone.Rclone{} }

    if folder != "" {
        folders = extractFolder(folders, folder)
        if len(folders) == 0 {
            fmt.Printf("Folder %s not found in the configuration. To add it run: otoclone config", folder)
            os.Exit(1)
        }
    }

    if errs := proc.Backup(folders, verbose); len(errs) > 0 {
        for _, e := range errs {
            fmt.Println("Error:", e)
        }
    }
}

func extractFolder(folders map[string]config.Folder, folder string) map[string]config.Folder {
    fols := make(map[string]config.Folder)

    for k, f := range folders {
        if k == folder {
            fols[k] = f
            return fols
        }
    }

    return fols
}
