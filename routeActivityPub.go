package main

func init() {
	h.GET("/.well-known/webfinger", funcWebFinger)
	h.GET("/acititypub/api/user/:userName", funcActivityPubUsers)
	h.GET("/acititypub/api/outbox/:userName", funcActivityPubOutBox)
}
