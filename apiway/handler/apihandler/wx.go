package apihandler

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"zg6zy5/apiway/inits"
)

var appid = "wxd25f742200bba77e"                //填你的
var secret = "89d24ad4d4feadb6b88fd86abaf363f0" //填你的

type WxUsers struct {
	Id       int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	OpenId   string `gorm:"column:open_id;type:varchar(80);comment:打开的ID;not null;" json:"open_id"` // 打开的ID
	NickName string `gorm:"column:nick_name;type:varchar(201);not null;" json:"nick_name"`
}

type SignReq struct {
	Signature string `json:"signature" form:"signature" binding:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" binding:"required"`
	Nonce     string `json:"nonce" form:"nonce" binding:"required"`
	Echostr   string `json:"echostr" form:"echostr" binding:"required"`
}

func Sign(c *gin.Context) {
	var req SignReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := "555" //与后台token相同

	// 1. 将 token、timestamp、nonce 放入切片
	tmpArr := []string{token, req.Timestamp, req.Nonce}

	// 2. 对切片进行字典序排序
	sort.Strings(tmpArr)

	// 3. 拼接字符串
	tmpStr := strings.Join(tmpArr, "")

	// 4. 计算 SHA1 哈希值
	hash := sha1.New()
	hash.Write([]byte(tmpStr))
	computedSignature := fmt.Sprintf("%x", hash.Sum(nil))

	// 5. 比较签名
	if computedSignature == req.Signature {
		c.Writer.Write([]byte(req.Echostr))
		return
	}

}

var (
	redirect = "http://687e2a90.r18.cpolar.top/v1/calblack" // 微信授权回调地址
)

func One(c *gin.Context) {

	redirectURI := url.QueryEscape(redirect) // redirect 是你 `/v1/one` 地址
	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect",
		appid, redirectURI)

	qr, err := qrcode.New(authURL, qrcode.Medium)
	if err != nil {
		c.JSON(500, gin.H{"error": "二维码生成失败", "detail": err.Error()})
		return
	}

	png, err := qr.PNG(256)
	if err != nil {
		c.JSON(500, gin.H{"error": "二维码转码失败", "detail": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(png)
	c.Writer.Write([]byte(authURL))

}

type CalblackReq struct {
	Code string `form:"code" binding:"required"`
}

func Calblack(c *gin.Context) {
	// 解析请求参数
	var req CalblackReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 获取 access_token 和 openid
	token, openid := GetAccessTokenByCode(req.Code)
	if token == "" || openid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取access_token或openid失败"})
		return
	}
	fmt.Println(token)
	fmt.Println(openid)

	// 拉取用户信息
	res, err := http.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", token, openid,
	))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取用户信息失败: " + err.Error()})
		return
	}
	defer res.Body.Close()

	resBytes, _ := io.ReadAll(res.Body)
	var userInfo map[string]interface{}
	if err := json.Unmarshal(resBytes, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息解析失败"})
		return
	}

	// 构造用户对象
	user := WxUsers{
		OpenId:   fmt.Sprintf("%v", userInfo["open_id"]),
		NickName: fmt.Sprintf("%v", userInfo["nick_name"]),
	}

	// 写入数据库
	if err := inits.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库插入失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetAccessTokenByCode(code string) (accessToken string, openid string) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appid, secret, code,
	)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("请求 access_token 失败:", err)
		return "", ""
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("access_token 响应解析失败:", err)
		return "", ""
	}

	// 错误处理
	if _, ok := data["errcode"]; ok {
		fmt.Printf("微信接口错误: %+v\n", data)
		return "", ""
	}

	// 获取 token 和 openid
	accessToken = fmt.Sprintf("%v", data["access_token"])
	openid = fmt.Sprintf("%v", data["openid"])
	return accessToken, openid
}
