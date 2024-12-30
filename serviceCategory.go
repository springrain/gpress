// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"

	"gitee.com/chunanyong/zorm"
)

// findAllCategory 查找所有的导航菜单
func findAllCategory(ctx context.Context) ([]Category, error) {
	finder := zorm.NewSelectFinder(tableCategoryName)
	categorys := make([]Category, 0)
	err := zorm.Query(ctx, finder, &categorys, nil)
	return categorys, err
}