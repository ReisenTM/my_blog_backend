// Package middleware 权限验证
package middleware

import (
	"blog/internal/common/resp"
	"blog/internal/model/enum"
	"blog/internal/service/redis_service/redis_jwt"
	"blog/internal/utils/jwts"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 权限验证
func AuthMiddleware(c *gin.Context) {

	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		resp.FailWithError(err, c)
		c.Abort()
		return
	}
	//确认用户不在黑名单
	blcType, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		resp.FailWithMsg(blcType.Msg(), c)
		c.Abort()
		return
	}
	//保存验证过的用户信息
	c.Set("claims", claims)
	return
}

// AdminMiddleware 管理员级验证
func AdminMiddleware(c *gin.Context) {

	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		resp.FailWithError(err, c)
		c.Abort()
		return
	}

	if claims.Role != enum.RoleAdminType {
		//不是管理员
		resp.FailWithMsg("权限错误", c)
		c.Abort()
		return
	}
	//确认用户不在黑名单
	blcType, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		resp.FailWithMsg(blcType.Msg(), c)
		c.Abort()
		return
	}
	//保存验证过的用户信息
	c.Set("claims", claims)
	return
}
