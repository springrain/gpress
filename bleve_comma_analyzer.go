package main

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/character"
	"github.com/blevesearch/bleve/v2/registry"
)

const commaName = "comma"

func commaTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return character.NewCharacterTokenizer(isComma), nil
}

// 是否是分号(,) 如果是分号 返回false
func isComma(r rune) bool {
	//44 就是 ,
	/*
		if r-44 == 0 {
			//fmt.Println(r, strconv.QuoteRune(r))
			return false
		}
		return true
	*/
	return r != 44

}
func commaAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	commaTokenizer, err := cache.TokenizerNamed(commaName)
	if err != nil {
		return nil, err
	}
	rv := analysis.Analyzer{
		Tokenizer: commaTokenizer,
	}
	return &rv, nil
}

// 注册分词器
func init() {
	registry.RegisterTokenizer(commaName, commaTokenizerConstructor)
	registry.RegisterAnalyzer(commaName, commaAnalyzerConstructor)
}
