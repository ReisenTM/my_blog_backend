package article_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/custom"
	"blog/internal/utils/jwts"
	"blog/internal/utils/str"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostArticleReq struct {
	Title      string             `json:"title" binding:"required"`
	Categories custom.StringArray `json:"categories"`
	Tags       custom.StringArray `json:"tags"`
	Content    string             `json:"content"`
}

func (ArticleApi) PostArticle(c *gin.Context) {
	var req PostArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.FailWithMsg("请求结构错误", c)
		return
	}
	claims := jwts.GetClaims(c)
	newArticle := model.ArticleModel{
		Slug:       uuid.NewString(),
		Title:      req.Title,
		Summary:    str.Substr(req.Content, 50),
		Categories: req.Categories,
		Tags:       req.Tags,
		AuthorID:   claims.Username,
		Content:    req.Content,
	}
	err := global.DB.Create(&newArticle).Error
	if err != nil {
		resp.FailWithMsg("发布文章失败", c)
		return
	}
	resp.OKWithMsg("发布文章成功", c)
}
