package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"otoclone/config"
	"otoclone/fsnotify"
	"otoclone/processor"
	"otoclone/validator"
)

var (
    configFile = ""
    verbose = false
)

func init() {
    parseFlags(rootCmd.Flags())
}

var rootCmd = &cobra.Command{
    Use: "otoclone",
    Short: "Otoclone is an automatic backup and sync utility",
    Long: `A backup and sync utility that automatically reacts to filesystem
    events and copies watched folders to various remotes.`,
    Run: watch,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func watch(cmd *cobra.Command, args []string) {
    folders := loadConfig(configFile)
    validator.Examine(folders)

    paths := extractPaths(folders)
    foldersToWatch := strings.Join(paths, " ")

    for {
        fmt.Println("Watching", foldersToWatch)
        event, err := fsnotify.Watch(paths)

        if err != nil {
            fmt.Println("Error:",  err)
            os.Exit(1)
        }

        path, errors := processor.Handle(event, folders, verbose)

        if errors != nil {
            fmt.Println("Errors:")
            for e := range errors {
                fmt.Println(e)
            }
            continue
        }

        if path == "" {
            fmt.Println("Ignored", event.File)
        } else {
            fmt.Println("Backed up", path)
        }
    }
}

func parseFlags(flags *pflag.FlagSet) {
    flags.BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
    flags.StringVarP(&configFile, "config", "c", "", "Path to the configuration file")
}

func loadConfig(configFile string) map[string]config.Folder {
    var folders map[string]config.Folder
    var err error

    if configFile != "" {
        folders, err = config.LoadFile(configFile)
    } else {
        folders, err = config.Load()
    }

    if err != nil {
        fmt.Println("Error:",  err)
        os.Exit(1)
    }
    return folders
}

func extractPaths(folders map[string]config.Folder) []string {
    fKeys := map[string]bool{}
    var paths []string

    for _, f := range folders {
        if _, value := fKeys[f.Path]; !value {
            fKeys[f.Path] = true
            paths = append(paths, f.Path)
        }
    }
    return paths
}

