package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	kerr "k8s.io/apimachinery/pkg/api/errors"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a kubernetes secret from aws ssm param store",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("error: no secret name provided")
		}
		return cli.Import(args)
	},
}

func (c *CommandOptions) Import(args []string) error {

	c.SetNamespace()
	secretname := args[0]
	secrets, err := c.ssm.GetSecrets(c.ssmPath)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		return fmt.Errorf(fmt.Sprintf("no parameters found at path: %s\n", c.ssmPath))
	}
	if c.encode {
		decoded, err := c.ssm.DecodeSecrets(secrets)
		if err != nil {
			return err
		}
		secrets = decoded
	}
	err = c.k8s.CreateSecret(secretname, secrets, c.tls)
	if err != nil {
		if kerr.IsAlreadyExists(err) {
			if c.overwrite {
				err = c.k8s.UpdateSecret(secretname, secrets)
				if err != nil {
					return err
				}
				fmt.Printf("imported secret: %s\n", secretname)
			}
		}
		return err
	}
	fmt.Printf("imported secret: %s\n", secretname)
	return nil
}
