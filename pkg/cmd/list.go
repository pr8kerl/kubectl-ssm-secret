package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:          "list <ssm parameter store path> ...",
	Short:        "list ssm parameters by path ",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no ssm param store path provided")
		}
		return cli.List(args)
	},
}

func (c *CommandOptions) List(args []string) error {
	for _, key := range args {
		secrets, err := c.ssm.GetSecrets(key)
		if err != nil {
			return err
		}

		if len(secrets) == 0 {
			return fmt.Errorf(fmt.Sprintf("no parameters found at path: %s", key))
		}
		for k, v := range secrets {
			fmt.Printf("%s: %s\n", k, v)
		}

	}
	return nil
}
