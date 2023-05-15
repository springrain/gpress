package main

func init() {
	h.GET("/.well-known/webfinger", funcWebFinger)
	h.GET("/api/acititypub/users/:userName", funcActivityPubUserInfo)
}
