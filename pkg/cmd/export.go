package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export a kubernetes secret to aws ssm param store",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("error: no secret name provided")
		}
		return cli.Export(args)
	},
}

func (c *CommandOptions) Export(args []string) error {

	c.SetNamespace()
	secretname := args[0]
	secrets, err := c.k8s.GetSecret(secretname)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("no data found in secret: %s\n", secretname))
	}
	err = c.ssm.PutSecrets(c.ssmPath, secrets, c.overwrite)
	if err != nil {
		return err
	}
	fmt.Printf("exported secret: %s\n", secretname)
	return nil
}
