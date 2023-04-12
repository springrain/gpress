package main

import (
	"context"

	"gitee.com/chunanyong/zorm"
)

func insertUser(ctx context.Context, account string, password string) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, tableUserName)
	// 初始化数据
	user := zorm.NewEntityMap(tableUserName)
	id := FuncGenerateStringID()
	user.PkColumnName = "id"
	user.Set("id", id)
	user.Set("account", account)
	user.Set("password", password)
	user.Set("userName", account)
	_, err := saveEntityMap(ctx, user)
	return err
}

func findUserId(ctx context.Context, account string, password string) (string, error) {
	finder := zorm.NewSelectFinder(tableUserName, "id").Append(" WHERE account=? and password=?", account, password)
	userId := ""
	_, err := zorm.QueryRow(ctx, finder, &userId)
	return userId, err
}
