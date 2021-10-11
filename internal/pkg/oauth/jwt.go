package auth

import (
	"errors"
	"github.com/102er/go-apiserver/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	Account         string `json:"account"`
	Role            string `json:"role"`
	UUAPAccessToken string `json:"accessToken"`
	jwt.StandardClaims
}

var jwtSecret = []byte("秘钥你猜猜看")

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"reason":  "4",
				"message": "token is required",
			})
		}
		//解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer" && len(parts[1]) > 0) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   400,
				"reason": "4",
				"msg":    "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		//校验token
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":   403,
				"reason": "4",
				"msg":    "无效token",
			})
			c.Abort()
			return
		}
		c.Set("username", mc.Account)
		c.Set("access_token", mc.UUAPAccessToken)
		c.Next()
	}
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func CreateToken(accessToken string, userInfo domain.User) (token string, err error) {
	now := time.Now()
	account := userInfo.Name
	jwtId := account + strconv.FormatInt(now.Unix(), 10)
	role := "member"
	claims := Claims{
		Account:         account,
		Role:            role,
		UUAPAccessToken: accessToken,
		StandardClaims: jwt.StandardClaims{
			Audience:  account,
			ExpiresAt: now.Add(30 * 24 * time.Hour).Unix(),
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "api-server",
			Subject:   account,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(jwtSecret)
	return token, err
}
