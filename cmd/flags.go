package cmd

import "github.com/mirantiscontainers/boundless-cli/pkg/constants"

// PersistenceFlags represents configuration pFlags.
type PersistenceFlags struct {
	LogLevel string
}

func NewPersistenceFlags() *PersistenceFlags {
	return &PersistenceFlags{
		LogLevel: constants.DefaultLogLevel,
	}
}
