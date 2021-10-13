package auth

/*
 * 基于gin的oidc + oauth2的keycloak适配器
 */
import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strconv"
	"time"
)

type IAuthService interface {
	GetServerLoginUrl() string

	LoginHandler(c *gin.Context)
	LoginCallBackHandler(c *gin.Context)
}

type Client struct {
	OAuth2         *oauth2.Config
	OidcVerifier   *oidc.IDTokenVerifier
	State          string //随机状态位
	RedirectUrl    string //回调地址
	ServerLoginUrl string //登录地址
}

type P struct {
	ProviderIssuer     string
	Oauth2ClientID     string
	Oauth2ClientSecret string
	RedirectURL        string
	ServerLogin        string
}

func NewClient(p P) *Client {
	ctx := &gin.Context{}
	provider, err := oidc.NewProvider(ctx, p.ProviderIssuer) //获取endpoint信息，但是我们的iam架构之间加了一层， 所以并没有使用这个返回结果
	if err != nil {
		log.Fatal("iam provider new failed,error:", err)
	}
	oauth2Config := &oauth2.Config{
		ClientID:     p.Oauth2ClientID,
		ClientSecret: p.Oauth2ClientSecret,
		RedirectURL:  p.RedirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: oauth2.Endpoint{
			AuthURL:   provider.Endpoint().AuthURL,
			TokenURL:  provider.Endpoint().TokenURL,
			AuthStyle: provider.Endpoint().AuthStyle,
		},
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	oidcConfig := &oidc.Config{
		ClientID:          p.Oauth2ClientID,
		SkipClientIDCheck: true, //跳过 client id 校验
	}
	verifier := provider.Verifier(oidcConfig)
	return &Client{
		OAuth2:         oauth2Config,
		OidcVerifier:   verifier,
		State:          strconv.FormatInt(time.Now().Unix(), 10), //理论上需要随机生成，随机生成 需要用额外的存储 才能校验
		RedirectUrl:    p.RedirectURL,
		ServerLoginUrl: p.ServerLogin,
	}
}

// GetServerLoginUrl 获取本地服务登录地址
func (i *Client) GetServerLoginUrl() string {
	return i.ServerLoginUrl
}

// LoginHandler 登录接口
func (i *Client) LoginHandler(c *gin.Context) {
	loginUrl := i.OAuth2.AuthCodeURL(i.State)
	if len(loginUrl) == 0 {
		c.JSON(http.StatusInternalServerError, "get iam auth url failed")
		return
	}
	c.Redirect(http.StatusFound, loginUrl)
}

// LoginCallBackHandler 登录验证 以及 iam验证接口
func (i *Client) LoginCallBackHandler(c *gin.Context) {
	ctx := context.Background()
	r := c.Request
	var err error
	var ok bool
	var accessToken string
	var oauth2Token *oauth2.Token
	//不带这两个参数 则认为不是从iam回调回来的请求 默认未登录
	if c.Query("state") == "" || c.Query("code") == "" {
		//认证失败 自己定义返回信息
		return
	}
	//验证iam返回的state
	if r.URL.Query().Get("state") != i.State {
		//认证失败 自己定义返回信息
		return
	}
	//iam返回的授权码code 去交换token
	oauth2Token, err = i.OAuth2.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		//认证失败 自己定义返回信息
		return
	}
	//oidc模式 取id_token
	//oauth模式 取 access_token
	accessToken, ok = oauth2Token.Extra("access_token").(string)
	if !ok {
		//认证失败 自己定义返回信息
		return
	}
	//从access token中解析 用户信息
	idToken, err := i.OidcVerifier.Verify(ctx, accessToken)
	if err != nil {
		//解析用户信息失败 自定义返回信息
		return
	}
	IDTokenClaims := jwt.MapClaims{}
	if err = idToken.Claims(&IDTokenClaims); err != nil {
		//解析用户信息失败 自定义返回信息
		return
	}
	//拿到用户信息 就可以生成 本地token
	// nil 返回自己需要的数据
	c.JSON(http.StatusOK, nil)
}
