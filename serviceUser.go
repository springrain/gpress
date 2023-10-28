package main

import (
	"context"

	"gitee.com/chunanyong/zorm"
)

func insertUser(ctx context.Context, user User) error {
	// 清空用户,只能有一个管理员
	deleteAll(ctx, tableUserName)
	// 初始化数据
	user.Id = "gpress_admin"
	user.SortNo = 1
	user.Status = 1
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.Insert(ctx, &user)
	})
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
	if len(userMap) < 1 || userMap["id"] == nil || userMap["chainType"] == nil || userMap["chainAddress"] == nil { //没有数据
		return "", "", "", err
	}
	return userMap["id"].(string), userMap["chainType"].(string), userMap["chainAddress"].(string), err
}
