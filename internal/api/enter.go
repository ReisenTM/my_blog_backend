package api

import "blog/internal/api/user_api"

type Api struct {
	UserApi user_api.UserApi
}

// App 实例化 以供外部调用Api
var App = Api{}
