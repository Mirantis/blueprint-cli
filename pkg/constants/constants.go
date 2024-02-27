package constants

import "time"

const (
	// ManifestUrlLatest is the URL of the latest manifest YAML for the Boundless Operator
	ManifestUrlLatest = "https://raw.githubusercontent.com/mirantiscontainers/boundless/main/deploy/static/boundless-operator.yaml"

	// NamespaceBoundless is the system namespace where the Boundless Operator and its components are installed
	NamespaceBoundless = "boundless-system"

	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"

	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"

	// DryRunWaitInterval is the interval to wait between checks of resources when performing a dry run
	DryRunWaitInterval = 2 * time.Second

	// DryRunTimeout is the timeout for dry run operations
	DryRunTimeout = 2 * time.Minute
)
