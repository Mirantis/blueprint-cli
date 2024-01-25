package constants

const (
	// ManifestUrlLatest is the URL of the latest manifest YAML for the Boundless Operator
	ManifestUrlLatest = "https://raw.githubusercontent.com/mirantiscontainers/boundless/main/deploy/static/boundless-operator.yaml"

	// NamespaceBoundless is the system namespace where the Boundless Operator and its components are installed
	NamespaceBoundless = "boundless-system"

	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"

	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"
)
