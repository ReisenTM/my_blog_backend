package article_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/custom"
	"blog/internal/utils/jwts"
	utilsMarkdown "blog/internal/utils/markdown"
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
	renderedHTML := utilsMarkdown.RenderMarkdown(req.Content)
	preview, _ := utilsMarkdown.GetPreviewContent(req.Content, 50)

	newArticle := model.ArticleModel{
		Slug:       uuid.NewString(),
		Title:      req.Title,
		Summary:    preview,
		Categories: req.Categories,
		Tags:       req.Tags,
		AuthorID:   claims.Username,
		Content:    renderedHTML,
	}
	err := global.DB.Create(&newArticle).Error
	if err != nil {
		resp.FailWithMsg("发布文章失败", c)
		return
	}
	resp.OKWithMsg("发布文章成功", c)
}
