package main

import (
	"context"

	"gitee.com/chunanyong/zorm"
)

func insertUser(ctx context.Context, userMap map[string]string) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, tableUserName)
	// 初始化数据
	user := zorm.NewEntityMap(tableUserName)
	id := FuncGenerateStringID()
	user.PkColumnName = "id"
	user.Set("id", id)
	for k, v := range userMap {
		user.Set(k, v)
	}

	_, err := saveEntityMap(ctx, user)
	return err
}

func findUserId(ctx context.Context, account string, password string) (string, error) {
	finder := zorm.NewSelectFinder(tableUserName, "id").Append(" WHERE account=? and password=?", account, password)
	userId := ""
	_, err := zorm.QueryRow(ctx, finder, &userId)
	return userId, err
}

func findUserAddress(ctx context.Context) (string, string, string, error) {
	finder := zorm.NewSelectFinder(tableUserName, "id,chainType,chainAddress")
	userMap, err := zorm.QueryRowMap(ctx, finder)
	if len(userMap) < 1 { //没有数据
		return "", "", "", err
	}
	return userMap["id"].(string), userMap["chainType"].(string), userMap["chainAddress"].(string), err
}
