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
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gitee.com/chunanyong/zorm"

	// 00.引入数据库驱动
	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

var dbDao *zorm.DBDao

var dbDaoConfig zorm.DataSourceConfig

// checkDBStatus 初始化数据库,并检查是否成功
func checkDBStatus() bool {

	// 打开文件
	jsonFile, err := os.Open(datadir + "db.json")
	if err == nil {
		// 关闭文件
		defer jsonFile.Close()
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			FuncLogError(nil, err)
		} else {
			// Decode从输入流读取下一个json编码值并保存在v指向的值里
			err = json.Unmarshal(byteValue, &dbDaoConfig)
			if err != nil {
				FuncLogError(nil, err)
			}
		}

	}

	if dbDaoConfig.Dialect == "postgresql" {
		dbDao, err = zorm.NewDBDao(&dbDaoConfig)
		if dbDao == nil || err != nil { //数据库初始化失败
			return false
		}
		extname := ""
		finder := zorm.NewFinder().Append("SELECT extname FROM pg_extension WHERE extname = ?", "pg_search")
		_, err = zorm.QueryRow(context.Background(), finder, &extname)
		if err != nil {
			return false
		}
		// 使用pg_search扩展,进行全文检索
		return extname == "pg_search"
		//return checkPGSQLStatus()
	}
	dbDaoConfig = zorm.DataSourceConfig{
		DSN:                   sqliteDBfile,
		DriverName:            "sqlite3_simple", // 使用simple分词器会注册这个驱动名
		Dialect:               "sqlite",
		MaxOpenConns:          1,
		MaxIdleConns:          1,
		ConnMaxLifetimeSecond: 600,
		SlowSQLMillis:         -1,
	}
	return checkSQLiteStatus()

}

// checkSQLiteStatus 初始化sqlite数据库,并检查是否成功
func checkSQLiteStatus() bool {
	const failSuffix = ".fail"
	if failDB := datadir + "gpress.db" + failSuffix; pathExist(failDB) {
		FuncLogError(nil, fmt.Errorf(funcT("Please confirm if [%s] needs to be manually renamed to [gpress.db]. If not, please manually delete [%s]"), failDB, failDB))
		return false
	}

	fts5File := datadir + "fts5/libsimple"
	//注册fts5的simple分词器,建议使用jieba分词
	//需要  --tags "fts5"
	sql.Register("sqlite3_simple", &sqlite3.SQLiteDriver{
		Extensions: []string{
			fts5File, //不要加后缀,它会自己处理,这样代码也统一
		},
	})

	var err error
	dbDao, err = zorm.NewDBDao(&dbDaoConfig)
	if dbDao == nil || err != nil { //数据库初始化失败
		if db := datadir + "gpress.db"; pathExist(db) {
			_ = os.Rename(db, db+failSuffix)
		}
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
		panic("gpressdatadir/gpress.sql " + funcT("File anomaly"))
	}

	ctx := context.Background()
	_, err = execNativeSQL(ctx, createTableSQL)
	if err != nil {
		panic(err)
	}

	return true
}

// tableExist 数据表是否存在
func tableExist(tableName string) bool {
	finder := zorm.NewSelectFinder("sqlite_master", "count(*)").Append("WHERE type=? and name=?", "table", tableName)
	count := 0
	zorm.QueryRow(context.Background(), finder, &count)
	return count > 0
}

// updateTable 更新表数据
func updateTable(ctx context.Context, newMap zorm.IEntityMap) error {
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateEntityMap(ctx, newMap)
		return nil, err
	})
	return err
}

// deleteById 根据Id删除数据
func deleteById(ctx context.Context, tableName string, id string) error {
	finder := zorm.NewDeleteFinder(tableName).Append(" WHERE id=?", id)
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})

	return err
}

// deleteAll 删除所有数据
func deleteAll(ctx context.Context, tableName string) error {
	finder := zorm.NewDeleteFinder(tableName)
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.UpdateFinder(ctx, finder)
		return nil, err
	})

	return err
}

// execNativeSQL 执行SQL语句
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
