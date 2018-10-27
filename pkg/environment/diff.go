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
	"bufio"
	"strings"
)

type Diff struct {
	Path string
	Type string
}

func DiffFromOutput(s string) []Diff {
	diffs := []Diff{}

	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	var first = true
	var last_type string
	for scanner.Scan() {
		m := scanner.Text()
		if first {
			last_type = m
			first = false
		} else {
			diffs = append(diffs, Diff{Type: last_type, Path: m})
			first = true
		}
	}

	return diffs
}

func (d *Diff) IsAdd() bool {
	if d.Type == "A" {
		return true
	}
	return false
}
func (d *Diff) IsDeleted() bool {
	if d.Type == "D" {
		return true
	}
	return false
}

func (d *Diff) IsModified() bool {
	if d.Type == "M" {
		return true
	}
	return false
}
