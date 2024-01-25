package constants

// The constants in this file represent the values of the fields allowed in a
// Blueprint

const (
	// kinds
	KindBlueprint = `Blueprint`

	// providers
	// ProviderK0S is the name of the k0s distro
	ProviderK0s = "k0s"
	// ProviderKind is the name of the kind distro
	ProviderKind = "kind"
	// ProviderExisting is the name of an existing unofficial distro
	ProviderExisting = "existing"

	// addons
	// AddonManifest is the name of the manifest addon
	AddonManifest = "manifest"
	// AddonChart is the name of the chart addon
	AddonChart = "chart"
)
