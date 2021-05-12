package main

import (
	"fmt"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/character"
	"github.com/blevesearch/bleve/v2/registry"
)

const CommaName = "comma"

func TokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return character.NewCharacterTokenizer(IsComma), nil
}

// IsLetter reports whether the rune is a letter (category L).
func IsComma(r rune) bool {
	fmt.Println(r)
	return false

}
func AnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	keywordTokenizer, err := cache.TokenizerNamed(CommaName)
	if err != nil {
		return nil, err
	}
	rv := analysis.Analyzer{
		Tokenizer: keywordTokenizer,
	}
	return &rv, nil
}

func init() {
	registry.RegisterTokenizer(CommaName, TokenizerConstructor)
	registry.RegisterAnalyzer(CommaName, AnalyzerConstructor)
}
