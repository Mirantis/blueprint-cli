package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
	"github.com/rs/zerolog/log"
)

var operatorReleaseUri = "https://github.com/MirantisContainers/blueprint/releases/download/%s/blueprint-operator.yaml"

// Github uses a slightly different path for the latest release
var latestOperatorReleaseUri = "https://github.com/mirantiscontainers/blueprint/releases/latest/download/blueprint-operator.yaml"

// determineOperatorUri determines the URI of the operator based on the version
// If the version is a valid URI, it will be returned as is for dev testing
// Otherwise, it will be assumed to be a version and the URI will be constructed
func determineOperatorUri(version string) (string, error) {
	if version == "latest" {
		return latestOperatorReleaseUri, nil
	}

	// Check for a valid semver version
	regexWithoutV, err := regexp.Compile(constants.SemverRegexWithoutV)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}
	regexWithV, err := regexp.Compile(constants.SemverRegexWithV)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}

	if regexWithoutV.MatchString(version) {
		// We'll just add the v in this case and handle it with the same code as below
		version = fmt.Sprintf("v%s", version)
	}
	if regexWithV.MatchString(version) {
		return parseUri(version)
	}
	log.Debug().Msg("Version is not a valid semver version, assuming it is a URI")

	uri, err := url.ParseRequestURI(version)
	if err == nil {
		return uri.String(), nil
	}
	log.Debug().Msg("Version is not a valid URL")

	isFile := strings.HasPrefix(version, "file://")
	if isFile {
		return version, nil
	}
	log.Debug().Msg("Version is not a valid file URI")

	return "", fmt.Errorf("version is not a valid semver version or URI")
}

// parseUri parses the URI for the operator release
func parseUri(version string) (string, error) {
	uri, err := url.Parse(fmt.Sprintf(operatorReleaseUri, version))
	if err != nil {
		return "", fmt.Errorf("failed to parse uri: %w", err)
	}
	return uri.String(), nil
}

// setImageRegistry replaces the image registry in the BOP manifest with the provided one
// if the supplied URI points to remote file, it is downloaded first;
// an updated manifest is saved to a temporary file and its path is returned;
// if the image registry is not provided or is the default one, the original URI is returned;
// the second return value indicates whether the temporary file was created and should be removed later
func setImageRegistry(bopURI, imageRegistry string) (string, bool, error) {
	if bopURI == "" {
		return "", false, fmt.Errorf("empty BOP manifest URI")
	}
	if imageRegistry == "" || imageRegistry == constants.MirantisImageRegistry {
		return bopURI, false, nil
	}

	var manifestBytes []byte
	var err error
	if strings.HasPrefix(bopURI, "file://") {
		manifestBytes, err = readLocalManifest(strings.TrimPrefix(bopURI, "file://"))
	} else {
		manifestBytes, err = downloadRemoteManifest(bopURI)
	}
	if err != nil {
		return "", false, fmt.Errorf("unable to obtain BOP manifest: %w", err)
	}

	manifestBytes = bytes.ReplaceAll(manifestBytes, []byte(constants.MirantisImageRegistry), []byte(imageRegistry))

	tmpManifest, err := os.CreateTemp("", "bop-*.yaml")
	if err != nil {
		return "", false, fmt.Errorf("unable to create temporary manifest file: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tmpManifest.Close()
			_ = os.Remove(tmpManifest.Name())
		}
	}()

	writtenN, err := tmpManifest.Write(manifestBytes)
	if err != nil {
		return "", false, fmt.Errorf("unable to write temporary manifest file with updated image registry: %w", err)
	}
	if writtenN != len(manifestBytes) {
		err = fmt.Errorf("unable to write temporary manifest file with updated image registry: wrote %d bytes, expected %d", writtenN, len(manifestBytes))
		return "", false, err
	}
	if err = tmpManifest.Close(); err != nil {
		return "", false, fmt.Errorf("unable to close temporary manifest file: %w", err)
	}

	return fmt.Sprintf("file://%s", tmpManifest.Name()), true, nil
}

func readLocalManifest(bopPath string) ([]byte, error) {
	f, err := os.Open(bopPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open BOP manifest file %s: %w", bopPath, err)
	}
	defer f.Close()

	manifestBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read BOP manifest file %s: %w", bopPath, err)
	}

	return manifestBytes, nil
}

func downloadRemoteManifest(bopURI string) ([]byte, error) {
	resp, err := http.Get(bopURI)
	if err != nil {
		return nil, fmt.Errorf("unable to download BOP manifest from %s: %w", bopURI, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got an unexpected status code while downloading BOP manifest from %s: %s", bopURI, resp.Status)
	}

	manifestBytes := new(bytes.Buffer)

	_, err = io.Copy(manifestBytes, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read BOP manifest downloaded from %s: %w", bopURI, err)
	}

	return manifestBytes.Bytes(), nil
}
