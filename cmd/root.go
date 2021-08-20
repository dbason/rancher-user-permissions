package cmd

import (
	"context"
	"os"

	"github.com/dbason/rancher-user-permissions/cmd/commands"
	"github.com/dbason/rancher-user-permissions/pkg/util"
	"github.com/spf13/cobra"
)

func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rancher-user-permissions",
		Short: "helper for user permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	rootCmd.AddCommand(commands.BuildListCommand())
	return rootCmd
}

func Execute() {
	util.LoadDefaultClientConfig()
	if err := BuildRootCmd().ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
