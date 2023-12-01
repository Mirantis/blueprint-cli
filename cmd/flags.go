package cmd

// PersistenceFlags represents configuration pFlags.
type PersistenceFlags struct {
	LogLevel string
}

func NewPersistenceFlags() *PersistenceFlags {
	return &PersistenceFlags{
		LogLevel: DefaultLogLevel,
	}
}
