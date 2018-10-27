/*

Copyright (C) 2018 Ettore Di Giacinto <mudler@gentoo.org>

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

package state

import (
	"io/ioutil"
	"os"
	"testing"

	common "github.com/MottainaiCI/replicant/pkg/common"
)

func TestSave(test *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		test.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	state := &State{Source: "M"}
	err = state.Save(&common.Context{Database: tmpfile.Name()})
	if err != nil {
		test.Fatal(err)
	}
}

func TestFind(test *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		test.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())
	ctx := common.NewContext(tmpfile.Name())
	state := &State{Source: "Test"}
	err = state.Save(ctx)
	if err != nil {
		test.Fatal(err)
	}

	state, err = Find(ctx, "Source", "Test")
	if err != nil {
		test.Fatal(err)
	}
	if state.Source != "Test" {
		test.Fatal("Source differs")
	}

}

func TestDelete(test *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		test.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())
	ctx := common.NewContext(tmpfile.Name())
	state := &State{Source: "Test"}
	err = state.Save(ctx)
	if err != nil {
		test.Fatal(err)
	}

	state, err = Find(ctx, "Source", "Test")
	if err != nil {
		test.Fatal(err)
	}
	if state.Source != "Test" {
		test.Fatal("Source differs")
	}
	err = state.Delete(ctx)
	if err != nil {
		test.Fatal(err)
	}

	state, err = Find(ctx, "Source", "Test")
	if err.Error() != "not found" {
		test.Fatal(err)
	}
}
