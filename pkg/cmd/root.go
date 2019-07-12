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
	%[1]s list /param/path/foo

	# create a kubernetes secret called foo using key/values from parameter store path /param/path/foo
	%[1]s create /param/path/foo

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
		IOStreams:   streams,
	}
}

func init() {
	cli = NewCommandOptions(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&cli.toSsm, "secret-to-ssm", cli.toSsm, "copy from k8s secret to aws ssm param store path")
	createCmd.Flags().StringVar(&cli.ssmPath, "ssm-path", cli.ssmPath, "ssm parameter store path to read data from")
	createCmd.Flags().BoolVar(&cli.overwrite, "overwrite", cli.overwrite, "copy from k8s secret to aws ssm param store path")
	cli.configFlags.AddFlags(rootCmd.Flags())
	cli.configFlags.AddFlags(listCmd.Flags())
	cli.configFlags.AddFlags(createCmd.Flags())
}

var rootCmd = &cobra.Command{
	Use:              "ssm-secret list|create|update /param/store/path [flags]",
	Short:            "view or create secret from param store",
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
