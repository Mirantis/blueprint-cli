package commands

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/mirantiscontainers/boundless-cli/pkg/components"
	"github.com/mirantiscontainers/boundless-cli/pkg/types"
)

// Init initializes a new cluster
func Init(provider string) error {
	if provider == "kind" {
		return components.Encode(types.ConvertToClusterWithKind("blueprint-cluster", components.DefaultComponents))
	}

	// @TODO Include pFlags for k0sctl init
	cmd2 := exec.Command("k0sctl", "init")
	cmd2.Stdin = os.Stdin
	cmd2.Stderr = os.Stderr

	buf := new(bytes.Buffer)
	cmd2.Stdout = buf
	err := cmd2.Run()
	if err != nil {
		return err
	}

	k0sConfig, err := types.ParseK0sCluster(buf.Bytes())
	if err != nil {
		return err
	}

	return components.Encode(types.ConvertToClusterWithK0s(k0sConfig, components.DefaultComponents))
}
