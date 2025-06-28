// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"otoclone/config"
	"otoclone/rclone"
	"otoclone/utils"
	"otoclone/validator"
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
    if configFile != "" {
        ok, err := utils.PathExists(configFile)

        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        if !ok {
            fmt.Println("The provided config file doesn't exist:", configFile)
            return
        }
    }

    addF := "Add a new folder"
    editF := "Edit an existing folder"
    remF := "Remove an existing folder"
    q := "Quit"

    choice := ""
    prompt := &survey.Select{
        Message: "What do you want to do?",
        Options: []string{
            addF,
            editF,
            remF,
            q,
        },
    }
    survey.AskOne(prompt, &choice)

    switch choice {
    case addF:
        addFolder()
    case remF:
        deleteFolder()
    case editF:
        editFolder()
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
    if len(errs) > 0 {
        for _, e := range errs {
            fmt.Println(e.Error())
        }
        return
    }

    fmt.Println("Adding the following folder configuration:")
    fmt.Printf("%+v\n", f)
    var choice string
    prompt := &survey.Select{
        Message: "Do you confirm?",
        Options: []string{"Yes", "No"},
    }
    survey.AskOne(prompt, &choice)
    if choice == "No" {
        return
    }

    if configFile != "" {
        err = config.WriteTo(f, configFile)
    } else {
        err = config.Write(f)
    }

    if err != nil {
        fmt.Println(err.Error())
    }
}

func deleteFolder() {
    // List existing folder and prompt for selection
    var fols map[string]config.Folder
    var err error
    if configFile != "" {
        fols, err = config.LoadFile(configFile)
    } else {
        fols, err = config.Load()
    }

    if err != nil {
        fmt.Println("Error:", err)
    }

    var keys []string
    for k := range fols {
        keys = append(keys, k)
    }
    var folder string

    prompt := &survey.Select{
        Message: "Which folder do you want to remove?",
        Options: keys,
    }
    survey.AskOne(prompt, &folder)

    // Prompt for confirmation
    var choice string
    prompt = &survey.Select{
        Message: "Do you confirm?",
        Options: []string{"Yes", "No"},
    }
    survey.AskOne(prompt, &choice)
    if choice == "No" {
        return
    }

    // delete selected folder
    if configFile != "" {
        err = config.RemoveFrom(folder, configFile)
    } else {
        err = config.Remove(folder)
    }
    if err != nil {
        fmt.Println(err.Error())
    }
}

func editFolder() {
    var loc string

    if configFile != "" {
        loc = configFile
    } else {
        var err error
        loc, err = config.Location()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
    }

    fmt.Println("To make changes to the configuration, edit the config file", loc)
}
