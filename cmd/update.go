package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"boundless-cli/internal/boundless"
)

// updateCmd represents the apply command
func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the cluster according to the blueprint",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd)
		},
	}

	flags := cmd.Flags()
	addConfigFlags(flags)
	addKubeFlags(flags)

	return cmd
}

func runUpdate(cmd *cobra.Command) error {
	log.Info().Msgf("Applying Boundless Operator resources")
	if err := boundless.ApplyBlueprint(kubeConfig, blueprint); err != nil {
		return fmt.Errorf("failed to update components: %w", err)
	}

	log.Info().Msgf("Finished updating Boundless Operator")
	return nil
}
