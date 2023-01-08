package main

import (
	"context"

	"github.com/blevesearch/bleve/v2"
)

func insertUser(ctx context.Context, account string, password string) error {
	userIndex := IndexMap[userIndexName]

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
	userIndex := IndexMap[userIndexName]

	queryAccount := bleve.NewTermQuery(account)
	queryAccount.SetField("account")

	passwordAccount := bleve.NewTermQuery(password)
	passwordAccount.SetField("password")
	// 多个条件联查
	query := bleve.NewConjunctionQuery(queryAccount, passwordAccount)

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
