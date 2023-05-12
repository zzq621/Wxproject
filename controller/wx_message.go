package controller

import (
	"WxProject/config"
	"WxProject/dto"
	"WxProject/service"
	"WxProject/utils"
	"encoding/xml"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/gin-gonic/gin"
	"net/http"
)

// VerifyCallBack 回调验证
func VerifyCallBack(c *gin.Context) {
	var q dto.CallBackParams
	if err := c.Bind(&q); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	msg := utils.GetReVerifyCallBack(q)
	_, _ = c.Writer.Write(msg)
}

// WxChatCommand 实际处理用户消息
func WxChatCommand(c *gin.Context) {
	var dataStuc dto.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	// 解析请求体
	raw, err := c.GetRawData()
	if err != nil {
		xlog.Log.WithError(err).Error("解析微信回调参数失败")
		return
	}
	userData := dto.MsgContent{}
	userDataDecrypt := utils.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	// 解密失败返回空
	if userDataDecrypt == nil {
		xlog.Log.WithField("用户数据：", userData).Error("解密失败")
	}
	// 提前向微信返回成功接受，防止微信多次回调
	c.JSON(http.StatusOK, "")
	// 异步处理用户请求
	go func() {
		err = xml.Unmarshal(userDataDecrypt, &userData)
		if err != nil {
			xlog.Log.WithError(err).Error("反序列化用户数据错误")
			return
		}
		prompt := make(map[string]string)
		sdata := make([]map[string]string, 0)
		prompt["obj"] = "Human"
		prompt["value"] = userData.Content
		sdata = append(sdata, prompt)
		//接入fast_gpt的api
		reqData := utils.HttpRequest(config.GetGptConf().BotApiUrl, dto.VectorData{
			ModelId:  config.GetGptConf().ModelId,
			IsStream: false,
			Prompts:  sdata,
		}, config.GetGptConf().Apikey, "POST")
		xlog.Log.Info("fast_gpt返回的响应数据:", reqData)
		//返回微信信息
		wxresp := utils.SendTextToUser(userData.FromUsername, reqData)
		xlog.Log.Info("wxresp 返回的响应数据:", wxresp)
	}()
}

func TalkWeixin(c *gin.Context) {
	var dataStuc dto.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	// 解析请求体
	raw, err := c.GetRawData()
	if err != nil {
		xlog.Log.WithError(err).Error("解析微信回调参数失败")
		return
	}
	userDataDecrypt := utils.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	var weixinUserAskMsg dto.WeixinUserAskMsg
	err = xml.Unmarshal(userDataDecrypt, &weixinUserAskMsg)
	if err != nil {
		xlog.Log.Errorf("反序列化数据错误：%v", err)
		return
	}
	accessToken, err := service.AccessToken()
	if err != nil {
		xlog.Log.Errorf("获取accesstoken错误：%v", err)
		return
	}
	msgToken := weixinUserAskMsg.Token
	msgRet, err := service.GetMsgs(accessToken, msgToken)
	if err != nil {
		return
	}
	//企业微信重试缓存
	if service.IsRetry(dataStuc.MsgSignature) {
		c.JSON(http.StatusOK, "ok")
		return
	}
	go service.HandleMsgRet(msgRet)
	//提前向微信返回成功接受，防止微信多次回调
	c.JSON(http.StatusOK, "ok")
}
