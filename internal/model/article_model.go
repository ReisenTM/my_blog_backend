package model

import (
	"blog/internal/model/custom"
	"gorm.io/gorm"
)

// ArticleModel 表示一篇完整的文章数据结构，包含元数据、封面信息和内容分段。
type ArticleModel struct {
	gorm.Model
	// 文章用于路由的短链接（slug）
	Slug string `json:"slug"`
	// 文章标题
	Title string `json:"title"`
	// 文章摘要，用于文章简介展示
	Summary string `json:"summary"`
	// 文章所属分类
	Categories custom.StringArray `json:"categories"`
	// 文章标签列表
	Tags custom.StringArray `json:"tags"`
	// 作者的唯一标识
	AuthorID string `json:"authorId"`
	// 正文内容（Markdown 或 HTML）
	Content string `json:"content"`
}
