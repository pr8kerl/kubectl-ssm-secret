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
		return cli.List(args)
	},
}

func (c *CommandOptions) List(args []string) error {
	c.SetNamespace()
	err := c.ListSsmSecrets()
	if err != nil {
		return err
	}
	err = c.ListK8sSecrets(args)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandOptions) ListSsmSecrets() error {
	if len(c.ssmPath) > 0 {
		secrets, err := c.ssm.GetSecrets(c.ssmPath)
		if err != nil {
			return err
		}
		if len(secrets) == 0 {
			return fmt.Errorf(fmt.Sprintf("no parameters found at path: %s", c.ssmPath))
		}
		if c.toEnvironment {
			for k, v := range secrets {
				fmt.Printf("%s=%s\n", k, v)
			}
		} else {
			for k, v := range secrets {
				fmt.Printf("ssm:%s/%s: %s\n", c.ssmPath, k, v)
			}
		}
	}
	return nil
}

func (c *CommandOptions) ListK8sSecrets(args []string) error {
	for _, key := range args {
		secrets, err := c.k8s.GetSecret(key)
		if err != nil {
			return err
		}
		if len(secrets) == 0 {
			return fmt.Errorf(fmt.Sprintf("no secret data found in secret: %s", key))
		}
		if c.toEnvironment {
			for k, v := range secrets {
				fmt.Printf("%s=%s\n", k, v)
			}
		} else {
			for k, v := range secrets {
				fmt.Printf("k8s:%s/%s/%s: %s\n", c.namespace, key, k, v)
			}
		}
	}
	return nil
}
