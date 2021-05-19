package main

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/go-ego/gse"
)

const (
	gseAnalyzerName = "gse"
)

func gseAnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {

	tokenizer, err := cache.TokenizerNamed(gseAnalyzerName)
	if err != nil {
		return nil, err
	}
	alz := &analysis.Analyzer{Tokenizer: tokenizer}
	return alz, nil
}

func init() {
	registry.RegisterTokenizer(gseAnalyzerName, gseTokenizerConstructor)
	registry.RegisterAnalyzer(gseAnalyzerName, gseAnalyzerConstructor)
}

type gseTokenizer struct {
	segmenter *gse.Segmenter
}

func (t *gseTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	//segments := t.segmenter.ModeSegment(sentence, true)
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
	dicts := datadir + "/dict/zh/dict.txt"
	var segmenter gse.Segmenter
	segmenter.SkipLog = true
	segmenter.LoadDict(dicts)

	return &gseTokenizer{&segmenter}, nil

}