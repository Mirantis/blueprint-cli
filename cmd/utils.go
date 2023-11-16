package cmd

import (
	"github.com/spf13/cobra"
)

func actions(fs ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range fs {
			if err := f(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
