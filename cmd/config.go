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
        addFolder()
    case remF:
    case editF:
    default:
        fmt.Println("Unknown option", choice)
    }
}

func addFolder() {
    var qs = []*survey.Question{
        {
            Name:     "path",
            Prompt:   &survey.Input{Message: "Enter the folder path:"},
            Validate: survey.Required,
            //Transform: survey.Title,
        },
        {
            Name: "strategy",
            Prompt: &survey.Select{
                Message: "Choose a backup strategy:",
                Options: []string{"copy", "sync"},
                //Default: "copy",
            },
        },
    }

    answers := struct {
        Path string
        Strategy string
        Remote string
        Bucket string
    }{}

    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // TODO check if path exists

    remotes, err := promptRemotes()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    ignoreList := promptIgnoreList()
    excludePattern := promptExcludePattern()

    fmt.Printf("You selected")
    fmt.Printf("Path: %s", answers.Path)
    fmt.Printf("Strategy: %s", answers.Strategy)
    fmt.Printf("Remotes: %v", remotes)
    fmt.Printf("Ignore list: %v", ignoreList)
    fmt.Printf("Exclude pattern: %s", excludePattern)
}

func promptExcludePattern() string {
    choice := ""
    prompt := &survey.Select{
        Message: "Do you want to add a filename pattern to exclude from the backup?",
        Options: []string{"No", "Yes"},
        //Default: "No",
    }
    survey.AskOne(prompt, &choice)

    p := ""
    if choice == "Yes" {
        prompt := &survey.Input{
            Message: "Enter the file name to ignore:",
        }
        survey.AskOne(prompt, &p)
    }
    return p
}

func promptRemotes() ([]config.Remote, error) {
    var remotes []config.Remote

    r, err := promptRemote()

    if err != nil {
        return nil, err
    }
    remotes = append(remotes, r)

    choice := "Yes"
    for choice == "Yes" {
        prompt := &survey.Select{
            Message: "Do you want to add another remote?",
            Options: []string{"No", "Yes"},
            //Default: "No",
        }

        survey.AskOne(prompt, &choice)
        if choice == "Yes" {
            r, err := promptRemote()

            if err != nil {
                return nil, err
            }
            // TODO check if remote is configured
            remotes = append(remotes, r)
        }
    }

    return remotes, nil
}

func promptIgnoreList() []string {
    var list []string

    choice := "Yes"

    for choice == "Yes" {
        prompt := &survey.Select{
            Message: "Do you want to add a file to ignore on the watch list?",
            Options: []string{"No", "Yes"},
        }
        survey.AskOne(prompt, &choice)

        if choice == "Yes" {
            i := promptIgnoreListItem()
            list = append(list, i)
        }
    }

    return list
}

func promptIgnoreListItem() string {
    i := ""
    prompt := &survey.Input{
        Message: "Enter the file name to ignore:",
    }
    survey.AskOne(prompt, &i)
    return i
}

func promptRemote() (config.Remote, error) {

    var qs = []*survey.Question{
        {
            Name:     "remote",
            Prompt:   &survey.Input{Message: "Enter the remote to use as destination:"},
            Validate: survey.Required,
        },
        {
            Name: "bucket",
            Prompt:   &survey.Input{Message: "Enter the path on the remote:"},
            Validate: survey.Required,
        },

    }
    answers := struct {
        Remote string
        Bucket string
    }{}

    err := survey.Ask(qs, &answers)
    if err != nil {
        return config.Remote{}, err

    }

    return config.Remote{
        Name: answers.Remote,
        Bucket: answers.Bucket,
    }, nil
}
