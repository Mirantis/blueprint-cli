package constants

import "time"

const (
	// ManifestUrlLatest is the URL of the latest manifest YAML for the Boundless Operator
	ManifestUrlLatest = "https://raw.githubusercontent.com/mirantiscontainers/boundless/main/deploy/static/boundless-operator.yaml"

	// NamespaceBlueprint is the system namespace where the Boundless Operator and its components are installed
	NamespaceBlueprint = "blueprint-system"

	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"

	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"

	// DryRunWaitInterval is the interval to wait between checks of resources when performing a dry run
	DryRunWaitInterval = 2 * time.Second

	// DryRunTimeout is the timeout for dry run operations
	DryRunTimeout = 2 * time.Minute

	BlueprintOperatorDeployment = "blueprint-operator-controller-manager"

	// These semver regex come from the official semver spec: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	// but they have been modified into a version without a leading v, a version with a leading v, and a version where the leading v is optional
	SemverRegexWithV     = `^[v](0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`
	SemverRegexWithoutV  = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`
	SemverRegexOptionalV = `^[v]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`
	// The k0s semver regex is the same as the optional v semver regex, but with an optional "+k0s.0" at the end
	K0sSemverRegex = `^[v]?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:\+(k[0-9a-zA-Z-]s+(?:\.[0-9a-zA-Z-]+)*))?$`
)
