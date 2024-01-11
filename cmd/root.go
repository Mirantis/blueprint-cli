package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mirantiscontainers/boundless-cli/internal/boundless"
	"github.com/mirantiscontainers/boundless-cli/internal/distro"
	"github.com/mirantiscontainers/boundless-cli/internal/k8s"
	"github.com/mirantiscontainers/boundless-cli/internal/types"
	"github.com/mirantiscontainers/boundless-cli/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	appName      = "bctl"
	shortAppDesc = "A tool to manage boundless operator."
)

var (
	pFlags        *PersistenceFlags
	blueprintFlag string
	operatorUri   string

	blueprint  types.Blueprint
	kubeConfig *k8s.KubeConfig

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogger()
		},
		RunE:         runHelp,
		SilenceUsage: true,
	}

	kubeFlags = genericclioptions.NewConfigFlags(k8s.UsePersistentConfig)
	out       = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(
		versionCmd(),
		initCmd(),
		applyCmd(),
		updateCmd(),
		resetCmd(),
		upgradeCmd(),
		statusCmd(),
	)

	pFlags = NewPersistenceFlags()
	rootCmd.PersistentFlags().StringVarP(&pFlags.LogLevel, "logLevel", "l", DefaultLogLevel, "Specify a log level (info, warn, debug, trace, error)")

	// TODO (ranyodh): Add support for the other k0sctl commands
}

// Execute root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err)
	}
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func loadBlueprint(cmd *cobra.Command, args []string) error {
	var err error
	log.Debug().Msgf("Loading blueprint from %q", blueprintFlag)
	if blueprint, err = utils.LoadBlueprint(blueprintFlag); err != nil {
		return fmt.Errorf("failed to load blueprint file at %q: %w", blueprintFlag, err)
	}
	return nil
}

// loadKubeConfig loads the kubeconfig file
// This function should be added as a pre-run hook for all commands that connects to the cluster
func loadKubeConfig(cmd *cobra.Command, args []string) error {
	// unless context flag is passed, explicitly set the context to use for kubeconfig
	if kubeFlags.Context == nil || *kubeFlags.Context == "" {

		// Determine the distro
		provider, err := distro.GetProvider(&blueprint, kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to determine kubernetes provider: %w", err)
		}
		context := provider.GetKubeConfigContext()
		kubeFlags.Context = &context
	}
	kubeConfig = k8s.NewConfig(kubeFlags)

	// TODO (ranyodh): remove this hack
	// This is a hack to ensure that the kubeconfig file is not loaded for apply command
	// because the cluster is not yet created at this point
	if cmd.Name() == "apply" {
		return nil
	}

	log.Debug().Msgf("Loading kubeconfig from %q", kubeConfig.GetConfigPath())
	// Try to load kubeconfig file here, and fail early if it is not present
	if err := kubeConfig.TryLoad(); err != nil {
		return err
	}
	return nil
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}})
	zerolog.SetGlobalLevel(parseLevel(pFlags.LogLevel))
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func addOperatorUriFlag(flags *pflag.FlagSet) {
	flags.StringVarP(&operatorUri, "operator-uri", "", boundless.ManifestUrlLatest, "URL or path to the Boundless Operator manifest file")
}

func addBlueprintFileFlags(flags *pflag.FlagSet) {
	// @todo ranyodh: remove deprecated `config` flag before 1.0.0
	flags.StringVarP(&blueprintFlag, "config", "c", DefaultBlueprintFileName, "Path to the blueprint file")
	_ = flags.MarkDeprecated("config", "use --file (or -f)")

	flags.StringVarP(&blueprintFlag, "file", "f", DefaultBlueprintFileName, "Path to the blueprint file")
}

func addKubeFlags(flags *pflag.FlagSet) {
	// Exposing certain flags from k8s.io/cli-runtime/pkg/genericclioptions
	// To expose all flags, use kubeFlags.AddFlags(flags)
	flags.StringVar(kubeFlags.KubeConfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	flags.StringVar(kubeFlags.Timeout, "request-timeout", "", "The length of time to wait before giving up on a single server request")
	flags.StringVar(kubeFlags.Context, "context", "", "The name of the kubeconfig context to use")
	flags.StringVar(kubeFlags.ClusterName, "cluster", "", "The name of the kubeconfig cluster to use")
	flags.StringVar(kubeFlags.AuthInfoName, "user", "", "The name of the kubeconfig user to use")

	// as flags
	flags.StringVar(kubeFlags.Impersonate, "as", "", "Username to impersonate for the operation")
	flags.StringArrayVar(kubeFlags.ImpersonateGroup, "as-group", []string{}, "Group to impersonate for the operation")

	// cert flags
	flags.BoolVar(kubeFlags.Insecure, "insecure-skip-tls-verify", false, "If true, the server's caCertFile will not be checked for validity")
	flags.StringVar(kubeFlags.CAFile, "certificate-authority", "", "Path to a cert file for the certificate authority")
	flags.StringVar(kubeFlags.KeyFile, "client-key", "", "Path to a client key file for TLS")
	flags.StringVar(kubeFlags.CertFile, "client-certificate", "", "Path to a client certificate file for TLS")

	flags.StringVar(kubeFlags.BearerToken, "token", "", "Bearer token for authentication to the API server")
}

func strPtr(s string) *string {
	return &s
}
