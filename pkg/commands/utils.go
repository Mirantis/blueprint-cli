package commands

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/mirantiscontainers/boundless-cli/pkg/constants"
	"github.com/rs/zerolog/log"
)

var operatorReleaseUri = "https://github.com/MirantisContainers/blueprint/releases/download/%s/blueprint-operator.yaml"

// determineOperatorUri determines the URI of the operator based on the version
// If the version is a valid URI, it will be returned as is for dev testing
// Otherwise, it will be assumed to be a version and the URI will be constructed
func determineOperatorUri(version string) (string, error) {
	if version == "latest" {
		uri, err := url.Parse(fmt.Sprintf(operatorReleaseUri, version))
		if err != nil {
			return "", fmt.Errorf("failed to parse uri: %w", err)
		}
		return uri.String(), nil
	}

	// Check for a valid semver version
	re, err := regexp.Compile(constants.SemverRegexWithoutV)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}
	if re.MatchString(version) {
		// We'll just add the v in this case and handle it with the same code as below
		version = fmt.Sprintf("v%s", version)
	}
	rev, err := regexp.Compile(constants.SemverRegexWithV)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}
	if rev.MatchString(version) {
		uri, err := url.Parse(fmt.Sprintf(operatorReleaseUri, version))
		if err != nil {
			return "", fmt.Errorf("failed to parse uri: %w", err)
		}
		return uri.String(), nil
	}
	log.Debug().Msg("Version is not a valid semver version, assuming it is a URI")

	uri, err := url.ParseRequestURI(version)
	if err == nil {
		return uri.String(), nil
	}

	log.Debug().Msg("Version is not a valid URI")

	return "", fmt.Errorf("version is not a valid semver version or URI")
}