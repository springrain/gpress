package gbleve

import (
	"gitee.com/gpress/gpress/constant"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/character"
	"github.com/blevesearch/bleve/v2/registry"
)

// InitCommaAnalyzer 注册分词器
func InitCommaAnalyzer() {
	registry.RegisterTokenizer(constant.COMMA_ANALYZER_NAME, commaTokenizerConstructor)
	registry.RegisterAnalyzer(constant.COMMA_ANALYZER_NAME, commaAnalyzerConstructor)
}

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
	return r-44 != 0
}

func commaAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
	commaTokenizer, err := cache.TokenizerNamed(constant.COMMA_ANALYZER_NAME)
	if err != nil {
		return nil, err
	}
	rv := analysis.DefaultAnalyzer{
		Tokenizer: commaTokenizer,
	}
	return &rv, nil
}
