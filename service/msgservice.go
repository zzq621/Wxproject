package service

import (
	"WxProject/config"
	"WxProject/dto"
	"WxProject/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func AccessToken() (string, error) {
	urlBase := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	url := fmt.Sprintf(urlBase, config.GetWechatConf().Corpid, config.GetWechatConf().CorpSecret)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	s := string(body)
	var accessToken dto.AccessToken
	json.Unmarshal([]byte(s), &accessToken)
	token := accessToken.AccessToken
	return token, nil
}

func GetMsgs(accessToken, msgToken string) (dto.MsgRet, error) {
	var msgRet dto.MsgRet
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token=" + accessToken
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{"token" : "%s"}`, msgToken))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
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
		return msgRet, err
	}
	json.Unmarshal([]byte(string(body)), &msgRet)
	return msgRet, nil
}

func HandleMsgRet(msgRet dto.MsgRet) {
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
	prompt := make(map[string]string)
	sdata := make([]map[string]string, 0)
	prompt["obj"] = "Human"
	prompt["value"] = content
	sdata = append(sdata, prompt)
	//接入fast_gpt的api
	reqData := utils.HttpRequest(config.GetGptConf().BotApiUrl, dto.VectorData{
		ModelId:  config.GetGptConf().ModelId,
		IsStream: false,
		Prompts:  sdata,
	}, config.GetGptConf().Apikey, "POST")
	if reqData == "" {
		reqData = "抱歉，你的问题不在知识库中。。。"
	}
	TalkToUser(userId, kfId, content, strings.TrimSpace(reqData))
}

func TalkToUser(external_userid, open_kfid, ask, content string) {
	reply := dto.ReplyMsg{
		Touser:   external_userid,
		OpenKfid: open_kfid,
		Msgtype:  "text",
		Text: struct {
			Content string `json:"content,omitempty"`
		}{Content: content},
	}
	atoken, err := AccessToken()
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
	return nil
}
