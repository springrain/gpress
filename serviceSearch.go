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
	"strings"

	"gitee.com/chunanyong/zorm"
)

// findSearchContent searches public content and prepares the API response.
func findSearchContent(ctx context.Context, q string, pageNo int) (ResponseData, error) {
	responseData := ResponseData{StatusCode: 0}
	q = strings.TrimSpace(q)
	if pageNo < 1 {
		pageNo = 1
	}

	page := zorm.NewPage()
	page.PageNo = pageNo
	page.PageSize = 20
	finder := zorm.NewFinder().Append("SELECT id,title,href_url,summary,create_time,update_time,tag,category_name,description FROM content WHERE status in (1,2) ")
	if q != "" {
		if dbDaoConfig.Dialect == "postgresql" {
			finder.Append("AND markdown ||| ? ", q)
		} else {
			finder.Append("AND rowid in (select rowid from fts_content where fts_content match jieba_query(?) ) ", q)
		}
	}
	finder.Append("ORDER BY status desc, sortno desc")

	contents := make([]Content, 0)
	if err := zorm.Query(ctx, finder, &contents, page); err != nil {
		return responseData, err
	}
	for i := range contents {
		if contents[i].HrefURL == "" {
			contents[i].HrefURL = funcBasePath() + funcTrimPrefixSlash(contents[i].Id)
		} else {
			contents[i].HrefURL, _ = funcHrefURL(contents[i].HrefURL)
		}
	}

	responseData.StatusCode = 1
	responseData.Data = contents
	responseData.Page = page
	return responseData, nil
}
