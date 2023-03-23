package service

import (
	"context"
	"fmt"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/gbleve"
	"gitee.com/gpress/gpress/logger"
	"github.com/blevesearch/bleve/v2"
)

func GetNavMenu(pid string) (interface{}, error) {
	navIndex := gbleve.IndexMap[constant.INDEX_NAV_MENU_NAME]
	// PID 跟 Active 为查询字段
	queryPID := bleve.NewTermQuery(pid)
	queryPID.SetField("PID")
	queryActive := bleve.NewNumericRangeInclusiveQuery(&Active, &Active, &Inclusive, &Inclusive)
	queryActive.SetField("Active")
	query := bleve.NewConjunctionQuery(queryPID, queryActive)
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	serarch.Fields = []string{"*"}
	result, err := navIndex.SearchInContext(context.Background(), serarch)
	if err != nil {
		logger.FuncLogError(err)
		return nil, err
	}
	data := make([]map[string]interface{}, len(result.Hits))
	for i, v := range result.Hits {
		id := fmt.Sprintf("%v", v.Fields["ID"]) // 强转为string

		if id != "" && id != "nil" {
			value, _ := GetNavMenu(id)
			v.Fields["childs"] = value
		}
		data[i] = v.Fields
	}

	return data, err
}
