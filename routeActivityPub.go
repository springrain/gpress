package main

func init() {
	h.GET("/.well-known/webfinger", funcWebFinger)
	h.GET("/activitypub/api/user/:userName", funcActivityPubUsers)
	h.GET("/activitypub/api/outbox/:userName", funcActivityPubOutBox)
	h.GET("/activitypub/api/outbox_page/:userName/:pageNo", funcActivityPubOutBoxPage)
}
