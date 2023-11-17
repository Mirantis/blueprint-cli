package cmd

const (
	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"
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
