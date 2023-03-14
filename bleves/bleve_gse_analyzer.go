package bleves

import (
	"gitee.com/gpress/gpress/configs"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/go-ego/gse"
)

// 中文分词器的mapping
var gseAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

// 注册分词器
func init() {
	registry.RegisterTokenizer(configs.GSE_ANGLYZER_NAME, gseTokenizerConstructor)
	registry.RegisterAnalyzer(configs.GSE_ANGLYZER_NAME, gseAnalyzerConstructor)
}

func gseAnalyzerConstructor(interfaceMap map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed(configs.GSE_ANGLYZER_NAME)
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

func gseTokenizerConstructor(interfaceMap map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	var segmenter gse.Segmenter
	segmenter.SkipLog = true

	segmenter.LoadDict(configs.DATA_DIR+"dict/zh/dict.txt", configs.DATA_DIR+"dict/dictionary.txt")
	segmenter.LoadStop(configs.DATA_DIR+"dict/stop_word.txt", configs.DATA_DIR+"dict/stop_tokens.txt")

	return &gseTokenizer{&segmenter}, nil
}
