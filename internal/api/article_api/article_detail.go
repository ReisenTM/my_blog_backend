package article_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/custom"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleDetail struct {
	ID         uint               `json:"id"`
	Slug       string             `json:"slug"`
	Title      string             `json:"title"`
	Summary    string             `json:"summary"`
	Categories custom.StringArray `json:"categories"`
	Tags       custom.StringArray `json:"tags"`
	AuthorID   string             `json:"authorId"`
	Content    string             `json:"content"`
	CreatedAt  string             `json:"createdAt"`
	UpdatedAt  string             `json:"updatedAt"`
}

func (ArticleApi) ArticleDetail(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		resp.FailWithMsg("缺少文章标识", c)
		return
	}

	var article model.ArticleModel
	err := global.DB.Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.FailWithMsg("文章不存在", c)
			return
		}
		resp.FailWithMsg("获取文章详情失败", c)
		return
	}

	detail := ArticleDetail{
		ID:         article.ID,
		Slug:       article.Slug,
		Title:      article.Title,
		Summary:    article.Summary,
		Categories: article.Categories,
		Tags:       article.Tags,
		AuthorID:   article.AuthorID,
		Content:    article.Content,
		CreatedAt:  article.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  article.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	resp.OkWithData(detail, c)
}
