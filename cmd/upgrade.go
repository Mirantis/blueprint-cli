package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"boundless-cli/internal/k8s"
)

// updateCmd represents the apply command
func upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade boundless operator on the cluster",
		Args:    cobra.NoArgs,
		PreRunE: actions(loadBlueprint, loadKubeConfig),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(cmd)
		},
	}

	flags := cmd.Flags()
	addOperatorUriFlag(flags)
	addBlueprintFileFlags(flags)
	addKubeFlags(flags)
	return cmd
}

func runUpgrade(cmd *cobra.Command) error {
	log.Info().Msgf("Upgrading Boundless Operator using manifest file %q", operatorUri)
	if err := k8s.ApplyYaml(kubeConfig, operatorUri); err != nil {
		return fmt.Errorf("failed to upgrade boundless operator: %w", err)
	}

	log.Info().Msgf("Finished updating Boundless Operator")
	return nil
}
