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
		Short:   "Update the blueprints to the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		Run: func(cmd *cobra.Command, args []string) {
			err := updateFunc(cmd)
			if err != nil {
				return
			}
		},
	}

	flags := cmd.Flags()
	addConfigFlags(flags)
	addKubeFlags(flags)
	
	return cmd
}

func updateFunc(cmd *cobra.Command) error {
	// install components
	log.Info().Msgf("Applying Boundless Operator resource")
	if err := boundless.ApplyBlueprint(kubeConfig, blueprint); err != nil {
		return fmt.Errorf("failed to install components: %w", err)
	}

	log.Info().Msgf("Finished installing Boundless Operator")
	return nil
}
