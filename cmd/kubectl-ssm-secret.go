package main

import (
	"github.com/spf13/pflag"

	"github.com/pr8kerl/kubectl-ssm-secret/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-ssm-secret", pflag.ExitOnError)
	pflag.CommandLine = flags
	cmd.Execute()
}
