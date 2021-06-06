package main

import "fmt"

func saveNexIndex(newIndex map[string]interface{}, tableName string) (map[string]string, error) {

	SearchResult, err := findIndexFields(tableName, 1)
	m := make(map[string]string, 2)

	if err != nil {
		FuncLogError(err)
		m["code"] = "303"
		m["msg"] = "查询异常"
		return m, err
	}
	id := FuncGenerateStringID()
	newIndex["ID"] = "716ba31d-37c2-4db3-adaf-a682dfebae2d"
	result := SearchResult.Hits

	for _, v := range result {
		tmp := fmt.Sprintf("%v", v.Fields["FieldCode"]) //转为字符串
		_, ok := newIndex[tmp]
		if ok {
			if newIndex[tmp] == nil || fmt.Sprintf("%v", newIndex[tmp]) == "" {
				m["code"] = "401"
				m["msg"] = tmp + "不能为空"
				return m, nil
			}

		} else {
			m["code"] = "401"
			m["msg"] = tmp + "不能为空"
			return m, nil
		}

	}
	IndexMap[tableName].Index(id, newIndex)

	m["code"] = "200"
	m["msg"] = "保存成功"
	return m, nil

}
