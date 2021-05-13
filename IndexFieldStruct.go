package main

import "time"

// IndexFieldStruct 索引和字段(索引名:IndexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
// 这个可能是唯一的Struct......
type IndexFieldStruct struct {
	// ID 主键
	ID string `analyzer:"keyword"`
	// IndexCode 索引代码,类似表名 User,SiteInfo,PageTemplate,NavMenu,Module,Content
	IndexCode string `analyzer:"keyword"`
	// IndexCode 索引名称,类似表名中文说明
	IndexName string `analyzer:"keyword"`
	//BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
	BusinessID string `analyzer:"keyword"`
	// FieldCode  字段代码
	FieldCode string `analyzer:"keyword"`
	// FieldName  字段中文名称
	FieldName string `analyzer:"keyword"`
	// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
	FieldType int
	// FieldFormat 数据格式,用于日期或者数字
	FieldFormat string `analyzer:"keyword"`
	// DefaultValue  默认值
	DefaultValue string `analyzer:"keyword"`
	// AnalyzerName  分词器名称
	AnalyzerName string `analyzer:"keyword"`
	// CreateTime 创建时间
	CreateTime time.Time
	// CreateTime 更新时间
	UpdateTime time.Time
	// AnalyzerName  创建人,初始化 system
	CreateUser string `analyzer:"keyword"`
	// SortNo 排序
	SortNo int
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Active int
}
