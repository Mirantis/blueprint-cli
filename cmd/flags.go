package cmd

const (
	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"

	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"
)

// PersistenceFlags represents configuration pFlags.
type PersistenceFlags struct {
	LogLevel string
}

func NewPersistenceFlags() *PersistenceFlags {
	return &PersistenceFlags{
		LogLevel: DefaultLogLevel,
	}
}
