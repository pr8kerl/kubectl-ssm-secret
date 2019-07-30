package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pr8kerl/kubectl-ssm-secret/pkg/k8s"
	"github.com/pr8kerl/kubectl-ssm-secret/pkg/ssm"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	commandExample = `
	# view the param keys and values located in parameter store path /param/path/foo
	%[1]s list --ssm-path=/param/path/foo

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
	configFlags *genericclioptions.ConfigFlags

	ssmPath   string
	toSsm     bool
	args      []string
	ssm       *ssm.Client
	k8s       *k8s.K8sClient
	overwrite bool
	encode    bool
	tls       bool

	genericclioptions.IOStreams
}

// NewCommandOptions provides an instance of CommandOptions with default values
func NewCommandOptions(streams genericclioptions.IOStreams) *CommandOptions {
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
		configFlags: genericclioptions.NewConfigFlags(true),
		toSsm:       false,
		ssmPath:     "",
		ssm:         svc,
		k8s:         kclient,
		overwrite:   false,
		encode:      false,
		tls:         false,
		IOStreams:   streams,
	}
}

func init() {
	cli = NewCommandOptions(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	importCmd.Flags().StringVar(&cli.ssmPath, "ssm-path", cli.ssmPath, "ssm parameter store path to read data from")
	importCmd.MarkFlagRequired("ssm-path")
	importCmd.Flags().BoolVar(&cli.overwrite, "overwrite", cli.overwrite, "if k8s secret exists, overwite its values with those from param store")
	importCmd.Flags().BoolVar(&cli.encode, "decode", cli.encode, "treat store values in param store as gzipped, base64 encoded strings")
	importCmd.Flags().BoolVar(&cli.encode, "tls", cli.tls, "import ssm param store values to k8s tls secret")
	exportCmd.Flags().StringVar(&cli.ssmPath, "ssm-path", cli.ssmPath, "ssm parameter store path to write data to")
	exportCmd.MarkFlagRequired("ssm-path")
	exportCmd.Flags().BoolVar(&cli.overwrite, "overwrite", cli.overwrite, "if parameter store key exists, overwite its values with those from k8s secret")
	exportCmd.Flags().BoolVar(&cli.encode, "encode", cli.encode, "gzip, base64 encode values in parameter store")
	cli.configFlags.AddFlags(rootCmd.Flags())
	cli.configFlags.AddFlags(listCmd.Flags())
	cli.configFlags.AddFlags(importCmd.Flags())
	cli.configFlags.AddFlags(exportCmd.Flags())
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
