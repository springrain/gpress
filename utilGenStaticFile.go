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
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"gitee.com/chunanyong/zorm"
	"golang.org/x/crypto/sha3"
)

// onlyOnce控制并发
// var onlyOnce = make(chan struct{}, 1)
var searchDataLock = &sync.Mutex{}

// genSearchDataJson 生成flexSearch需要的json文件,默认2000条数据
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

// genStaticHtmlFile 生成全站静态文件
func genStaticHtmlFile() error {
	ctx := context.Background()
	postIds := make([]string, 0)
	f_post := zorm.NewSelectFinder(tableContentName, "id").Append(" WHERE status=1 order by sortNo desc")
	err := zorm.Query(ctx, f_post, &postIds, nil)
	if err != nil {
		return err
	}
	//删除整个目录
	os.RemoveAll(staticHtmlDir)
	//生成文章的文件
	for i := 0; i < len(postIds); i++ {
		postId := postIds[i]
		postURL := httpServerPath + "post/" + postId
		fileHash, err := writeStaticHtml(postURL, staticHtmlDir+"post/"+postId)
		if fileHash == "" || err != nil {
			continue
		}
	}

	return nil
}

// writeStaticHtml 写入静态html
func writeStaticHtml(httpurl string, filePath string) (string, error) {

	response, err := http.Get(httpurl)
	if err != nil {
		FuncLogError(err)
		return "", err
	}
	// 读取资源数据 body: []byte
	body, err := io.ReadAll(response.Body)
	// 关闭资源流
	response.Body.Close()
	if err != nil {
		FuncLogError(err)
		return "", err
	}
	//计算hash
	bytehex := sha3.Sum256(body)
	fileHash := hex.EncodeToString(bytehex[:])
	// 写入文件
	os.MkdirAll(filePath, os.ModePerm)
	err = os.WriteFile(filePath+"/index.html", body, os.ModePerm)
	return fileHash, err
}
