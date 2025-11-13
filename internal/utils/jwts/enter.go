package jwts

import (
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Claims struct {
	UserID   uint          `json:"userID"`
	Username string        `json:"username"`
	Role     enum.RoleType `json:"role"`
}

type MyClaims struct {
	Claims
	jwt.StandardClaims
}

// GetUser 取出UserID
func (m MyClaims) GetUser() (user model.UserModel, err error) {
	err = global.DB.Take(&user, m.UserID).Error
	return
}

// GetToken 转换 token
func GetToken(claims Claims) (string, error) {
	cla := MyClaims{
		Claims: claims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour).Unix(), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,                                                   // 签发人
		},
	}
	//设置签名算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	return token.SignedString([]byte(global.Config.Jwt.Secret)) // 进行签名生成对应的token
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*MyClaims, error) {
	if tokenString == "" {
		//如果未登录，直接返回
		return nil, errors.New("请登录")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil
	})
	if err != nil {
		//如果出错,判断出错类型
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token过期")
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return nil, errors.New("token无效")
		}
		if strings.Contains(err.Error(), "token contains an invalid") {
			return nil, errors.New("token非法")
		}
		return nil, err
	}
	//断言确定token有效
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseTokenByGin 从请求中获取 Token并解析
func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}

	return ParseToken(token)
}

// GetClaims 从上下文中获取解析完的claims
func GetClaims(c *gin.Context) (claims *MyClaims) {
	_claims, ok := c.Get("claims")
	if !ok {
		return
	}
	claims, ok = _claims.(*MyClaims)
	if !ok {
		return
	}
	return
}
