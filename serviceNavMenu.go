package main

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

//获取菜单树 pid 为0 为一级菜单
var active float64 = 1

func getNavMenu(pid string) (interface{}, error) {

	navIndex := IndexMap[indexNavMenuName]
	//PID 跟 Active 为查询字段
	queryPID := bleve.NewTermQuery(pid)
	queryPID.SetField("PID")
	queryActive := bleve.NewNumericRangeInclusiveQuery(&active, &active, &inclusive, &inclusive)
	queryActive.SetField("Active")
	query := bleve.NewConjunctionQuery(queryPID, queryActive)
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	serarch.Fields = []string{"*"}
	result, err := navIndex.SearchInContext(context.Background(), serarch)

	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	data := make([]map[string]interface{}, len(result.Hits))
	for i, v := range result.Hits {
		id := fmt.Sprintf("%v", v.Fields["ID"]) //强转为string

		if id != "" && id != "nil" {
			value, _ := getNavMenu(id)
			v.Fields["childs"] = value
		}
		data[i] = v.Fields
	}

	return data, err

}
