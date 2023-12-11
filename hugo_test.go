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
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gitee.com/chunanyong/zorm"
)

func TestCategory(t *testing.T) {
	deleteAll(context.Background(), tableCategoryName)
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	navs := []string{"About", "Web", "BlockChain", "CloudNative"}

	for i := 0; i < len(navs); i++ {
		nav := navs[i]
		menu := zorm.NewEntityMap(tableCategoryName)
		menu.Set("id", strings.ToLower(nav))
		menu.Set("name", nav)
		menu.Set("status", 1)
		menu.Set("sortNo", i+1)
		menu.Set("createTime", now)
		menu.Set("updateTime", now)
		_, err := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
			_, err := zorm.InsertEntityMap(ctx, menu)
			return nil, err
		})
		//_, err := saveEntityMap(context.Background(), menu)
		if err != nil {
			t.Error(err)
		}
	}

	siteMap := zorm.NewEntityMap(tableSiteName)
	siteMap.PkColumnName = "id"
	siteMap.Set("id", "gpress_site")
	siteMap.Set("title", "jiagou")
	siteMap.Set("name", "架构")
	siteMap.Set("domain", "jiagou.com")
	err := updateTable(context.Background(), siteMap)
	if err != nil {
		t.Error(err)
	}

}

func TestReadmks(t *testing.T) {
	deleteAll(context.Background(), "content")
	files, err := os.ReadDir("D:/post")
	if err != nil {
		t.Error("读取错误")
	}
	lists := make([]zorm.IEntityMap, 0)
	for i, file := range files {
		fileName := file.Name()
		source, err := os.ReadFile("D:/post/" + fileName)
		if err != nil {
			continue
		}
		sortNo := i + 1
		id := fileName[:strings.LastIndex(fileName, ".")]

		fmt.Println(id)
		markdown := strings.TrimSpace(string(source))
		smk := markdown[strings.Index(markdown, "---")+3:]
		start := strings.Index(smk, "---") + 3
		smk = strings.TrimSpace(smk[start:])
		smkRune := []rune(smk)
		end := 100
		if end > len(smkRune) {
			end = len(smkRune)
		}
		summary := string(smkRune[0:end])
		summary = strings.TrimSpace(summary)
		cMap := zorm.NewEntityMap(tableContentName)
		cMap.Set("id", id)
		cMap.Set("summary", summary)
		cMap.Set("markdown", markdown)
		cMap.Set("sortNo", sortNo)
		cMap.Set("status", 1)

		metaData, tocHtml, html, _ := conver2Html([]byte(markdown))
		dateStr := metaData["date"].(string)
		date, _ := time.Parse("2006-01-02T15:04:05+08:00", dateStr)
		categories := metaData["categories"].([]interface{})
		category := slice2string(categories)
		tags := metaData["tags"].([]interface{})
		tag := slice2string(tags)
		cMap.Set("title", metaData["title"])
		cMap.Set("author", metaData["author"])
		cMap.Set("updateTime", date)
		cMap.Set("createTime", date)
		cMap.Set("categoryName", category)
		cMap.Set("categoryID", category)
		cMap.Set("tag", tag)

		cMap.Set("content", html)
		cMap.Set("toc", tocHtml)
		lists = append(lists, cMap)
		//saveNewTable(context.Background(), "content", cMap)

	}

	fmt.Println("------------")

	var temp zorm.IEntityMap            // 定义临时变量,进行数据交换
	for j := 0; j < len(lists)-1; j++ { // 外循环 循环次数
		for i := 0; i < len(lists)-1; i++ { // 内循环 数组遍历
			m1 := lists[i]
			m2 := lists[i+1]
			d1 := m1.GetDBFieldMap()["updateTime"].(time.Time)
			d2 := m2.GetDBFieldMap()["updateTime"].(time.Time)
			if d1.After(d2) {
				temp = lists[i]
				lists[i] = lists[i+1]
				lists[i+1] = temp
			}
		}
	}

	for i := 0; i < len(lists); i++ { // 内循环 数组遍历
		cMap := lists[i]
		cMap.Set("sortNo", i+1)
		date := cMap.GetDBFieldMap()["updateTime"].(time.Time)
		dateStr := date.Format("2006-01-02 15:04:05")
		cMap.Set("updateTime", dateStr)
		cMap.Set("createTime", dateStr)

		_, err := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
			_, err := zorm.InsertEntityMap(ctx, cMap)
			return nil, err
		})

		//_, err := saveEntityMap(context.Background(), cMap)
		if err != nil {
			t.Error(err)
		}
	}

}

func TestAbout(t *testing.T) {
	ctx := context.Background()
	source, err := os.ReadFile("D:/about.md")
	if err != nil {
		t.Error(err)
	}
	markdown := strings.TrimSpace(string(source))
	cMap := zorm.NewEntityMap(tableContentName)
	cMap.Set("id", "about")
	cMap.Set("summary", "jiagou.com")
	cMap.Set("markdown", markdown)
	cMap.Set("sortNo", 100)
	cMap.Set("status", 0)

	_, tocHtml, html, _ := conver2Html([]byte(markdown))
	date := time.Now().Format("2006-01-02 15:04:05")
	cMap.Set("title", "about")
	cMap.Set("author", "springrain")
	cMap.Set("updateTime", date)
	cMap.Set("createTime", date)
	cMap.Set("categoryName", "about")
	cMap.Set("categoryID", "about")
	cMap.Set("content", html)
	cMap.Set("toc", tocHtml)
	cMap.Set("summary", `本站服务器配置:1核CPU,512M内存,20G硬盘,AnolisOS(ANCK).使用hugo和even模板,编译成静态文件,Nginx作为WEB服务器.我所见识过的一切都将消失一空,就如眼泪消逝在雨中......
	不妨大胆一些,大胆一些......`)
	_, err = zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		_, err := zorm.InsertEntityMap(ctx, cMap)
		return nil, err
	})
	//_, err = saveEntityMap(ctx, cMap)
	if err != nil {
		t.Error(err)
	}

	//更新about的hrefURL
	finder := zorm.NewUpdateFinder(tableCategoryName).Append("hrefURL=? WHERE id=?", "post/about", "about")
	_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.UpdateFinder(ctx, finder)
	})
	if err != nil {
		t.Error(err)
	}
}

func slice2string(slice []interface{}) string {
	len := len(slice)
	s := ""
	for i := 0; i < len; i++ {
		str := slice[i].(string)
		s = s + str
		if i+1 < len {
			s = s + ","
		}
	}
	return s
}
