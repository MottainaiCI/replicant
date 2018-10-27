/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>

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

package helpers

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// Git executes a git command with the given args as a []string, outputs as a string
func Git(cmdArgs []string, dir string) (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	cmdName := "git"
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Println("There was an error running git command: ", err)
		log.Println(strings.Join(cmdArgs, " "))
		log.Fatalln(string(cmdOut))
	}
	os.Chdir(cwd)
	result := string(cmdOut)
	return strings.TrimSpace(result), err
}

func GitClone(repo, dir string) {
	cwd, _ := os.Getwd()
	log.Println(Git([]string{"clone", repo, dir}, cwd))
}

// GitAlign executes a fetch --all and reset --hard to origin/master on the given git repository
func GitAlign(branch, dir string) {
	if branch == "" {
		branch = "origin/master"
	}
	log.Println(Git([]string{"fetch", "--all"}, dir))
	log.Println(Git([]string{"reset", "--hard", branch}, dir))
}

func GitPrevCommit(dir string) (string, error) {
	result, err := Git([]string{"log", "-2", `--pretty=format:"%h"`}, dir)
	temp := strings.Split(result, "\n")
	return temp[1], err
}

func GitHash(revision, dir string) string {
	head, _ := Git([]string{"rev-parse", revision}, dir)
	return head
}

func GitFileDiff(revision, dir string) string {
	head, _ := Git([]string{"diff-tree", "--name-status", "-r", "--no-commit-id", revision}, dir)
	return head
}

// GitHead returns the Head of the given repository
func GitHead(dir string) string {
	return GitHash("HEAD", dir)
}
