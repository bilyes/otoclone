// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package config

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func MakeFolder() (Folder, error) {
	var qs = []*survey.Question{
		{
			Name:     "path",
			Prompt:   &survey.Input{Message: "Enter the folder path:"},
			Validate: survey.Required,
		},
		{
			Name: "strategy",
			Prompt: &survey.Select{
				Message: "Choose a backup strategy:",
				Options: []string{"copy", "sync"},
			},
		},
	}

	answers := struct {
		Path     string
		Strategy string
		Remote   string
		Bucket   string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return Folder{}, err
	}

	if _, err := os.Stat(answers.Path); os.IsNotExist(err) {
		return Folder{}, fmt.Errorf("The folder %s doesn't exist.", answers.Path)
	}
	if answers.Strategy != "copy" && answers.Strategy != "sync" {
		return Folder{}, fmt.Errorf("The strategy %s is not supported.", answers.Strategy)
	}

	remotes, err := promptRemotes()

	if err != nil {
		return Folder{}, err
	}

	ignoreList := promptIgnoreList()
	excludePattern := promptExcludePattern()

	return Folder{
		Path:           answers.Path,
		Strategy:       answers.Strategy,
		Remotes:        remotes,
		IgnoreList:     ignoreList,
		ExcludePattern: excludePattern,
	}, nil
}

func promptExcludePattern() string {
	choice := ""
	prompt := &survey.Select{
		Message: "Do you want to add a filename pattern to exclude from the backup?",
		Options: []string{"No", "Yes"},
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

func promptRemotes() ([]Remote, error) {
	var remotes []Remote

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
		}

		survey.AskOne(prompt, &choice)
		if choice == "Yes" {
			r, err := promptRemote()

			if err != nil {
				return nil, err
			}
			// TODO Validate remote
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

func promptRemote() (Remote, error) {
	var qs = []*survey.Question{
		{
			Name:     "remote",
			Prompt:   &survey.Input{Message: "Enter the remote to use as destination:"},
			Validate: survey.Required,
		},
		{
			Name:     "bucket",
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
		return Remote{}, err
	}

	return Remote{
		Name:   answers.Remote,
		Bucket: answers.Bucket,
	}, nil
}
