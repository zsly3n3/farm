package thirdParty

import (
	"bytes"
	"encoding/json" //json封装解析
	"farm/log"
	"fmt"
	"io/ioutil"
	"net/http"
)

const wx_appid = "wxc011f016e4d7fe45"
const wx_appsecret = "c677332b922a827e83de110195d77dae"

type WX_OPENID struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
}

func GetWXOpenID(code string) string {
	str := ""
	var buf bytes.Buffer
	buf.WriteString("https://api.weixin.qq.com/sns/jscode2session?appid=" + wx_appid)
	buf.WriteString("&secret=" + wx_appsecret)
	buf.WriteString("&js_code=" + code)
	buf.WriteString("&grant_type=authorization_code")
	url := buf.String()
	p_body := httpGet(url)
	wx_data := new(WX_OPENID)
	if json_err := json.Unmarshal(*p_body, wx_data); json_err == nil {
		str = wx_data.OpenId
	}
	return str
}

func httpGet(url string) *[]byte {
	log.Debug("httpGet_url:%v", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &body
}

/*
type WX_ACCESS_TOKEN struct {
	Expires int `json:"expires_in"`
	Token string `json:"access_token"`
}
func getWXToken() string{
	str:=""
	var buf bytes.Buffer
	buf.WriteString("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+wx_appid)
	buf.WriteString("&secret="+wx_appsecret)
	url:=buf.String()
	p_body:=httpGet(url)
	wx_data := new(WX_ACCESS_TOKEN)
	if json_err := json.Unmarshal(*p_body, wx_data); json_err == nil {
     str = wx_data.Token
	}
	return str
}


type InviteQRCode struct {
	QRCode string `json:"qrcode"`
}

func GetQRCode(key string)string{
	 token:=getWXToken()
	 var buf bytes.Buffer
	 buf.WriteString(conf.Server.LocalHttpServer)
	 buf.WriteString("/generateQRCode/"+key)
	 buf.WriteString("/"+token)
	 url:=buf.String()
	 p_body:=httpGet(url)
	 str:=""
	 data := new(InviteQRCode)
	 if json_err := json.Unmarshal(*p_body, data); json_err == nil {
	  str = conf.Server.RemoteHttpServer+"/"+ data.QRCode
	 }
	 return str
}

func RemoveQRCode(key string){
	var buf bytes.Buffer
	buf.WriteString(conf.Server.LocalHttpServer)
	buf.WriteString("/deleteQRCode/"+key)
	url:=buf.String()
	httpGet(url)
}
*/
