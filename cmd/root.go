/*

Copyright (C) 2018  Ettore Di Giacinto <mudler@gentoo.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package cmd

import (
	"fmt"
	"os"

	//	"reflect"
	clicommon "github.com/MottainaiCI/mottainai-cli/common"
	logger "github.com/MottainaiCI/mottainai-server/pkg/logging"
	common "github.com/MottainaiCI/replicant/pkg/common"

	environment "github.com/MottainaiCI/replicant/cmd/environment"
	validate "github.com/MottainaiCI/replicant/cmd/validate"

	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"

	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
)

const (
	cliName = `Mottainai Replicant
Copyright (c) 2017-2021 Mottainai

Command line interface for Mottainai replicant`

	cliExamples = `$> replicant -m http://127.0.0.1:8080 environment apply --revision origin/master

$> replicant -m http://127.0.0.1:8080 environment deploy --revision origin/master
`
	version = setting.MOTTAINAI_VERSION + ".1"
)

func initConfig(config *setting.Config) {
	// Set env variable
	config.Viper.SetEnvPrefix(common.MREPLICANT_ENV_PREFIX)
	config.Viper.BindEnv("config")
	config.Viper.SetDefault("master", "http://localhost:8080")
	config.Viper.SetDefault("profile", "")
	config.Viper.SetDefault("config", "")
	config.Viper.SetDefault("etcd-config", false)

	config.Viper.AutomaticEnv()

	// Set config file name (without extension)
	config.Viper.SetConfigName(common.MREPLICANT_CONFIG_NAME)

	// Set Config paths list
	config.Viper.AddConfigPath(common.MREPLICANT_LOCAL_PATH)
	config.Viper.AddConfigPath(fmt.Sprintf("$HOME/%s", common.MREPLICANT_HOME_PATH))

	config.Viper.SetTypeByDefaultValue(true)
}

var Logger *logger.Logger

func initCommand(rootCmd *cobra.Command, config *setting.Config) {
	var pflags = rootCmd.PersistentFlags()
	v := config.Viper

	pflags.StringP("master", "m", "http://localhost:8080", "MottainaiCI webUI URL")
	pflags.StringP("apikey", "k", "fb4h3bhgv4421355", "Mottainai API key")

	pflags.StringP("profile", "p", "", "Use specific profile for call API.")

	v.BindPFlag("master", rootCmd.PersistentFlags().Lookup("master"))
	v.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	v.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	Logger = logger.New()
	Logger.SetupWithConfig(false, config)
	rootCmd.AddCommand(
		environment.NewEnvironmentCommand(config),
		validate.NewValidateCommand(config),
	)
}

func Execute() {
	// Create Main Instance Config object
	var config *setting.Config = setting.NewConfig(nil)

	initConfig(config)

	var rootCmd = &cobra.Command{
		Short:        cliName,
		Version:      version,
		Example:      cliExamples,
		Args:         cobra.OnlyValidArgs,
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			var v *viper.Viper = config.Viper

			// Parse configuration file
			err = config.Unmarshal()
			// TODO: Add loglevel in debug that said no config file processed.
			// if err != nil {
			//	fmt.Println(err)
			//}

			// Load profile data and override master if not present.
			if v.Get("profiles") != nil && !cmd.Flag("master").Changed {

				// PRE: profiles contains a map
				//      map[
				//        <NAME_PROFILE1>:<PROFILE INTERFACE>
				//        <NAME_PROFILE2>:<PROFILE INTERFACE>
				//     ]

				var conf clicommon.ProfileConf
				var profile *clicommon.Profile
				if err = v.Unmarshal(&conf); err != nil {
					fmt.Println("Ignore config: ", err)
				} else {
					if v.GetString("profile") != "" {
						profile, err = conf.GetProfile(v.GetString("profile"))

						if profile != nil {
							v.Set("master", profile.GetMaster())
							if profile.GetApiKey() != "" && !cmd.Flag("apikey").Changed {
								v.Set("apikey", profile.GetApiKey())
							}
						} else {
							fmt.Printf("No profile with name %s. I use default value.\n", v.GetString("profile"))
						}
					}
				}

			}
		},
	}

	initCommand(rootCmd, config)

	// Start command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
