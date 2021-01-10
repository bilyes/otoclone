// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"otoclone/config"
	"otoclone/validator"
	"otoclone/rclone"
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
    //editF := "Edit an existing folder"
    //remF := "Remove an existing folder"
    q := "Quit"

    choice := ""
    prompt := &survey.Select{
        Message: "What do you want to do?",
        Options: []string{
            addF,
            //editF,
            //remF,
            q,
        },
    }
    survey.AskOne(prompt, &choice)

    switch choice {
    case addF:
        addFolder()
    //case remF:
    //case editF:
    case q:
        fmt.Println("Bye.")
    default:
        fmt.Println("Unknown option", choice)
    }
}

func addFolder() {
    f, err := config.MakeFolder()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    val := &validator.Validator{ Cloner: &rclone.Rclone{} }
    errs := val.ExamineOne(f)
    if errs != nil && len(errs) > 0 {
        for _, e := range errs {
            fmt.Println(e.Error())
        }
        return
    }

    // TODO confirm all is good
    if err := config.Write(f); err != nil {
        fmt.Println(err.Error())
        return
    }
}
