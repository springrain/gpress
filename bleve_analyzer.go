package main

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/char/html"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/character"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/go-ego/gse"
)

// 注册逗号分词器
func initRegisterAnalyzer() {
	registry.RegisterTokenizer(commaAnalyzerName, commaTokenizerConstructor)
	registry.RegisterAnalyzer(commaAnalyzerName, commaAnalyzerConstructor)
	registry.RegisterTokenizer(gseAnalyzerName, gseTokenizerConstructor)
	registry.RegisterAnalyzer(gseAnalyzerName, gseAnalyzerConstructor)
	registry.RegisterAnalyzer(keywordAnalyzerName, keywordlowerAnalyzerConstructor)
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
	commaTokenizer, err := cache.TokenizerNamed(commaAnalyzerName)
	if err != nil {
		return nil, err
	}
	rv := analysis.DefaultAnalyzer{
		Tokenizer: commaTokenizer,
	}
	return &rv, nil
}

func gseAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed(gseAnalyzerName)
	if err != nil {
		return nil, err
	}
	alz := &analysis.DefaultAnalyzer{Tokenizer: tokenizer}
	return alz, nil
}

type gseTokenizer struct {
	segmenter *gse.Segmenter
}

func (t *gseTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	// segments := t.segmenter.ModeSegment(sentence, true)
	segments := t.segmenter.Segment(sentence)
	for _, seg := range segments {
		token := analysis.Token{
			Term:     []byte(seg.Token().Text()),
			Start:    seg.Start(),
			End:      seg.End(),
			Position: pos,
			Type:     analysis.Ideographic,
		}
		result = append(result, &token)
		pos++
	}
	return result
}

func gseTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	var segmenter gse.Segmenter
	segmenter.SkipLog = true

	segmenter.LoadDict(datadir+"dict/zh/dict.txt", datadir+"dict/dictionary.txt")
	segmenter.LoadStop(datadir+"dict/stop_word.txt", datadir+"dict/stop_tokens.txt")

	return &gseTokenizer{&segmenter}, nil
}

func keywordlowerAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
	keywordTokenizer, err := cache.TokenizerNamed(single.Name)
	if err != nil {
		return nil, err
	}
	// 定义去除前后空格的 CharFilter
	//trimCharFilter := regexpcharfilter.New("\\A\\s+|\\s+\\z", []byte{})
	rv := analysis.DefaultAnalyzer{

		CharFilters: []analysis.CharFilter{
			html.New(),
		},

		Tokenizer: keywordTokenizer,
		TokenFilters: []analysis.TokenFilter{
			lowercase.NewLowerCaseFilter(),
		},
	}
	return &rv, nil
}
