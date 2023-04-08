package main

import (
	"time"
)

// 默认模型 indexInfo indexType="module". 只是记录,并不创建index,全部保存到context里,用于全局检索
// 对比content新增的字段
func init() {

	// 获取当前时间
	now := time.Now()

	//保存表信息
	indexInfo, _, _ := openBleveIndex(indexInfoName)
	indexInfo.Index(indexModuleDefaultName, IndexInfoStruct{
		ID:         indexModuleDefaultName,
		Name:       "默认模型",
		Code:       indexModuleDefaultName,
		IndexType:  "module",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     4,
		Status:     1,
	})
}
