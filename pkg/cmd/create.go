package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a kubernetes secret from aws ssm param store",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("error: no ssm param store path provided")
		}
		return cli.Create(args)
	},
}

func (c *CommandOptions) Create(args []string) error {

	secretname := args[0]
	secrets, err := c.ssm.GetSecrets(c.ssmPath)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("no parameters found at path: %s\n", c.ssmPath))
	}
	fmt.Printf("creating secret: %s\n", secretname)
	err = c.k8s.CreateSecret(secretname, secrets)
	if err != nil {
		return err
	}
	return nil
}
