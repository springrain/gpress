package main

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

// 获取菜单树 pid 为0 为一级菜单
var status float64 = 1

func getNavMenu(pid string) (interface{}, error) {
	// PID 跟 Status 为查询字段
	queryPID := bleveNewTermQuery(pid)
	queryPID.SetField("PID")
	queryActive := bleve.NewNumericRangeInclusiveQuery(&status, &status, &inclusive, &inclusive)
	queryActive.SetField("status")
	query := bleve.NewConjunctionQuery(queryPID, queryActive)
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	serarch.Fields = []string{"*"}
	result, err := bleveSearchInContext(context.Background(), indexNavMenuName, serarch)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	data := make([]map[string]interface{}, len(result.Hits))
	for i, v := range result.Hits {
		id := fmt.Sprintf("%v", v.Fields["id"]) // 强转为string

		if id != "" && id != "nil" {
			value, _ := getNavMenu(id)
			v.Fields["childs"] = value
		}
		data[i] = v.Fields
	}

	return data, err
}
