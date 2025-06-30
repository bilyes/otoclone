// Author: Ilyess Bachiri
// Copyright (c) 2020-present Ilyess Bachiri

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "0.4.4"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Otoclone",
	Long:  `All software has versions. This is Otoclone's.`,
	Run:   printVersion,
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}
