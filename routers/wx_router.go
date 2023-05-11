package routers

import (
	"WxProject/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegistryWXRouter(r *gin.Engine) {
	gptApi := r.Group("/gpt")
	{
		gptApi.GET("", controller.VerifyCallBack)
		gptApi.POST("", controller.WxChatCommand)
		//gptApi.GET("", controller.WxChatCommand)
	}

	botApi := r.Group("/bot")
	{
		botApi.GET("", controller.VerifyCallBack)
		botApi.POST("", controller.TalkWeixin)
		//gptApi.GET("", controller.WxChatCommand)

	}
}

func TestRouter(r *gin.Engine) {
	testGroup := r.Group("/test")
	testGroup.GET("", func(context *gin.Context) {
		context.String(http.StatusOK, "Pong")
	})
}
