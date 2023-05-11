package controller

import (
	"WxProject/config"
	"WxProject/dto"
	"WxProject/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
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
	//
	// 异步处理用户请求
	go func() {
		err = xml.Unmarshal(userDataDecrypt, &userData)
		if err != nil {
			xlog.Log.WithError(err).Error("反序列化用户数据错误")
			return
		}
		prompt := make(map[string]string)
		prompt["obj"] = "Human"
		prompt["value"] = userData.Content
		sdata := make([]map[string]string, 0)
		sdata = append(sdata, prompt)
		//接入fast_gpt的api
		reqData := utils.HttpRequest("https://bot.botai.ai/api/openapi/chat/vectorGpt", dto.VectorData{
			ModelId:  "645311bf8a5dbd2eb7eebc6a",
			IsStream: false,
			Prompts:  sdata,
		}, "645309008a5dbd2eb7eebba9-jfjo9cnvqr6p3lh34xp87", "POST")
		xlog.Log.Info("fast_gpt返回的响应数据:", reqData)
		//返回微信信息
		wxresp := utils.SendTextToUser(userData.FromUsername, reqData)
		xlog.Log.Info("wxresp 返回的响应数据:", wxresp)
	}()
}

func TalkWeixin(c *gin.Context) {
	//verifyMsgSign := c.Query("msg_signature")
	//verifyTimestamp := c.Query("timestamp")
	//verifyNonce := c.Query("nonce")
	//bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	//data := utils.DeCryptMsg(bodyBytes, verifyMsgSign, verifyTimestamp, verifyNonce)
	fmt.Println("进来了。。。001")
	var dataStuc dto.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	// 解析请求体
	raw, err := c.GetRawData()
	fmt.Println("解析数据。。。002", raw)
	if err != nil {
		xlog.Log.WithError(err).Error("解析微信回调参数失败")
		return
	}
	//userData := dto.MsgContent{}
	userDataDecrypt := utils.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	var weixinUserAskMsg dto.WeixinUserAskMsg
	err = xml.Unmarshal(userDataDecrypt, &weixinUserAskMsg)
	if err != nil {
		fmt.Println("err:  " + err.Error())
	}
	accessToken, err := accessToken()
	if err != nil {
		c.JSON(500, "ok")
		return
	}
	msgToken := weixinUserAskMsg.Token
	msgRet, err := getMsgs(accessToken, msgToken)
	if err != nil {
		c.JSON(500, "ok")
		return
	}
	go handleMsgRet(msgRet)
	c.JSON(200, "ok")
}

func accessToken() (string, error) {
	//var tokenCacheKey = "tokenCache"
	//data, found := tokenCache.Get(tokenCacheKey)
	//if found {
	//	return fmt.Sprintf("%v", data), nil
	//}
	urlBase := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	url := fmt.Sprintf(urlBase, config.GetWechatConf().Corpid, config.GetWechatConf().CorpSecret)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	s := string(body)
	var accessToken dto.AccessToken
	json.Unmarshal([]byte(s), &accessToken)
	token := accessToken.AccessToken
	//tokenCache.Set(tokenCacheKey, token, 5*time.Minute)
	return token, nil
}

func getMsgs(accessToken, msgToken string) (dto.MsgRet, error) {
	var msgRet dto.MsgRet
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token=" + accessToken
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{"token" : "%s"}`, msgToken))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return msgRet, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return msgRet, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return msgRet, err
	}
	json.Unmarshal([]byte(string(body)), &msgRet)
	return msgRet, nil
}

func handleMsgRet(msgRet dto.MsgRet) {
	fmt.Println(msgRet)
	size := len(msgRet.MsgList)
	if size < 1 {
		return
	}
	current := msgRet.MsgList[size-1]
	userId := current.ExternalUserid
	kfId := current.OpenKfid
	content := current.Text.Content
	if content == "" {
		return
	}
	//
	prompt := make(map[string]string)
	prompt["obj"] = "Human"
	prompt["value"] = content
	sdata := make([]map[string]string, 0)
	sdata = append(sdata, prompt)
	//接入fast_gpt的api
	reqData := utils.HttpRequest("https://bot.botai.ai/api/openapi/chat/vectorGpt", dto.VectorData{
		ModelId:  "645311bf8a5dbd2eb7eebc6a",
		IsStream: false,
		Prompts:  sdata,
	}, "645309008a5dbd2eb7eebba9-jfjo9cnvqr6p3lh34xp87", "POST")
	TalkToUser(userId, kfId, content, reqData)
}

//func AskOnConversation(question, conversationId string, size int) (string, error) {
//	var messages = []gptgo.ChatCompletionMessage{}
//	//key := fmt.Sprintf("cache:conversation:%s", conversationId)
//	//data, found := conversationCache.Get(key)
//	//if found {
//	//	messages = data.([]gptgo.ChatCompletionMessage)
//	//}
//	messages = append(messages, gptgo.ChatCompletionMessage{
//		Role:    "system",
//		Content: question,
//	})
//	fmt.Println(messages)
//	pivot := size
//	if pivot > len(messages) {
//		pivot = len(messages)
//	}
//	messages = messages[len(messages)-pivot:]
//	conversationCache.Set(key, messages, 12*time.Hour)
//	chat := NewGPT(openAiKey, conversationId)
//	defer chat.Close()
//	answer, err := chat.Chat(messages)
//	if err != nil {
//		fmt.Print(err.Error())
//	}
//	return answer, err
//}

func TalkToUser(external_userid, open_kfid, ask, content string) {
	reply := dto.ReplyMsg{
		Touser:   external_userid,
		OpenKfid: open_kfid,
		Msgtype:  "text",
		Text: struct {
			Content string `json:"content,omitempty"`
		}{Content: content},
	}
	atoken, err := accessToken()
	if err != nil {
		return
	}
	callTalk(reply, atoken)
}

func callTalk(reply dto.ReplyMsg, accessToken string) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg?access_token=" + accessToken
	method := "POST"
	data, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	reqBody := string(data)
	fmt.Println(reqBody)
	payload := strings.NewReader(reqBody)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s := string(body)
	fmt.Println(s)
	return nil
}
