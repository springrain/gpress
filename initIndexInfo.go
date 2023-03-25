package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// IndexInfoStruct 记录所有的表信息(索引名:indexInfo)
type IndexInfoStruct struct {
	// ID 主键 值为 IndexName,也就是表名
	ID string
	// Name 索引名称,类似表名中文说明
	Name string
	// Code 索引代码
	Code string
	// IndexType index/module 索引和模型,两种类型
	IndexType string
	// CreateTime 创建时间
	CreateTime time.Time
	// UpdateTime 更新时间
	UpdateTime time.Time
	// CreateUser  创建人,初始化 system
	CreateUser string
	// SortNo 排序
	SortNo int
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Active int
}

// initIndexInfo 初始化创建indexInfo索引
func initIndexInfo() (bool, error) {
	// 创建配置表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName
	indexInfo, err := bleve.New(indexInfoName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[indexInfoName] = indexInfo
	return true, nil
}
