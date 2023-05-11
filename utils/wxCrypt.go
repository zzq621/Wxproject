package utils

import (
	"WxProject/config"
	"WxProject/dto"
	"WxProject/utils/xlog"
	"context"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/response"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

var wxCrypt *wxbizmsgcrypt.WXBizMsgCrypt

// var wkCrypt *wxbizmsgcrypt.WXBizMsgCrypt
var WeComApp *work.Work

func init() {
	LoadWeComAppConf()
	LoadWxUtils()
}

func LoadWeComAppConf() {
	xlog.Log.Info("初始化企业微信助手......")
	app, err := work.NewWork(&work.UserConfig{
		CorpID:  config.GetWechatConf().Corpid,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: config.GetWechatConf().AgentId,    // 内部应用的app id
		Secret:  config.GetWechatConf().CorpSecret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: config.GetSystemConf().CallBackUrl, //
			Scopes:   nil,
		},
		HttpDebug: false,
	})
	if err != nil {
		xlog.Log.WithError(err).Error("初始化企业微信助手失败！")
		panic(err)
	}
	WeComApp = app
}

func LoadWxUtils() {
	xlog.Log.Info("初始化微信工具包......")
	wxCrypt = wxbizmsgcrypt.NewWXBizMsgCrypt(config.GetWechatConf().WeApiRCallToken, config.GetWechatConf().WeApiEncodingKey, config.GetWechatConf().Corpid, wxbizmsgcrypt.XmlType)
	//wkCrypt = wxbizmsgcrypt.NewWXBizMsgCrypt(config.GetWechatConf().WkApiRCallToken, config.GetWechatConf().WkApiEncodingKey, config.GetWechatConf().Corpid, wxbizmsgcrypt.XmlType)
}

// GetReVerifyCallBack 从微信回调解析请求数据
func GetReVerifyCallBack(q dto.CallBackParams) []byte {
	msg, cryptErr := wxCrypt.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.Echostr)
	if cryptErr != nil {
		xlog.Log.Errorf("验证Url出错（回调消息解密错误）：%v", cryptErr)
		return []byte("")
	}
	xlog.Log.Info("解析的回调字符为：", string(msg))
	return msg
}

// DeCryptMsg 解密消息
func DeCryptMsg(cryptMsg []byte, msgSignature, timeStamp, nonce string) []byte {
	msg, cryptErr := wxCrypt.DecryptMsg(msgSignature, timeStamp, nonce, cryptMsg)
	if cryptErr != nil {
		xlog.Log.Errorf("回调消息解密错误：%v", cryptErr)
		return nil
	}
	return msg
}

// CryptMessage 加密消息
func CryptMessage(respData, reqTimestamp, reqNonce string) string {
	encryptMsg, cryptErr := wxCrypt.EncryptMsg(respData, reqTimestamp, reqNonce)
	if cryptErr != nil {
		//xlog.Log.Errorf("消息加密错误：%v", cryptErr)
		return ""
	}
	return string(encryptMsg)
}

// SendTextToUser 发送text消息
func SendTextToUser(userName string, respMsg string) *response.ResponseMessageSend {
	// 封装微信消息体
	messages := &request.RequestMessageSendText{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 userName,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "text",
			AgentID:                config.GetWechatConf().AgentId,
			Safe:                   0,
			EnableIDTrans:          0,
			EnableDuplicateCheck:   0,
			DuplicateCheckInterval: 1800,
		},
		Text: &request.RequestText{
			Content: respMsg,
		},
	}
	// 发送微信消息
	resp, err := WeComApp.Message.SendText(context.Background(), messages)
	if err != nil {
		xlog.Log.Errorf("创建微信发送消息内容失败：%v", err)
		return nil
	}
	return resp
}
