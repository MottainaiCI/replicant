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

package state

import (
	common "github.com/MottainaiCI/replicant/pkg/common"
	"github.com/asdine/storm"
)

type State struct {
	ID     int    `storm:"increment"`
	PlanID string `storm:"unique"`
	TaskID string `storm:"unique"`
	Source string `storm:"unique"`
}

func (s *State) Save(ctx *common.Context) error {
	return ctx.DBOpen(func(db *storm.DB) error {
		return db.Save(s)
	})
}

func (s *State) Delete(ctx *common.Context) error {
	return ctx.DBOpen(func(db *storm.DB) error {
		return db.DeleteStruct(s)
	})
}
func Find(ctx *common.Context, field, value string) (*State, error) {
	var state State
	err := ctx.DBOpen(func(db *storm.DB) error {
		return db.One(field, value, &state)
	})
	return &state, err
}
