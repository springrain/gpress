package main

func init() {
	h.GET("/.well-known/webfinger", funcWebFinger)
}
