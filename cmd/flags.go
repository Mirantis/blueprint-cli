package cmd

const (
	// DefaultConfigFilename represents the default blueprints filename.
	DefaultConfigFilename = "blueprint.yaml"
)

// PersistenceFlags represents configuration pFlags.
type PersistenceFlags struct {
	Debug bool
}

func NewPersistenceFlags() *PersistenceFlags {
	return &PersistenceFlags{
		Debug: false,
	}
}
