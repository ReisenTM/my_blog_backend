package redis_jwt

import (
	"blog/internal/global"
	"blog/internal/utils/jwts"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// BanType 黑名单分类
type BanType uint8

const (
	UserBanType   BanType = iota + 1 //用户token自动过期
	AdminBanType                     //管理员强制下线
	DeviceBanType                    // 其他设备把自己挤下来了
)

func (b BanType) String() string {
	return fmt.Sprintf("%d", b)
}

func (b BanType) Msg() string {
	switch b {
	case UserBanType:
		return "已注销"
	case AdminBanType:
		return "禁止登录"
	case DeviceBanType:
		return "设备下线"
	}
	return "已注销"
}
func ParseBanType(val string) BanType {
	switch val {
	case "1":
		return UserBanType
	case "2":
		return AdminBanType
	case "3":
		return DeviceBanType
	}
	return UserBanType
}

// RedisBlackList 设置黑名单，在黑名单的用户无法登录
func RedisBlackList(token string, value BanType) {
	fmtedToken := fmt.Sprintf("token_black_%s", token)
	claims, err := jwts.ParseToken(token)
	if err != nil {
		logrus.Errorf("jwts.ParseToken err: %s", err.Error())
		return
	}
	//得到过期时间
	//在现在到token过期这段时间设置黑名单
	timeLast := claims.ExpiresAt - time.Now().Unix()
	_, err = global.Redis.Set(context.Background(), fmtedToken, value.String(), time.Duration(timeLast)*time.Second).Result()
	if err != nil {
		logrus.Errorf("Redis.Set err: %s", err.Error())
		return
	}
	return
}

// HasBlkList 判断是否在黑名单
func HasBlkList(token string) (banType BanType, ok bool) {
	fmtedToken := fmt.Sprintf("token_black_%s", token)
	isExist, err := global.Redis.Get(context.Background(), fmtedToken).Result()
	if err != nil {
		return
	}
	return ParseBanType(isExist), true
}

// HasTokenBlackByGin 从gin上下文获取token
func HasTokenBlackByGin(c *gin.Context) (blk BanType, ok bool) {
	// 1. 从 HTTP Header 获取 token
	token := c.GetHeader("token")
	if token == "" {
		// 2. 如果 Header 中没有，尝试从 URL Query 参数获取
		token = c.Query("token")
	}
	// 3. 调用 HasTokenBlack 检查 token 是否在黑名单中
	return HasBlkList(token)
}
