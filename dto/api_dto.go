package dto

// CallBackParams 回调参数验证model
type CallBackParams struct {
	MsgSignature string `form:"msg_signature"`
	TimeStamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
	Echostr      string `form:"echostr"`
}

// CallBackData 回调用户消息
type CallBackData struct {
	MsgSignature string `form:"msg_signature"`
	TimeStamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
}

// MsgContent  微信回调Body
type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      string `xml:"AgentId"`
}

type SimpleEvent struct {
	MsgContent
	Event    string `xml:"Event"`
	EventKey string `xml:"EventKey"`
}

// ChatCompletionMessage 请求与响应Message
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest 请求体
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionResponse 响应体
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

// Usage Token Usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type VectorData struct {
	ModelId  string              `json:"modelId"`
	IsStream bool                `json:"isStream"`
	Prompts  []map[string]string `json:"prompts"`
}

type WeixinUserAskMsg struct {
	ToUserName string `xml:"ToUserName"`
	CreateTime int64  `xml:"CreateTime"`
	MsgType    string `xml:"MsgType"`
	Event      string `xml:"Event"`
	Token      string `xml:"Token"`
	OpenKfId   string `xml:"OpenKfId"`
}

type AccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type MsgRet struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	NextCursor string `json:"next_cursor"`
	MsgList    []Msg  `json:"msg_list"`
}

type Msg struct {
	Msgid    string `json:"msgid"`
	SendTime int64  `json:"send_time"`
	Origin   int    `json:"origin"`
	Msgtype  string `json:"msgtype"`
	Event    struct {
		EventType      string `json:"event_type"`
		Scene          string `json:"scene"`
		OpenKfid       string `json:"open_kfid"`
		ExternalUserid string `json:"external_userid"`
		WelcomeCode    string `json:"welcome_code"`
	} `json:"event"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	OpenKfid       string `json:"open_kfid"`
	ExternalUserid string `json:"external_userid"`
}
type ReplyMsg struct {
	Touser   string `json:"touser,omitempty"`
	OpenKfid string `json:"open_kfid,omitempty"`
	Msgid    string `json:"msgid,omitempty"`
	Msgtype  string `json:"msgtype,omitempty"`
	Text     struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
}
