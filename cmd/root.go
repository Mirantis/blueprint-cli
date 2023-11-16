package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/mattn/go-colorable"

	"boundless-cli/internal/k8s"
	"boundless-cli/internal/types"
	"boundless-cli/internal/utils"
)

const (
	appName      = "bctl"
	shortAppDesc = "A tool to manage boundless operator."
)

var (
	pFlags    *PersistenceFlags
	kubeFlags *genericclioptions.ConfigFlags

	blueprint  types.Blueprint
	kubeConfig *k8s.KubeConfig

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogger()
		},
		Run:          runHelp,
		SilenceUsage: true,
	}

	out = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(
		versionCmd(),
		initCmd(),
		applyCmd(),
		updateCmd(),
		resetCmd(),
	)

	pFlags = NewPersistenceFlags()
	rootCmd.PersistentFlags().BoolVarP(&pFlags.Debug, "debug", "d", false, " Enable debug logging (default: false)")

	// TODO (ranyodh): Add support for the other k0sctl commands
}

// Execute root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err)
	}
}

func loadBlueprint(cmd *cobra.Command, args []string) error {
	// load the blueprint file
	var err error
	if blueprint, err = utils.LoadBlueprint(blueprintConfig); err != nil {
		return err
	}
	return nil
}

func loadKubeConfig(cmd *cobra.Command, args []string) error {
	// TODO (ranyodh): check if kubeconfig file is present
	// TODO (ranyodh): if multiple contexts are present, ensure we load the one that is created by bctl
	kubeConfig = k8s.NewConfig(kubeFlags)
	return nil
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}})

	if pFlags.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func addConfigFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&blueprintConfig, "config", "", DefaultConfigFilename, "Path to the blueprint file")
}

func addKubeFlags(flags *pflag.FlagSet) {
	kubeFlags = genericclioptions.NewConfigFlags(k8s.UsePersistentConfig)

	flags.StringVar(kubeFlags.KubeConfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	flags.StringVar(kubeFlags.Timeout, "request-timeout", "", "The length of time to wait before giving up on a single server request")
	flags.StringVar(kubeFlags.Context, "context", "", "The name of the kubeconfig context to use")
	flags.StringVar(kubeFlags.ClusterName, "cluster", "", "The name of the kubeconfig cluster to use")
	flags.StringVar(kubeFlags.AuthInfoName, "user", "", "The name of the kubeconfig user to use")

	// as pFlags
	flags.StringVar(kubeFlags.Impersonate, "as", "", "Username to impersonate for the operation")
	flags.StringArrayVar(kubeFlags.ImpersonateGroup, "as-group", []string{}, "Group to impersonate for the operation")

	// cert pFlags
	flags.BoolVar(kubeFlags.Insecure, "insecure-skip-tls-verify", false, "If true, the server's caCertFile will not be checked for validity")
	flags.StringVar(kubeFlags.CAFile, "certificate-authority", "", "Path to a cert file for the certificate authority")
	flags.StringVar(kubeFlags.KeyFile, "client-key", "", "Path to a client key file for TLS")
	flags.StringVar(kubeFlags.CertFile, "client-certificate", "", "Path to a client certificate file for TLS")

	flags.StringVar(kubeFlags.BearerToken, "token", "", "Bearer token for authentication to the API server")
}
