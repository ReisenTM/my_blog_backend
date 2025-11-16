# Reisen`s Blog · Go 后端

基于 **Go + Gin + Gorm** 的博客服务，提供文章管理、邮箱验证码登录/注册、JWT 鉴权等接口。配套的 React 前端位于 `../frontend`。

---

## ✨ 功能亮点

- **文章模块**：支持分页查询、关键词/作者过滤；登录后可创建文章，后端自动生成 `slug`、摘要并将 Markdown 渲染为 HTML。
- **用户鉴权**：邮箱验证码注册 + 密码登录，登录成功返回 JWT；中间件会校验黑名单、过期等状态。
- **统一响应**：所有接口返回 `code/message/data`，便于前端统一处理。
- **CLI 工具**：通过 `-t user -s create` 子命令可在命令行创建账号（支持角色选择）。

---

## ⚙️ 技术栈

| 领域 | 技术 |
| --- | --- |
| HTTP 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [Gorm](https://gorm.io) + PostgreSQL driver |
| JWT | 自定义 `internal/utils/jwts` |
| Markdown | `gomarkdown/markdown`（渲染 HTML、生成摘要） |
| 其他 | Redis、Logrus、Cron 等 |

---

## 🚀 快速开始

1. **克隆仓库并进入后端目录**
   ```bash
   cd backend/blog
   ```
2. **复制配置**（可直接使用示例配置或修改自己的参数）
   ```bash
   cp settings-dev.yaml settings.yaml
   ```
3. **安装依赖并运行**
   ```bash
   go mod tidy
   go run main.go -f settings.yaml
   ```
   程序启动后默认监听 `settings.yaml` 中的 `system.addr`，对外暴露 `/api` 路由。

### 数据库迁移

```
go run main.go -f settings.yaml -db
```

---

## 👤 命令行创建用户

命令行支持交互式创建账号，包含角色（普通用户 / 管理员）选择与密码确认。适合初始化管理员。

```bash
go run main.go -f settings.yaml -t user -s create
```

执行后按提示输入角色编号、用户名、密码即可。用户的 `RegSource` 会被标记为 “命令行创建”，并写入数据库。

---

## 📡 API 概览

| 方法 & 路径 | 说明 |
| --- | --- |
| `GET /api/posts` | 文章列表（分页/搜索/作者过滤） |
| `POST /api/post` | 发布文章（需 `token`） |
| `POST /api/auth/email-code` | 发送邮箱验证码（注册/重置） |
| `POST /api/auth/register` | 邮箱注册（验证码 + 密码） |
| `POST /api/auth/email-login` | 邮箱密码登录，返回 JWT |
| `GET /api/post/:slug` | 单篇文章详情（含 HTML 内容） |

完整字段定义可参考 `doc/api-reference.md`。

---

## 📁 目录结构

```
backend/blog
├─ internal/
│  ├─ api/              # 文章、用户等业务接口
│  ├─ middleware/       # JWT、日志、CORS
│  ├─ model/            # Gorm 模型与枚举
│  ├─ flags/            # 命令行子命令（迁移、创建用户）
│  └─ utils/            # JWT、密码、Markdown、邮箱等工具
├─ doc/                 # API 文档
├─ settings-*.yaml      # 配置示例
├─ main.go              # 入口，负责加载配置、初始化依赖
└─ README.md
```

---

## 🤝 开发建议

- 运行前请确认数据库/Redis 连接信息在 `settings.yaml` 中配置正确。
- 新增字段后需同步更新 `internal/model` 和相应的 API/响应结构。
- 如需扩展接口，可在 `internal/api` 中新增模块，并在 `internal/routers` 下注册路由。

欢迎结合前端一起体验完整的 “Reisen`s Blog”！ 🎉
