package codeInHere

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
)

func BenchmarkExample(b *testing.B) {
	//调用runtime/prof包 进行性能测试
	// CPU Profile
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	for i := 0; i < 10; i++ {
		Example()
	}

	// Memory Profile
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func Example() {
	INDEX_DIR := "gojieba.bleve"
	messages := []struct {
		Id   string
		Body string
	}{
		{
			Id:   "1",
			Body: "你好",
		},
		{
			Id:   "2",
			Body: "交代",
		},
		{
			Id:   "3",
			Body: "长江大桥",
		},
	}

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(INDEX_DIR)
	// clean index when example finished
	defer os.RemoveAll(INDEX_DIR)

	//加入gojieba 分词器
	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":     gojieba.DICT_PATH,
			"hmmpath":      gojieba.HMM_PATH,
			"userdictpath": gojieba.USER_DICT_PATH,
			"idf":          gojieba.IDF_PATH,
			"stop_words":   gojieba.STOP_WORDS_PATH,
			"type":         "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	//添加分析器
	// 这里会去找分析器 analyzer.go
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	indexMapping.DefaultAnalyzer = "gojieba"

	index, err := bleve.New(INDEX_DIR, indexMapping)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		if err := index.Index(msg.Id, msg); err != nil {
			panic(err)
		}
	}

	querys := []string{
		"你好世界",
		"亲口交代",
		"长江",
	}

	for _, q := range querys {
		//循环查询。
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("搜索关键字:", q)
		fmt.Println("搜索结果", res)
		fmt.Println(prettify(res))
		fmt.Println("_________________________________")
	}

	//cleanup cgo allocated heap memory
	if jieba, ok := (index.Mapping().AnalyzerNamed("gojieba").Tokenizer).(*JiebaTokenizer); !ok {
		panic("jieba.Free() failed")
	} else {
		jieba.Free()
	}
	index.Close()
	// Output:
	// [{"id":"1","score":0.27650412875470115}]
	// [{"id":"2","score":0.27650412875470115}]
	// [{"id":"3","score":0.7027325540540822}]
}

func prettify(res *bleve.SearchResult) string {
	type Result struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
	results := []Result{}
	for _, item := range res.Hits {
		results = append(results, Result{item.ID, item.Score})
	}
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(b)
}
