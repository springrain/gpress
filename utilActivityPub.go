package main

import (
	"github.com/go-ap/activitypub"
)

func createAPObject() error {
	apObject := activitypub.ObjectNew(activitypub.ArticleType)
	apObject.ID = activitypub.ID("gpress")

	return nil

}
