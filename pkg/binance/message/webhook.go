package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebHook interface {
	Send(keyword string, hookUrl string) error
}

/*
钉钉自定义webHook文档: https://developers.dingtalk.com/document/app/custom-robot-access?spm=ding_open_doc.document.0.0.e2b328e18qB7bT#topic-2026027
*/

/*
发送文本消息
*/

type WebHookTextMessage struct {
	MsgType string                    `json:"msgtype"` // text
	Text    WebHookMessageTextContent `json:"text"`
	At      WebHookMessageAt          `json:"at"`
}

type WebHookMessageTextContent struct {
	Content string `json:"content"`
}

type WebHookMessageAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

/*
发送链接信息
*/

type WebHookLinkMessage struct {
	MsgType string                    `json:"msgtype"` // link
	Link    WebHookMessageLinkContent `json:"link"`
}

type WebHookMessageLinkContent struct {
	Title      string `json:"title"`      // 标题
	Text       string `json:"text"`       // 正文文本
	PicUrl     string `json:"picUrl"`     // 展示图片地址
	MessageUrl string `json:"messageUrl"` // 标题跳转链接
}

/*
HTTP RESPONSE
*/

type DingResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (dingResult DingResult) String() string {
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", dingResult.ErrCode, dingResult.ErrMsg)
}

/*
发送钉钉文本消息
*/

func (m WebHookTextMessage) Send(keyword string, hookUrl string) error {
	var dingResult DingResult
	
	if keyword != "" {
		m.Text.Content = keyword + ":\n" + m.Text.Content
	}
	
	webHookMessageJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(hookUrl, "application/json", bytes.NewReader(webHookMessageJson))
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()
	
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(respByte, &dingResult)
	if err != nil {
		return err
	}
	
	if dingResult.ErrCode != 0 {
		return errors.New(dingResult.String())
	}
	
	return nil
}

/*
发送钉钉链接消息
*/

func (m WebHookLinkMessage) Send(keyword string, hookUrl string) error {
	var dingResult DingResult
	
	if keyword != "" {
		m.Link.Text = keyword + ":\n" + m.Link.Text
	}
	
	webHookMessageJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(hookUrl, "application/json", bytes.NewReader(webHookMessageJson))
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()
	
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(respByte, &dingResult)
	if err != nil {
		return err
	}
	
	return nil
}
