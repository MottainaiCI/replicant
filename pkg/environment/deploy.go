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

package environment

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	logrus "github.com/sirupsen/logrus"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	task "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	common "github.com/MottainaiCI/replicant/pkg/common"
	helpers "github.com/MottainaiCI/replicant/pkg/helpers"
	state "github.com/MottainaiCI/replicant/pkg/state"
)

type Deployment struct {
	Client  *client.Fetcher
	Context *common.Context
}

type PlanHandler func(plan *task.Plan, path string) error

func (d *Deployment) AddPlan(plan *task.Plan, path string) error {
	st, err := state.Find(d.Context, "Source", path)
	if err == nil {
		st.Delete(d.Context)
	}
	st.Source = path

	// Create a plan remotely in such case
	return d.createPlan(plan, st)
}

func (d *Deployment) createPlan(plan *task.Plan, state *state.State) error {
	plan_data := plan.ToMap()
	res, err := d.Client.GenericForm("/api/tasks/plan", plan_data)
	if err != nil {
		return err
	}
	state.PlanID = string(res)
	logrus.WithFields(logrus.Fields{
		"component": "deploy",
		"id":        state.PlanID,
	}).Info("Created")
	return state.Save(d.Context)
}

func (d *Deployment) deletePlan(planID string, state *state.State) error {
	_, err := d.Client.GetOptions("/api/tasks/plan/delete/"+planID, map[string]string{})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"component": "deploy",
		"id":        planID,
	}).Info("Deleted")
	return state.Delete(d.Context)
}

func (d *Deployment) Destroy() {
	var tlist []task.Plan
	d.Client.GetJSONOptions("/api/tasks/planned", map[string]string{}, &tlist)
	for _, i := range tlist {
		_, err := d.Client.GetOptions("/api/tasks/plan/delete/"+i.ID, map[string]string{})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"component": "status",
				"error":     err,
				"id":        i.ID,
			}).Error("Could not delete plan)")
		}
		logrus.WithFields(logrus.Fields{
			"component": "deploy",
			"id":        i.ID,
		}).Info("Deleted")
	}
}

func (d *Deployment) Apply(revision string) (*Environment, error) {
	environment, _ := Find(d.Context, "Revision", revision)
	//if err != nil {
	//	return Generate(d.Context, revision)
	//}
	//helpers.GitAlign(revision, d.Context.ControlRepoPath)
	hash := helpers.GitHash(revision, d.Context.ControlRepoPath)
	current_hash := environment.Hash
	if hash == environment.Hash {
		return environment, nil
	}

	environment.Delete(d.Context)
	environment.Hash = hash
	environment.Revision = revision
	environment.Save(d.Context)
	// Sync changed file and check changed/deleted states, and update/delete plans accordingly
	changed := helpers.GitFileDiff("remotes/"+revision+".."+current_hash, d.Context.ControlRepoPath)

	diffs := DiffFromOutput(changed)
	for _, diff := range diffs {

		var plan *task.Plan
		var err error
		if strings.HasSuffix(diff.Path, ".yaml") || strings.HasSuffix(diff.Path, ".yml") {
			plan, err = task.PlanFromYaml(path.Join(d.Context.ControlRepoPath, diff.Path))
		} else if strings.HasSuffix(diff.Path, ".json") {
			plan, err = task.PlanFromJSON(path.Join(d.Context.ControlRepoPath, diff.Path))
		}

		if err != nil || plan.Planned == "" {
			continue
		}

		if diff.IsAdd() {
			logrus.WithFields(logrus.Fields{
				"component": "deploy",
				"plan":      plan,
				"path":      diff.Path,
			}).Info("Added")
			st, err := state.Find(d.Context, "Source", diff.Path)
			if err == nil {
				// Delete with client also plan in instance
				d.deletePlan(st.PlanID, st)
			}

			newState := &state.State{}
			newState.Source = diff.Path

			// Add it online and store Task/PlanID
			d.createPlan(plan, newState)
		}
		if diff.IsDeleted() {
			logrus.WithFields(logrus.Fields{
				"component": "deploy",
				"path":      diff.Path,
			}).Info("Deleted")
			st, err := state.Find(d.Context, "Source", diff.Path)
			if err == nil {
				// Delete with client also plan in instance
				d.deletePlan(st.PlanID, st)
			}
		}
		if diff.IsModified() {
			logrus.WithFields(logrus.Fields{
				"component": "deploy",
				"path":      diff.Path,
				"plan":      plan,
				"task":      plan.Task,
			}).Info("Modified")
			st, err := state.Find(d.Context, "Source", diff.Path)
			if err == nil {
				// Delete with client also plan in instance
				d.deletePlan(st.PlanID, st)
			}
			newState := &state.State{}
			newState.Source = diff.Path

			// Add it online and store Task/PlanID
			d.createPlan(plan, newState)
		}
	}

	return environment, nil
}

func (d *Deployment) Generate(revision string) (*Environment, error) {
	// First time generation
	env := &Environment{}
	env.Hash = helpers.GitHash(revision, d.Context.ControlRepoPath)
	env.Revision = revision
	env.Save(d.Context)
	_, err := d.Validate()
	if err != nil {
		return nil, err
	}
	// generate state from scratch and add plans
	cwd, _ := os.Getwd() // for having rel paths
	os.Chdir(d.Context.ControlRepoPath)
	err = filepath.Walk(".", d.generateFromPathHandle(d.AddPlan, func(plan *task.Plan, path string) error { return nil }))
	os.Chdir(cwd)

	return env, err
}

func (d *Deployment) Validate() (*Environment, error) {
	// First time generation
	env := &Environment{}
	var failed []string

	// generate state from scratch and add plans
	cwd, _ := os.Getwd() // for having rel paths
	os.Chdir(d.Context.ControlRepoPath)
	err := filepath.Walk(".", d.generateFromPathHandle(
		func(plan *task.Plan, path string) error { return nil },
		func(plan *task.Plan, path string) error {
			failed = append(failed, path)
			return nil
		},
	))
	os.Chdir(cwd)

	for _, v := range failed {
		logrus.WithFields(logrus.Fields{
			"component": "validate",
			"path":      v,
		}).Error("Task/Plan syntax check failed")
	}
	if len(failed) != 0 {
		err = errors.New("Tasks/Plans that failed validation: " + strings.Join(failed, ", "))
	}

	return env, err
}

func (d *Deployment) generateFromPathHandle(h PlanHandler, errorHandler PlanHandler) func(string, os.FileInfo, error) error {
	return func(path string, f os.FileInfo, err error) error {

		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
			plan, err := task.PlanFromYaml(path)
			if err != nil {
				return errorHandler(plan, path)
			} else if plan.Planned != "" {
				logrus.WithFields(logrus.Fields{
					"component": "deploy",
					"path":      path,
				}).Info("Plan found")
				return h(plan, path)
			}
		}
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			plan, err := task.PlanFromJSON(path)
			if err != nil {
				return errorHandler(plan, path)
			} else if plan.Planned != "" {
				logrus.WithFields(logrus.Fields{
					"component": "deploy",
					"path":      path,
				}).Info("Plan found")
				return h(plan, path)
			}
		}
		return nil
	}
}
