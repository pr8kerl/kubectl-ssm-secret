package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "snapshot"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the ssm-secret version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("ssm-secret %s", version))
	},
}
