# Reisen`s Blog 后端实现指导

本文档汇总前端所需的核心 API、数据模型与集成要点，帮助你快速实现与当前 React/TypeScript 前端兼容的后端服务。

## 1. 概览
- **框架自由**：可选 Node.js/Express、Go、NestJS、Spring Boot 等，只需满足 REST 兼容。
- **数据源**：短期可用 JSON/Mock，长期建议对接数据库或 CMS。
- **返回格式**：统一使用 JSON，编码 UTF-8。

## 2. 路由与响应规范
| Method | Path | 说明 | 备注 |
| --- | --- | --- | --- |
| GET | `/api/posts` | 列出所有文章（首页/列表页使用）。 | 支持 `?category=` `?tag=` `?limit=` 过滤。 |
| GET | `/api/posts/:slug` | 获取单篇文章详情。 | 需包含 `sections`、`recommendedSlugs`。 |
| GET | `/api/topics` | 获取分类主题及文章统计。 | 与前端 `topics` 模块对应。 |
| GET | `/api/changelog` | 获取最新写作/更新日志。 | Home 页面底部使用。 |
| GET | `/api/authors/:id` | 作者信息。 | 可在文章详情中内联。 |
| POST | `/api/posts` (可选) | 后台发布文章。 | 需鉴权，可后续实现。 |
| POST | `/api/auth/email-code` | 发送邮箱验证码。 | 请求体 `{ email: string }`。 |
| POST | `/api/auth/verify` | 校验验证码并返回 token。 | 请求体 `{ email: string; code: string }`；响应携带 JWT/Session。 |
| GET | `/api/auth/github` | GitHub OAuth 跳转。 | 重定向到 GitHub 授权。 |
| GET | `/api/auth/github/callback` | GitHub OAuth 回调。 | 根据 `code` 换取 token，返回前端。

> 未来如需登录/交互，可新增 `/api/session`、`/api/comments` 等路由。

## 3. 数据模型
与 `src/types/blog.ts` 对应：
```ts
interface BlogPost {
  id: string;
  slug: string; // URL 标识
  title: string;
  summary: string;
  categories: string[];
  tags: string[];
  publishedAt: string; // ISO8601
  readingTime: string; // e.g. "8 min"
  hero: { badge: string; headline: string; description: string };
  authorId: string;
  sections: PostSection[];
  recommendedSlugs: string[];
}

interface PostSection {
  id: string;
  title: string;
  level: number; // 2=H2, 3=H3
  paragraphs: string[];
  list?: string[];
  alert?: { type: 'tip' | 'warning'; content: string };
  code?: { language: string; content: string }; // 用于 CodeBlock
}
```

### 示例响应：`GET /api/posts/golang-markdown-testing`
```json
{
  "id": "post-004",
  "slug": "golang-markdown-testing",
  "title": "Golang Markdown 代码块渲染测试",
  "summary": "在这篇随笔里，我用 Go 写了一个极简的 Markdown 渲染器…",
  "categories": ["后端编程", "编程随笔"],
  "tags": ["Go", "Markdown", "CLI"],
  "publishedAt": "2025-03-18",
  "readingTime": "5 min",
  "hero": {
    "badge": "Go · CLI Tooling",
    "headline": "用 Go 解析 Markdown 代码块",
    "description": "为了验证 CodeBlock 的表现，我写了段 Go 代码…"
  },
  "authorId": "gemini",
  "sections": [
    { "id": "motivation", "title": "写这段工具的小动机", "level": 2, "paragraphs": ["…"] },
    {
      "id": "code-demo",
      "title": "核心代码片段",
      "level": 2,
      "paragraphs": ["下面的代码展示了…"],
      "code": { "language": "go", "content": "package main\n..." }
    }
  ],
  "recommendedSlugs": ["dockerfile-best-practice"]
}
```

## 4. 鉴权与登录 API 示例

### 4.1 邮箱验证码
**请求：** `POST /api/auth/email-code`
```json
{
  "email": "s3068272040@gmail.com"
}
```
**响应：**
```json
{ "code": 0, "message": "验证码已发送" }
```

**请求：** `POST /api/auth/verify`
```json
{
  "email": "s3068272040@gmail.com",
  "code": "184203"
}
```
**响应：**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "jwt-or-session-token",
    "user": { "id": "reisen", "name": "Reisen", "avatar": "https://cdn/avatar.png" }
  }
}
```

### 4.2 GitHub OAuth
1. 前端点击 `/api/auth/github`，后端重定向 GitHub 并在回调中携带 `code`。
2. 回调地址 `/api/auth/github/callback?code=xxx`：
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "jwt-or-session-token",
    "user": {
      "id": "github-123",
      "name": "Reisen",
      "avatar": "https://avatars.githubusercontent.com/u/..."
    }
  }
}
```

> 建议所有响应都包含 `code`（0 表示成功）、`message`、`data`，方便前端统一处理。

## 5. 开发流程建议
1. **Mock 阶段**：将 `src/data/mockData.ts` 中数据迁移到后端 JSON 文件或数据库。
2. **接口联调**：前端用 `fetch`/React Query 调用上述 API，并用 loading/error 状态兜底。
3. **部署**：后端部署到云服务（Vercel Serverless、Netlify Functions、Render、Railway 等）。
4. **文档**：保持与本文件同步更新，尤其是新增字段、鉴权策略时。

## 6. Checklist
- [ ] 所有 API 采用 HTTPS。
- [ ] 响应头 `Content-Type: application/json; charset=utf-8`。
- [ ] 支持 CORS（允许前端域名或 `*`）。
- [ ] 文章内容中包含代码、提示、列表等结构时，后端直接返回统一格式，前端无需再次解析 Markdown。

如需我基于该文档生成后端代码，请明确后端技术栈、持久化方案及鉴权需求。EOF
