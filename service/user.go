package service

import (
	"context"
	"gitee.com/gpress/gpress/configs"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
)

func insertUser(ctx context.Context, account string, password string) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, configs.USER_INDEX_NAME)
	userIndex := configs.IndexMap[configs.USER_INDEX_NAME]
	// 初始化数据
	user := make(map[string]string)
	id := util.FuncGenerateStringID()
	user["id"] = id
	user["account"] = account
	user["password"] = password
	user["userName"] = account

	return userIndex.Index(id, user)
}

func findUserId(ctx context.Context, account string, password string) (string, error) {
	userIndex := configs.IndexMap[configs.USER_INDEX_NAME]

	accountQuery := bleve.NewTermQuery(account)
	accountQuery.SetField("account")

	passwordQuery := bleve.NewTermQuery(password)
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
