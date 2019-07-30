package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the ssm-secret version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssm-secret v1.0.0")
	},
}
