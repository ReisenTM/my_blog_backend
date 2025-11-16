# 博客服务 API 文档

> 最近更新：2025-11-16

## 1. 服务概览
- **基础地址**：默认从 `settings-dev.yaml` 读取，开发环境为 `http://127.0.0.1:8083`.
- **统一前缀**：所有业务接口都挂载在 `/api` 路由组之下（例如 `/api/auth/email-code`）。
- **静态资源**：`/uploads` 目录通过 `r.Static("/uploads","uploads")` 暴露，可直接以 `/uploads/<文件>` 访问已上传文件。

## 2. 通用交互约定
### 2.1 请求规范
- 建议所有请求携带 `Content-Type: application/json` 并使用 UTF-8 编码。
- 需要登录态的接口通过 `token` 头（或 `?token=` 查询参数）传递 JWT；Token 来源于邮箱登录接口。

### 2.2 统一响应格式
所有接口返回 HTTP 200 状态码，具体成功/失败通过 JSON 体中的 `code` 字段确认：

```json
{
  "code": 0,
  "message": "成功或错误提示",
  "data": {}
}
```

| code | 含义   |
|------|--------|
| 0    | 成功   |
| -1   | 业务失败（参数、鉴权、服务异常等） |

### 2.3 JWT 认证
- 登录成功后会得到一个字符串格式的 Token，内部载荷 `userID、username、role`。
- Token 默认有效期等于 `jwt.expire`（小时），鉴权中间件会校验是否过期或被拉黑。
- 黑名单记录存放在 Redis，以 `token_black_<token>` 为 key，常出现在用户注销或被管理员踢下线时。

### 2.4 邮箱验证码
- `/api/auth/email-code` 会生成 4 位数字验证码，通过邮件发送，默认 10 分钟有效。
- 验证码以邮箱为 key 缓存在内存 `email_store` 中，验证通过后立即删除，使用一次即失效。
- 邮箱注册接口由 `EmailVerifyMiddleware` 自动验证 `email + emailCode` 是否匹配，无需重复调用验证接口。

## 3. API 详细说明

### 3.1 发送邮箱验证码
- **方法 / 路径**：`POST /api/auth/email-code`
- **是否鉴权**：否
- **用途**：向目标邮箱发送注册或重置密码验证码。

#### 请求体
| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| `email` | string | 是 | `"user@example.com"` | 收码邮箱，系统会检查是否存在/不存在来决定是否允许发送。 |
| `type` | uint8 | 否（默认注册） | `1` | `1` 注册验证码，请求会拒绝已存在的邮箱；`2` 重置密码验证码，仅允许邮箱注册或绑定邮箱的账号。 |

#### 成功响应
```json
{
  "code": 0,
  "message": "验证码已发送",
  "data": {}
}
```

#### 常见失败
- `邮箱已存在，请登录`：type=1 时目标邮箱已经注册。
- `该邮箱不存在`、`仅支持邮箱注册或绑定邮箱的用户重置密码`：type=2 时的校验失败。
- `发送邮件失败,xxx`：SMTP 配置异常。

#### 示例
```bash
curl -X POST http://127.0.0.1:8083/api/auth/email-code \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","type":1}'
```

### 3.2 邮箱注册
- **方法 / 路径**：`POST /api/auth/register`
- **是否鉴权**：否，但请求会先经过 `EmailVerifyMiddleware`，自动验证邮箱验证码。
- **用途**：以邮箱 + 验证码 + 密码完成账号创建，并随机生成昵称（形如 `用户1234`）。

#### 请求体
| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| `email` | string | 是 | `"user@example.com"` | 需与验证码发送时的邮箱一致。 |
| `emailCode` | string | 是 | `"8456"` | 最近一次发送的 4 位验证码，仅可使用一次。 |
| `password` | string | 是 | `"Passw0rd!"` | 密码会使用 `bcrypt` 哈希后入库。 |

#### 成功响应
```json
{ "code": 0, "message": "注册成功", "data": {} }
```

#### 常见失败
- `邮箱验证失败`：验证码错误/过期/重复使用。
- `用户创建失败`：数据库写入报错，例如唯一键冲突。

#### 示例
```bash
curl -X POST http://127.0.0.1:8083/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","emailCode":"8456","password":"Passw0rd!"}'
```

### 3.3 邮箱密码登录
- **方法 / 路径**：`POST /api/auth/email-login`
- **是否鉴权**：否
- **用途**：以邮箱 + 密码换取 JWT，用于后续需要登录态的接口。

#### 请求体
| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| `email` | string | 是 | `"user@example.com"` | 必须是已注册邮箱。 |
| `password` | string | 是 | `"Passw0rd!"` | 明文密码，服务端与哈希对比。 |

#### 成功响应
```json
{
  "code": 0,
  "message": "成功",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```

#### 常见失败
- `请求体结构错误`：JSON 无法解析。
- `邮箱未注册`、`密码错误`：凭证错误。
- `生成token失败`：JWT 配置缺失或密钥异常。

#### 示例
```bash
curl -X POST http://127.0.0.1:8083/api/auth/email-login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"Passw0rd!"}'
```

### 3.4 发布文章
- **方法 / 路径**：`POST /api/post`
- **是否鉴权**：是，需通过 `AuthMiddleware`，在 Header 中携带 `token`.
- **用途**：创建文章草稿/正式文章，自动生成 `slug`（UUID）及 50 字摘要。

#### 请求体
| 字段 | 类型 | 必填 | 示例 | 说明 |
|------|------|------|------|------|
| `title` | string | 是 | `"Go 语言中的 JWT 实践"` | 最终展示的文章标题。 |
| `categories` | string[] | 否 | `["后端","Go"]` | 会以 JSON 数组存储/返回。 |
| `tags` | string[] | 否 | `["JWT","安全"]` | 同上。 |
| `content` | string | 否 | `"文章正文..."` | 若为空将导致 `summary` 为空；若有内容，摘要取前 50 个字符。 |

#### 成功响应
```json
{ "code": 0, "message": "发布文章成功", "data": {} }
```
数据库内会记录：
- 随机生成的 `slug`（UUID）。
- `authorId` = 当前登录用户的 `username`（注意：而非用户 ID）。
- `summary` 由 `content` 截断得到。

#### 常见失败
- `请求结构错误`：JSON 不合法或缺少 `title`。
- `请登录` / `token过期` 等：Token 缺失、无效或被拉黑。
- `发布文章失败`：数据库写入异常。

#### 示例
```bash
curl -X POST http://127.0.0.1:8083/api/post \
  -H 'Content-Type: application/json' \
  -H 'token: <登录获得的JWT>' \
  -d '{
        "title":"Go 语言中的 JWT 实践",
        "categories":["后端","Go"],
        "tags":["jwt","安全"],
	        "content":"JWT 是一种开放标准..."
	      }'
```

### 3.5 文章总览列表
- **方法 / 路径**：`GET /api/articles`
- **是否鉴权**：否
- **用途**：分页查看文章总览信息（标题、摘要、分类、作者与时间等），支持关键词与作者过滤。

#### 查询参数
| 参数 | 类型 | 必填 | 默认 | 说明 |
|------|------|------|------|------|
| `page` | int | 否 | `1` | 页码，最小 1。 |
| `pageSize` | int | 否 | `10` | 每页数量，最大 50。 |
| `keyword` | string | 否 | - | 按标题或摘要模糊搜索。 |
| `author` | string | 否 | - | 指定 `authorId` 的文章。 |

#### 成功响应
```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "slug": "c0c7c86e-5d27-45f6-ab36-3c35d6f756f3",
        "title": "Go 语言中的 JWT 实践",
        "summary": "JWT 是一种开放标准...",
        "categories": ["后端","Go"],
        "tags": ["jwt","安全"],
        "authorId": "用户1234",
        "createdAt": "2025-11-16 12:00:00",
        "updatedAt": "2025-11-16 12:05:00"
      }
    ],
    "count": 42
  }
}
```

#### 常见失败
- `查询参数错误`：Query 参数格式不符合要求。
- `统计文章数量失败` / `获取文章列表失败`：数据库查询异常。

#### 示例
```bash
curl "http://127.0.0.1:8083/api/articles?page=1&pageSize=5&keyword=Go"
```

## 4. 推荐业务流程
1. 调用 **发送邮箱验证码** 获取注册验证码。
2. 在验证码有效期内（10 分钟）调用 **邮箱注册** 完成开户。
3. 使用新账号调用 **邮箱密码登录** 获取 Token。
4. 携带 Token 调用 **发布文章** 等需要鉴权的接口。

## 5. 其他注意事项
- 邮件内容和发送账号取决于 `settings-*.yaml` 的 `email` 配置，生产环境务必替换默认凭证。
- 除文章发布外暂未开放更多业务接口，如需扩展请在 `internal/routers` 中新增路由并更新本文档。
- 如果未来提供密码重置、文章查询等接口，请沿用本文档的响应约定并补充示例，确保测试环境与文档同步。
