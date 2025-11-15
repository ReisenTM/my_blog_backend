package api

import (
	"blog/internal/api/article_api"
	"blog/internal/api/user_api"
)

type Api struct {
	UserApi    user_api.UserApi
	ArticleApi article_api.ArticleApi
}

// App 实例化 以供外部调用Api
var App = Api{}
