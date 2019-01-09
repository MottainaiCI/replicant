/*

Copyright (C) 2019  Ettore Di Giacinto <mudler@gentoo.org>

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

package environment

import (
	"os"

	logrus "github.com/sirupsen/logrus"

	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	common "github.com/MottainaiCI/replicant/pkg/common"
	environment "github.com/MottainaiCI/replicant/pkg/environment"
	cobra "github.com/spf13/cobra"
)

func newEnvironmentStatus(config *setting.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status [OPTIONS]",
		Short: "Display environment status",
		Args:  cobra.OnlyValidArgs,
		// TODO: PreRun check of minimal args if --json is not present
		Run: func(cmd *cobra.Command, args []string) {
			//var v *viper.Viper = config.Viper
			repopath, err := cmd.Flags().GetString("environment")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"component": "status",
					"error":     err,
				}).Error("You must specify an environment to deploy ( your git control repo )")
				return
			}
			//	client := client.NewTokenClient(v.GetString("master"), v.GetString("apikey"), config)
			ctx := common.NewContext(repopath + ".replicant.db")
			ctx.ControlRepoPath = repopath
			environment.PrintInfo(repopath)

		},
	}
	cwd, _ := os.Getwd()
	var flags = cmd.Flags()
	flags.StringP("environment", "e", cwd, "Environment control repo path")

	return cmd
}
