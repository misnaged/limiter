package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"limiter/config"
	"limiter/internal"
)

func Cmd(app *internal.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "limiter",
		Short:            "limiter",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, app.Config())
		},
	}

	return cmd
}

func initializeConfig(cmd *cobra.Command, cfg *config.Scheme) error {
	viper.SetConfigName("")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config file: %w", err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	bindFlags(cmd)

	return viper.Unmarshal(cfg)
}

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
