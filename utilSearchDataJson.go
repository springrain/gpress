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
	"encoding/json"
	"os"
	"sync"

	"gitee.com/chunanyong/zorm"
)

// onlyOnce控制并发
// var onlyOnce = make(chan struct{}, 1)
var searchDataLock = &sync.Mutex{}

func genSearchDataJson() error {
	//onlyOnce <- struct{}{}
	//defer func() { <-onlyOnce }()
	searchDataLock.Lock()
	defer searchDataLock.Unlock()

	finder := zorm.NewSelectFinder(tableContentName, "id,title,hrefURL,categoryName,content,description").Append("WHERE status=1 order by sortNo desc")
	page := zorm.NewPage()
	page.PageSize = 2000
	datas := make([]Content, 0)
	err := zorm.Query(context.Background(), finder, &datas, page)
	if err != nil {
		return err
	}
	for i := 0; i < len(datas); i++ {
		if datas[i].HrefURL == "" {
			datas[i].HrefURL = funcBasePath() + "post/" + datas[i].Id
		}
	}
	dataBytes, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	err = os.WriteFile(searchDataJsonFile, dataBytes, os.ModePerm)
	return err
}
