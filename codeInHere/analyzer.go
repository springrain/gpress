package codeInHere

import (
	"errors"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

// 当调用analysis这个包时  会执行这里的方法

type JiebaAnalyzer struct {
}

//构造器
func analyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizerName, ok := config["tokenizer"].(string)
	if !ok {
		return nil, errors.New("must specify tokenizer")
	}
	tokenizer, err := cache.TokenizerNamed(tokenizerName)
	if err != nil {
		return nil, err
	}
	alz := &analysis.Analyzer{
		Tokenizer: tokenizer,
	}
	return alz, nil
}

//注册gojieba分析器
//这里如果不注入 会找不到分析器
func init() {
	registry.RegisterAnalyzer("gojieba", analyzerConstructor)
}
