package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mirantiscontainers/blueprint-cli/pkg/commands"
	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a blueprint file template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Init(constants.ProviderK0s)
		},
	}

	cmd.AddCommand(initKindCmd())
	cmd.AddCommand(initK0sCmd())

	return cmd
}

func initKindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kind",
		Short: "Creates a kind blueprint file template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Init(constants.ProviderKind)
		},
	}
}

func initK0sCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "k0s",
		Short: "Creates a k0s blueprint file template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Init(constants.ProviderK0s)
		},
	}
}
