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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	task "github.com/MottainaiCI/mottainai-server/pkg/tasks"
)

func PrintInfo(dirPath string) error {

	fullPath, err := filepath.Abs(dirPath)

	if err != nil {
		return err
	}

	callback := func(path string, fi os.FileInfo, err error) error {
		return checkPlan(fullPath, path, fi, err)
	}

	return filepath.Walk(fullPath, callback)
}

func checkPlan(root string, path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		return nil
	}
	var plan *task.Plan

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		plan, err = task.PlanFromYaml(path)
	} else if strings.HasSuffix(path, ".json") {
		plan, err = task.PlanFromJSON(path)
	}
	if err != nil {
		return err
	}
	if plan == nil {
		return nil
	}
	var kind string
	if plan.Planned == "" {
		kind = "task"
	} else {
		kind = "plan"
	}

	fmt.Println("["+kind+"]", rel)

	return nil
}
