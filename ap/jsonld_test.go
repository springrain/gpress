package ap

import (
	"fmt"
	"testing"

	"github.com/piprate/json-gold/ld"
)

/**
https://lawrenceli.me/blog/activitypub
https://wangqiao.me/posts/activitypub-from-decentralized-to-distributed-social-networks/
https://blog.joinmastodon.org/2018/06/how-to-implement-a-basic-activitypub-server/
https://www.w3.org/TR/activitystreams-vocabulary/
**/

func TestExpand(t *testing.T) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// expanding in-memory document
	doc := map[string]interface{}{
		"@context":          "https://www.w3.org/ns/activitystreams",
		"id":                "https://lawrenceli.me/api/activitypub/actor",
		"type":              "Person",
		"name":              "Lawrence Li",
		"preferredUsername": "lawrence",
		"summary":           "Blog",
		"inbox":             "https://lawrenceli.me/api/activitypub/inbox",
		"outbox":            "https://lawrenceli.me/api/activitypub/outbox",
		"followers":         "https://lawrenceli.me/api/activitypub/followers",
	}

	expanded, err := proc.Expand(doc, options)
	if err != nil {
		t.Error("Error when expanding JSON-LD document:", err)
		return
	}
	fmt.Println(expanded)
}
