package cmd

import (
	"fmt"
	"os"
	"os/exec"

	_ "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
	"github.com/urfave/cli/v2"

	"boundless-cli/version"
)

const BoundlessManifestUrl = "https://raw.githubusercontent.com/mirantis/boundless/main/deploy/static/boundless-operator.yaml"

// App is the main urfave/cli.App for bctl
var App = &cli.App{
	Name:  "bctl",
	Usage: "boundless operator management tool",
	Flags: []cli.Flag{
		debugFlag,
		traceFlag,
		redactFlag,
	},
	Before: actions(initLogging),
	Commands: []*cli.Command{
		{
			Name:  "version",
			Usage: "Output bctl version",
			Action: func(ctx *cli.Context) error {
				fmt.Printf("version: %s\n", version.Version)
				return nil
			},
		},
		cmdInit,
		cmdApply,
		cmdUpdate,
		cmdReset,
	},
	Action: func(c *cli.Context) error {
		return cmdWrapper(c)
	},
}

func cmdWrapper(c *cli.Context) error {
	// Print help if no arguments are given
	if len(c.Args().Slice()) == 0 {
		return cli.ShowAppHelp(c)
	}

	k0sctlConfigFile, err := getK0sctlConfigPath(c)
	if err != nil {
		return err
	}
	args := append(argInsert(c.Command.Name, c.Args().Slice()), "--config", k0sctlConfigFile)
	cmd := exec.Command("k0sctl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func argInsert(arg string, args []string) []string {
	return append([]string{arg}, args...)
}
