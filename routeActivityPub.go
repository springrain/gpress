package main

func init() {
	//resource根信息查询
	h.GET("/.well-known/webfinger", funcWebFinger)

	//用户信息查询
	h.GET("/activitypub/api/user/:userName", funcActivityPubUsers)
	//outbox信息查询
	h.GET("/activitypub/api/outbox/:userName", funcActivityPubOutBox)
	h.GET("/activitypub/api/outbox_page/:userName/:pageNo", funcActivityPubOutBoxPage)

	//inbox信息查询
	h.POST("/activitypub/api/inbox/:userName", activitySignatureHandler, funcActivityPubInBox)

	//关注
	h.GET("/activitypub/api/following/:userName", funcActivityPubOutBox)
	//粉丝
	h.GET("/activitypub/api/followers/:userName", funcActivityPubOutBox)

}
