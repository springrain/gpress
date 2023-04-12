package main

import (
	"context"

	"gitee.com/chunanyong/zorm"

	// 00.引入数据库驱动
	_ "modernc.org/sqlite"
)

// dbDaoConfig 数据库的配置.这里只是模拟,生产应该是读取配置配置文件,构造DataSourceConfig
var dbDaoConfig = zorm.DataSourceConfig{
	// DSN 数据库的连接字符串,parseTime=true会自动转换为time格式,默认查询出来的是[]byte数组.&loc=Local用于设置时区
	DSN: sqliteDBfile,
	// DriverName 数据库驱动名称:mysql,postgres,oracle(go-ora),sqlserver,sqlite3,go_ibm_db,clickhouse,dm,kingbase,aci,taosSql|taosRestful 和Dialect对应
	// sql.Open(DriverName,DSN) DriverName就是驱动的sql.Open第一个字符串参数,根据驱动实际情况获取
	DriverName: "sqlite",
	// Dialect 数据库方言:mysql,postgresql,oracle,mssql,sqlite,db2,clickhouse,dm,kingbase,shentong,tdengine 和 DriverName 对应
	Dialect: "sqlite",
	// MaxOpenConns 数据库最大连接数 默认50
	MaxOpenConns: 50,
	// MaxIdleConns 数据库最大空闲连接数 默认50
	MaxIdleConns: 50,
	// ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
	ConnMaxLifetimeSecond: 600,
	// SlowSQLMillis 慢sql的时间阈值,单位毫秒.小于0是禁用SQL语句输出;等于0是只输出SQL语句,不计算执行时间;大于0是计算SQL执行时间,并且>=SlowSQLMillis值
	SlowSQLMillis: -1,
}

var dbDao, _ = zorm.NewDBDao(&dbDaoConfig)

// 全局存放 表对象,启动之后,所有的表都通过这个map获取,一个表只能打开一次,类似数据库连接,用一个对象操作
//var TableMap map[string]bleve.Table = make(map[string]bleve.Table)

//var TableMap sync.Map

func tableExist(tableName string) bool {
	finder := zorm.NewSelectFinder("sqlite_master", "count(*)").Append("WHERE type=? and name=?", "table", tableName)
	count := 0
	zorm.QueryRow(context.Background(), finder, &count)
	return count > 0
}

// 初始化 sqlite数据库
func checkSQLliteStatus() bool {
	if dbDao == nil { //数据库初始化失败
		return false
	}
	isInit := pathExist(datadir + "gpress.db")
	if !isInit { //需要初始化数据库
		return isInit
	}

	// 初始化indexField
	_, err := initTableField()
	if err != nil {
		return false
	}
	// 初始化tableInfo
	_, err = initTableInfo()
	if err != nil {
		return false
	}
	// 初始化 config
	ok, err := initConfig()
	if err != nil {
		return false
	}
	return ok
}

// findTableFieldStruct 获取表中符合条件字段,返回Struct对象
// tableName: 表名/表名
func findTableFieldStruct(ctx context.Context, tableName string, required int) ([]TableFieldStruct, error) {
	finder := zorm.NewSelectFinder(tableFieldName, "*").Append(" WHERE tableCode=? ")
	if required != 0 {
		finder.Append(" and required=? ", required)
	}
	finder.Append("order by sortNo asc,id desc", tableName)
	page := zorm.NewPage()
	page.PageNo = 1
	page.PageSize = 1000
	fields := make([]TableFieldStruct, 0)
	err := zorm.Query(ctx, finder, &fields, page)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	return fields, nil
}

// 保存新表
func saveEntityMap(ctx context.Context, newTable zorm.IEntityMap) (ResponseData, error) {
	tableFields, err := findTableFieldStruct(ctx, newTable.GetTableName(), 1)

	responseData := ResponseData{StatusCode: 1}
	if err != nil {
		FuncLogError(err)
		responseData.StatusCode = 303
		responseData.Message = "查询异常"
		return responseData, err
	}
	id := ""
	newId, ok := newTable.GetDBFieldMap()["id"]
	if ok {
		id = newId.(string)
	}
	if id == "" {
		id = FuncGenerateStringID()
	}

	newTable.Set("id", id)

	for _, v := range tableFields {
		tmp := v.FieldCode
		_, ok := newTable.GetDBFieldMap()[tmp]
		if !ok {
			responseData.StatusCode = 401
			responseData.Message = tmp + "不能为空"
			return responseData, err
		}
	}

	if newTable.GetDBFieldMap()["sortNo"] == 0 {
		count, _ := selectTableCount(ctx, newTable.GetTableName())
		newTable.Set("sortNo", count)
	}

	_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.InsertEntityMap(ctx, newTable)
		return nil, err
	})

	if err != nil {
		FuncLogError(err)
		responseData.StatusCode = 304
		responseData.Message = "建立表异常"
		return responseData, err
	}
	responseData.StatusCode = 200
	responseData.Message = "保存成功"
	responseData.Data = id
	return responseData, err
}

func updateTable(ctx context.Context, newMap zorm.IEntityMap) error {
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateEntityMap(ctx, newMap)
		return nil, err
	})
	return err
}
func deleteById(ctx context.Context, tableName string, id string) error {
	finder := zorm.NewDeleteFinder(tableName).Append(" WHERE id=?", id)
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})

	return err
}
func deleteAll(ctx context.Context, tableName string) error {
	finder := zorm.NewDeleteFinder(tableName)
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})

	return err
}

func crateTable(ctx context.Context, createTableSQL string) (bool, error) {
	finder := zorm.NewFinder().Append(createTableSQL)
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func saveTableInfo(ctx context.Context, tableInfoStruct TableInfoStruct) error {
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.Insert(ctx, &tableInfoStruct)
		return nil, err
	})
	return err
}

func saveTableField(ctx context.Context, tableFiledStruct TableFieldStruct) {
	zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.Insert(ctx, &tableFiledStruct)
		return nil, err
	})
}

func selectTableCount(ctx context.Context, tableName string) (int, error) {

	finder := zorm.NewSelectFinder(tableName, "count(*)")
	count := 0
	_, err := zorm.QueryRow(ctx, finder, &count)

	return count, err
}
