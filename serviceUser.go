package main

import (
	"context"

	"github.com/blevesearch/bleve/v2"
)

func insertUser(ctx context.Context, account string, password string) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, indexUserName)
	userIndex := IndexMap[indexUserName]
	// 初始化数据
	user := make(map[string]string)
	id := FuncGenerateStringID()
	user["id"] = id
	user["account"] = account
	user["password"] = password
	user["userName"] = account

	return userIndex.Index(id, user)
}

func findUserId(ctx context.Context, account string, password string) (string, error) {
	userIndex := IndexMap[indexUserName]

	accountQuery := bleveNewTermQuery(account)
	accountQuery.SetField("account")

	passwordQuery := bleveNewTermQuery(password)
	passwordQuery.SetField("password")
	// 多个条件联查
	query := bleve.NewConjunctionQuery(accountQuery, passwordQuery)
	// 只查一条
	serarchRequest := bleve.NewSearchRequestOptions(query, 1, 0, false)
	// 只查询id
	serarchRequest.Fields = []string{"id"}

	result, err := userIndex.SearchInContext(ctx, serarchRequest)
	if err != nil {
		return "", err
	}

	userId := ""
	if len(result.Hits) > 0 {
		userId = result.Hits[0].Fields["id"].(string)
	}
	return userId, nil
}
