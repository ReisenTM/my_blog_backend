package article_api

import (
	"blog/internal/common/resp"
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/custom"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	defaultArticlePage     = 1
	defaultArticlePageSize = 10
	maxArticlePageSize     = 50
)

// ArticleOverviewQuery 文章总览查询参数
type ArticleOverviewQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Author   string `form:"author"`
}

// ArticleOverviewItem 返回给前端的文章概览信息
type ArticleOverviewItem struct {
	ID         uint               `json:"id"`
	Slug       string             `json:"slug"`
	Title      string             `json:"title"`
	Summary    string             `json:"summary"`
	Categories custom.StringArray `json:"categories"`
	Tags       custom.StringArray `json:"tags"`
	AuthorID   string             `json:"authorId"`
	CreatedAt  string             `json:"createdAt"`
	UpdatedAt  string             `json:"updatedAt"`
}

// ArticleOverview 文章总览列表
func (ArticleApi) ArticleOverview(c *gin.Context) {
	var q ArticleOverviewQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		resp.FailWithMsg("查询参数错误", c)
		return
	}

	if q.Page <= 0 {
		q.Page = defaultArticlePage
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultArticlePageSize
	}
	if q.PageSize > maxArticlePageSize {
		q.PageSize = maxArticlePageSize
	}

	db := global.DB.Model(&model.ArticleModel{})
	if q.Keyword != "" {
		like := "%" + strings.TrimSpace(q.Keyword) + "%"
		db = db.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	if q.Author != "" {
		db = db.Where("author_id = ?", q.Author)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		resp.FailWithMsg("统计文章数量失败", c)
		return
	}

	var articles []model.ArticleModel
	err := db.Order("created_at DESC").
		Offset((q.Page - 1) * q.PageSize).
		Limit(q.PageSize).
		Find(&articles).Error
	if err != nil {
		resp.FailWithMsg("获取文章列表失败", c)
		return
	}

	list := make([]ArticleOverviewItem, 0, len(articles))
	for _, art := range articles {
		list = append(list, ArticleOverviewItem{
			ID:         art.ID,
			Slug:       art.Slug,
			Title:      art.Title,
			Summary:    art.Summary,
			Categories: art.Categories,
			Tags:       art.Tags,
			AuthorID:   art.AuthorID,
			CreatedAt:  art.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  art.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp.OkWithList(list, int(total), c)
}
