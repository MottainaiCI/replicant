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

package validate

import (
	"os"
	"path"

	logrus "github.com/sirupsen/logrus"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	common "github.com/MottainaiCI/replicant/pkg/common"
	environment "github.com/MottainaiCI/replicant/pkg/environment"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newEnvironmentValidation(config *setting.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "repo [OPTIONS]",
		Short: "Validates a task control repository",
		Args:  cobra.OnlyValidArgs,
		// TODO: PreRun check of minimal args if --json is not present
		Run: func(cmd *cobra.Command, args []string) {
			var v *viper.Viper = config.Viper
			repopath, err := cmd.Flags().GetString("environment")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"component": "validate",
					"error":     err,
				}).Error("You must specify an environment to validate ( your git control repo )")
				return
			}
			client := client.NewTokenClient(v.GetString("master"), v.GetString("apikey"), config)
			ctx := common.NewContext(path.Join(repopath, ".replicant.db"))
			ctx.ControlRepoPath = repopath

			dep := &environment.Deployment{Client: client, Context: ctx}
			_, err = dep.Validate()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"component": "deploy",
					"error":     err,
				}).Fatal("Error while generating deployment for the supplied environment")
			}
		},
	}
	cwd, _ := os.Getwd()
	var flags = cmd.Flags()
	flags.StringP("environment", "e", cwd, "Environment control repo path")

	return cmd
}
