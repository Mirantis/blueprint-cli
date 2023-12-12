package utils

import (
	"github.com/a8m/envsubst"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/rs/zerolog/log"
)

const DefaultBlueprintPath = "blueprint.yaml"

func LoadBlueprint(path string) (types.Blueprint, error) {
	if path == "" {
		path = DefaultBlueprintPath
	}

	content, err := ReadFile(path)
	if err != nil {
		return types.Blueprint{}, err
	}

	subst, err := envsubst.Bytes(content)
	if err != nil {
		return types.Blueprint{}, err
	}

	log.Debug().Msgf("Loaded configuration:\n%s", subst)
	cfg, err := types.ParseBoundlessCluster(subst)
	if err != nil {
		return types.Blueprint{}, err
	}

	return cfg, nil
}
