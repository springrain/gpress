package main

import (
	"context"

	"gitee.com/chunanyong/zorm"
	"github.com/blevesearch/bleve/v2"
)

func insertUser(ctx context.Context, account string, password string) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, indexUserName)
	// 初始化数据
	user := zorm.NewEntityMap(indexUserName)
	id := FuncGenerateStringID()
	user.PkColumnName = "id"
	user.Set("id", id)
	user.Set("account", account)
	user.Set("password", password)
	user.Set("userName", account)

	return bleveSaveEntityMap(indexUserName, user)
}

func findUserId(ctx context.Context, account string, password string) (string, error) {

	accountQuery := bleveNewTermQuery(account)
	accountQuery.SetField("account")

	passwordQuery := bleveNewTermQuery(password)
	passwordQuery.SetField("password")
	// 多个条件联查
	query := bleve.NewConjunctionQuery(accountQuery, passwordQuery)
	// 只查一条
	searchRequest := bleve.NewSearchRequestOptions(query, 1, 0, false)
	// 只查询id
	searchRequest.Fields = []string{"id"}

	result, err := bleveSearchInContext(ctx, indexUserName, searchRequest)
	if err != nil {
		return "", err
	}

	userId := ""
	if len(result.Hits) > 0 {
		userId = result.Hits[0].Fields["id"].(string)
	}
	return userId, nil
}
