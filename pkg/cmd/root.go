package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pr8kerl/kubectl-ssm-secret/pkg/k8s"
	"github.com/pr8kerl/kubectl-ssm-secret/pkg/ssm"
)

var (
	commandExample = `
	# view the parameter store keys and values located in parameter store path /param/path/foo
	%[1]s list /param/path/foo

	# import to a kubernetes secret called foo from key/values stored at parameter store path /param/path/foo
	%[1]s import foo --ssm-path /param/path/foo

	# export a kubernetes secret called foo to aws ssm parameter store path /param/path/foo
	%[1]s export foo --ssm-path /param/path/foo

	# display the plugin version
	%[1]s version
`
	cli *CommandOptions
)

type CommandOptions struct {
	ssmPath   string
	toSsm     bool
	args      []string
	ssm       *ssm.Client
	k8s       *k8s.K8sClient
	overwrite bool
	encode    bool
	tls       bool
	namespace string
}

// NewCommandOptions provides an instance of CommandOptions with default values
func NewCommandOptions() *CommandOptions {
	svc, err := ssm.New()
	if err != nil {
		fmt.Printf("error: cannot create aws ssm client: %s\n", err)
		os.Exit(1)
	}
	kconfig, err := k8s.NewK8sConfig()
	if err != nil {
		fmt.Printf("error: cannot configure k8s client: %s\n", err)
		os.Exit(1)
	}
	kclient, err := k8s.NewK8sClientFromConfig(kconfig)
	if err != nil {
		fmt.Printf("error: cannot init k8s client: %s\n", err)
		os.Exit(1)
	}
	return &CommandOptions{
		toSsm:     false,
		ssmPath:   "",
		ssm:       svc,
		k8s:       kclient,
		overwrite: false,
		encode:    false,
		tls:       false,
		namespace: "",
	}
}

func (c *CommandOptions) SetNamespace() {
	c.k8s.SetNamespace(c.namespace)
}

func init() {
	cli = NewCommandOptions()
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	importCmd.Flags().StringVarP(&cli.ssmPath, "ssm-path", "s", cli.ssmPath, "ssm parameter store path to read data from")
	importCmd.MarkFlagRequired("ssm-path")
	importCmd.Flags().BoolVarP(&cli.overwrite, "overwrite", "o", cli.overwrite, "if k8s secret exists, overwite its values with those from param store")
	importCmd.Flags().BoolVarP(&cli.encode, "decode", "d", cli.encode, "treat store values in param store as gzipped, base64 encoded strings")
	importCmd.Flags().BoolVarP(&cli.tls, "tls", "t", cli.tls, "import ssm param store values to k8s tls secret")
	importCmd.Flags().StringVarP(&cli.namespace, "namespace", "n", cli.namespace, "kubernetes namespace")
	exportCmd.Flags().StringVarP(&cli.ssmPath, "ssm-path", "s", cli.ssmPath, "ssm parameter store path to write data to")
	exportCmd.MarkFlagRequired("ssm-path")
	exportCmd.Flags().BoolVarP(&cli.overwrite, "overwrite", "o", cli.overwrite, "if parameter store key exists, overwite its values with those from k8s secret")
	exportCmd.Flags().BoolVarP(&cli.encode, "encode", "e", cli.encode, "gzip, base64 encode values in parameter store")
	exportCmd.Flags().StringVarP(&cli.namespace, "namespace", "n", cli.namespace, "kubernetes namespace")
}

var rootCmd = &cobra.Command{
	Use:              "ssm-secret list|import|export secret [flags]",
	Short:            "view or import/export k8s secrets from/to aws ssm param store",
	Example:          fmt.Sprintf(commandExample, "kubectl ssm-secret"),
	SilenceUsage:     true,
	TraverseChildren: true,
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no sub command provided")
		}
		return nil
	},
}

// Execute is used to run the command logic in a vein similar to the Cobra package
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
