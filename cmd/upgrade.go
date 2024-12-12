package cmd

import (
	"github.com/mirantiscontainers/blueprint-cli/pkg/commands"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// updateCmd represents the apply command
func upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade blueprint operator on the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msgf("Upgrading blueprint at %s", blueprintFlag)
			return commands.Upgrade(&blueprint, kubeConfig)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)

	return cmd
}
