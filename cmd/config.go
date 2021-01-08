// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package cmd

import (
    "fmt"

    "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

    "otoclone/config"
)

func init() {
    rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
    Use: "config",
    Short: "Update configuration",
    Long: `Create or make modifications to Otoclon's configuration global configuration`,
    Run: configure,
}

func configure(cmd *cobra.Command, args []string) {
    addF := "Add a new folder"
    editF := "Edit an existing folder"
    remF := "Remove an existing folder"
    q := "Quit"

    choice := ""
    prompt := &survey.Select{
        Message: "What do you want to do?",
        Options: []string{addF, editF, remF, q},
    }
    survey.AskOne(prompt, &choice)

    switch choice {
    case addF:
        f, err := config.MakeFolder()
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        // TODO confirm all is good
        if err := config.Write(f); err != nil {
            fmt.Println(err.Error())
            return
        }
        fmt.Printf("You selected")
        fmt.Printf("Path: %s", f.Path)
        fmt.Printf("Strategy: %s", f.Strategy)
        fmt.Printf("Remotes: %v", f.Remotes)
        fmt.Printf("Ignore list: %v", f.IgnoreList)
        fmt.Printf("Exclude pattern: %s", f.ExcludePattern)
    case remF:
    case editF:
    default:
        fmt.Println("Unknown option", choice)
    }
}

