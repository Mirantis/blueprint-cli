package cmd

import (
	"github.com/mirantiscontainers/blueprint-cli/pkg/commands"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply the blueprint to the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msgf("Applying blueprint at %s", blueprintFlag)
			return commands.Apply(&blueprint, kubeConfig, false, imageRegistry)
		},
	}

	flags := cmd.Flags()
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)
	addImageRegistryFlag(flags)

	return cmd
}
