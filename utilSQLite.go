// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"database/sql"
	"os"
	"runtime"

	"gitee.com/chunanyong/zorm"

	// 00.引入数据库驱动
	"github.com/mattn/go-sqlite3"
)

var dbDao *zorm.DBDao

var dbDaoConfig = zorm.DataSourceConfig{
	DSN:                   sqliteDBfile,
	DriverName:            "sqlite3_simple", // 使用simple分词器会注册这个驱动名
	Dialect:               "sqlite",
	MaxOpenConns:          50,
	MaxIdleConns:          50,
	ConnMaxLifetimeSecond: 600,
	SlowSQLMillis:         -1,
}

// 全局存放 表对象,启动之后,所有的表都通过这个map获取,一个表只能打开一次,类似数据库连接,用一个对象操作
//var TableMap map[string]bleve.Table = make(map[string]bleve.Table)

//var TableMap sync.Map

// 初始化 sqlite数据库
func checkSQLiteStatus() bool {

	defaultFtsFile := datadir + "fts5/libsimple"

	//CPU架构
	goarch := runtime.GOARCH

	ftsFile := defaultFtsFile + "-" + goarch
	if !pathExist(ftsFile) { //文件不存在,使用默认的地址
		ftsFile = defaultFtsFile
	}
	//注册fts5的simple分词器,建议使用jieba分词
	//需要  --tags "fts5"
	sql.Register("sqlite3_simple", &sqlite3.SQLiteDriver{
		Extensions: []string{
			ftsFile, //不要加后缀,它会自己处理,这样代码也统一
		},
	})

	var err error
	dbDao, err = zorm.NewDBDao(&dbDaoConfig)
	if dbDao == nil || err != nil { //数据库初始化失败
		return false
	}

	//初始化结巴分词的字典
	finder := zorm.NewFinder().Append("SELECT jieba_dict(?)", datadir+"dict")
	fts5jieba := ""
	_, err = zorm.QueryRow(context.Background(), finder, &fts5jieba)
	if err != nil {
		return false
	}

	finder = zorm.NewFinder().Append("select jieba_query(?)", "让数据自由一点点,让世界美好一点点")
	_, err = zorm.QueryRow(context.Background(), finder, &fts5jieba)
	if err != nil {
		return false
	}
	//fmt.Println(fts5jieba)
	isInit := pathExist(datadir + "gpress.db")
	if !isInit { //需要初始化数据库
		return isInit
	}

	if tableExist(tableContentName) {
		return true
	}

	sqlByte, err := os.ReadFile("gpressdatadir/gpress.sql")
	if err != nil {
		panic(err)
	}
	createTableSQL := string(sqlByte)
	if createTableSQL == "" {
		panic("gpressdatadir/gpress.sql 文件异常")
	}

	ctx := context.Background()
	_, err = execNativeSQL(ctx, createTableSQL)
	if err != nil {
		panic(err)
	}

	return true
}

func tableExist(tableName string) bool {
	finder := zorm.NewSelectFinder("sqlite_master", "count(*)").Append("WHERE type=? and name=?", "table", tableName)
	count := 0
	zorm.QueryRow(context.Background(), finder, &count)
	return count > 0
}

/*
// 保存新表
func saveEntityMap(ctx context.Context, newTable zorm.IEntityMap) (ResponseData, error) {

	responseData := ResponseData{StatusCode: 1}

	id := ""
	newId, ok := newTable.GetDBFieldMap()["id"]
	if ok {
		id = newId.(string)
	}
	if id == "" {
		id = FuncGenerateStringID()
	}

	newTable.Set("id", id)

	if newTable.GetDBFieldMap()["sortNo"] == 0 {
		count, _ := selectTableCount(ctx, newTable.GetTableName())
		newTable.Set("sortNo", count)
	}

	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
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
*/

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

func execNativeSQL(ctx context.Context, nativeSQL string) (bool, error) {
	finder := zorm.NewFinder().Append(nativeSQL)
	finder.InjectionCheck = false
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func selectTableCount(ctx context.Context, tableName string) (int, error) {

	finder := zorm.NewSelectFinder(tableName, "count(*)")
	count := 0
	_, err := zorm.QueryRow(ctx, finder, &count)

	return count, err
}
