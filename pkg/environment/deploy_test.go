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

package environment

//
// func TestGenerate(test *testing.T) {
// 	tmpdir, err := ioutil.TempDir("", "example")
// 	if err != nil {
// 		test.Fatal(err)
// 	}
//
// 	defer os.RemoveAll(tmpdir)
// 	helpers.GitClone("https://github.com/Sabayon/sbi-tasks", tmpdir)
// 	ctx := common.NewContext(tmpdir + ".replicant.db")
// 	ctx.ControlRepoPath = tmpdir
// 	cwd, _ := os.Getwd()
// 	dep := &Deployment{Client: &client.Fetcher{}, Context: ctx}
// 	_, err = dep.Generate("origin/master")
// 	if err != nil {
// 		test.Fatal(err)
// 	}
// 	os.Chdir(tmpdir)
// 	err = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
//
// 		if !f.IsDir() && strings.HasSuffix(f.Name(), ".yaml") {
// 			_, err := state.Find(ctx, "Source", path)
// 			if strings.Contains(path, "monitor-spinbase") {
// 				test.Log(path)
// 				if err != nil {
// 					//test.Fatal(err)
// 					test.Log(err)
// 				}
// 			}
//
// 		}
// 		return nil
// 	})
// 	os.Chdir(cwd)
//
// }
//
// func TestApply(test *testing.T) {
// 	tmpdir, err := ioutil.TempDir("", "example")
// 	if err != nil {
// 		test.Fatal(err)
// 	}
// 	defer os.RemoveAll(tmpdir)
// 	ctx := common.NewContext(tmpdir + ".replicant.db")
// 	ctx.ControlRepoPath = tmpdir
//
// 	helpers.GitClone("https://github.com/Sabayon/sbi-tasks", tmpdir)
// 	helpers.GitAlign("origin/master", tmpdir)
//
// 	//Generate(ctx, "origin/master")
// 	dep := &Deployment{Client: &client.Fetcher{}, Context: ctx}
// 	helpers.Git([]string{"reset", "--hard", "HEAD~40"}, tmpdir)
// 	dep.Apply("origin/master")
// }
