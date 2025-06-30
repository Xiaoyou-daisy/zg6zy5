package pkg

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

// 微信开放平台的配置
var (
	wechatConfig = &oauth2.Config{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		Scopes:       []string{"snsapi_userinfo"},
		RedirectURL:  "your_redirect_url",
	}
)

func main() {
	http.HandleFunc("/login/wechat", wechatLoginHandler)
	http.HandleFunc("/callback/wechat", wechatCallbackHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

// 引导用户到微信授权页面
func wechatLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := wechatConfig.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}

// 处理微信回调
func wechatCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := wechatConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Access Token: %s", token.AccessToken)
}
