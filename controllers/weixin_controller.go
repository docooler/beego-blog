package controllers

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ulricqin/beego-blog/models/loveletter"
	"io/ioutil"
	"math/rand"
	"sort"
	"time"
)

const (
	TOKEN    = "wechat_print_only_key"
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
)

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

type Request struct {
	XMLName                xml.Name `xml:"xml"`
	msgBase                         // base struct
	Location_X, Location_Y float32
	Scale                  int
	Label                  string
	PicUrl                 string
	MsgId                  int
}

type Response struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
	FuncFlag     int
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type WeiXinController struct {
	beego.Controller
}

func (this *WeiXinController) IsValid() (bool, string) {
	signature := this.Input().Get("signature")
	beego.Info(signature)
	timestamp := this.Input().Get("timestamp")
	beego.Info(timestamp)
	nonce := this.Input().Get("nonce")
	beego.Info(nonce)
	echostr := this.Input().Get("echostr")
	beego.Info(echostr)
	beego.Info(Signature(timestamp, nonce))
	if Signature(timestamp, nonce) == signature {
		return true, echostr
	} else {
		return false, ""
	}

}

func (this *WeiXinController) Get() {
	isValied, echostr := this.IsValid()
	if isValied {
		this.Ctx.WriteString(echostr)
	} else {
		this.Ctx.WriteString("")
	}
}

func (this *WeiXinController) Post() {

	isValid, _ := this.IsValid()
	if !isValid {
		beego.Info("wrong user post request")
		return
	}
	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(wreq.Content)
	wresp, err := dealwith(wreq)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	data, err := wresp.Encode()
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	this.Ctx.WriteString(string(data))
	return
}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text
	beego.Info(req.MsgType)
	beego.Info(req.Content)
	if req.MsgType == Text {
		maxLetter := loveletter.GetNrLetter()
		// 回复消息

		index := rand.Intn(maxLetter)

		resp.Content = loveletter.Get(index)
	} else {
		resp.Content = "暂时还不支持其他的类型"
	}
	return resp, nil
}

func Signature(timestamp, nonce string) string {
	strs := sort.StringSlice{TOKEN, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp Response) Encode() (data []byte, err error) {
	resp.CreateTime = time.Second
	data, err = xml.Marshal(resp)
	return
}
