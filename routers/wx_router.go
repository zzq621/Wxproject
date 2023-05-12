package routers

import (
	"WxProject/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegistryWXRouter(r *gin.Engine) {
	//gptApi := r.Group("/gpt")
	//{
	//	gptApi.GET("", controller.VerifyCallBack)
	//	gptApi.POST("", controller.WxChatCommand)
	//}
	botApi := r.Group("/bot")
	{
		botApi.GET("", controller.VerifyCallBack)
		botApi.POST("", controller.TalkWeixin)

	}
}

func TestRouter(r *gin.Engine) {
	testGroup := r.Group("/test")
	testGroup.GET("", func(context *gin.Context) {
		context.String(http.StatusOK, "Pong")
	})
}
