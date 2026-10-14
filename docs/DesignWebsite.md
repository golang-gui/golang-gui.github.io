# GOUI 展示网站设计

## 目标与范围

以 v0.1.0 为首个展示版本，帮助搜索 Go／Golang GUI 的开发者了解 GOUI、
体验真实 Demo、编写第一个应用。网站独立于框架仓库，不修改框架实现。

公开仓库：`golang-gui/golang-gui.github.io`。
正式地址：`https://golang-gui.github.io/goui/`。

展示内容遵循网站设计技能的建议：以真实截图、具体能力、可运行示例和可执行的
下一步形成完整阅读路径，不使用虚构评价、用户数量或未经验证的性能比较。
托管只使用用户指定的 GitHub Pages，不使用 Sites 托管。

## 页面与维护

| 路径                           | 作用                                             |
| ------------------------------ | ------------------------------------------------ |
| `/`                            | golang-gui 组织入口，链接 GOUI、nanovgo 和 oksvg |
| `/goui/`                       | 英文展示首页，主要搜索入口                       |
| `/goui/zh-cn/`                 | 中文展示首页                                     |
| `/goui/getting-started/`       | 英文入门教程                                     |
| `/goui/zh-cn/getting-started/` | 中文入门教程                                     |
| `/404.html`                    | 明确的找不到页面说明，不伪装成成功页面           |

Go 标准库小工具复用模板和中英文 JSON 内容，生成 `_site` 中的静态文件。
生成器没有第三方 Go 依赖，不引入 npm 或前端框架；浏览器不需要执行 JavaScript
才能获取正文、截图、链接和教程。JavaScript 仅增强复制代码与截图切换。

生成目录不提交。main 推送触发 GitHub Actions；PR 只检查，不部署。
Pages 使用官方 configure-pages、upload-pages-artifact、deploy-pages 操作。

## 内容与视觉

- 沿用 GOUI 原色 Logo 与三平台截图，不生成或重绘品牌图形。
- 首屏明确项目定位、当前版本、入门和 Demo 入口。
- 截图不放大超出原生尺寸，默认 Windows，其他平台可切换；禁用 JS 仍可阅读全部截图。
- 用实际功能组织能力介绍，避免只列抽象营销口号。
- 展示已发布资源的实际字节大小；macOS ZIP 大小不是可执行文件大小。
- 若使用已有开发测量，必须明确 Go 版本、编译参数、OS、后端和测量范围，不能冒充发布产物。
- 明确 Linux X11／运行库要求、macOS 签名现状、支持的架构和 v0.1.0 的 API 演进状态。
- 使用系统字体、响应式布局、键盘可操作链接与控件、可见焦点和 reduced-motion。
- Logo 保留 Renee French／CC BY 4.0 的来源归属，不把全部素材简单声明为 MIT。

## SEO

- 标题和自然正文覆盖 Go GUI、Golang GUI、desktop applications、no cgo 等真实搜索意图。
- 每页有独立 title／description、自引用 canonical、清晰标题和正常的 a 链接。
- 中英文使用独立 URL；对应页面互相声明 en／zh-CN／x-default，不按浏览器语言强制跳转。
- 提供根目录 robots.txt 与只包含正式页面的 sitemap.xml；不填充虚构 lastmod。
- 提供真实的 Organization／SoftwareSourceCode JSON-LD，不虚构评分或承诺搜索富结果。
- 配置 Open Graph、分享图片、描述性 alt、图片尺寸；避免外部字体和第三方追踪脚本。
- 网站完整生成后检查所有站内路径、锚点、元数据、结构化数据、语言配对及图片。

搜索收录和排名不是部署验收条件，也不能保证。Google Search Console 与 Bing
Webmaster Tools 需要站点所有者账号；没有真实验证码时不输出假验证标签。
部署完成后由所有者提供 HTML 文件或 meta 验证值，再提交 sitemap 和检查 URL。
README／About／发布说明反链单独授权，不在本次改动框架仓库。

参考：

- [GitHub Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/what-is-github-pages)
- [Pages 自定义工作流](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)
- [Google SEO 指南](https://developers.google.com/search/docs/fundamentals/seo-starter-guide)
- [Google 多语言页面](https://developers.google.com/search/docs/specialty/international/localized-versions)
- [Bing 接入流程](https://www.bing.com/webmasters/help/getting-started-checklist-66a806de)
- [Go Gopher 来源与许可](https://go.dev/blog/gopher)

## 验证与进度

- [x] 确认组织管理员权限、仓库名及 v0.1.0 发布资源。
- [x] 创建独立公开仓库，固化本文。
- [x] 实现五个正式页面、404、素材、sitemap 和 robots。
- [x] 生成器／结构校验测试、Go 入门代码对发布 tag 编译验证。
- [x] 桌面与移动端浏览器检查，禁用 JS、截图切换和复制代码验证。
- [ ] 推送并启用 Pages，确认部署结果与线上返回内容。
- [ ] Google／Bing 所有权验证和 sitemap 提交，需要所有者账号。

具体结果在交付时补充，区分本地验证、线上部署、搜索工具接入和实际收录。

### 2026-10-14 本地验证

- Linux amd64、Go 1.25.8：`CGO_ENABLED=0 go test ./...`、生成器执行、
  `go vet ./...`、JavaScript 语法与空白检查通过。测试检查六个页面、所有站内
  链接／锚点、独立标题、canonical、JSON-LD、互相对应的语言链接和五个 sitemap URL。
- 初次 `CGO_ENABLED=0 go test -race ./...` 被 Go 拒绝：race 工具需要 cgo；
  不是用例失败，也不计作通过。之后在已有 C 编译器环境下以 `CGO_ENABLED=1`
  执行相同包的 race 测试通过；生产生成与常规测试仍不需要 cgo。
- 浏览器使用本机已有 Chromium，通过临时目录内的 Playwright 验证；工具依赖
  不加入网站构建。五个正式页面分别检查 1440／1024／768／390／320px 宽度，
  无页面级横向溢出。中英文三平台截图切换、图片不放大、复制按钮、键盘跳过
  导航和禁用 JS 后的正文／全部截图可读均通过，无浏览器脚本或资源错误。
- 初次 320px 检查发现代码区的网格子项最小宽度撑宽页面；修正子项 min-width
  后复测通过，没有通过隐藏整个页面溢出规避。截图复核后缩短中文版首屏标题，
  修正孤字换行，并复跑相同检查。
- 入门例程从页面使用的 `content/hello.go.txt` 原样复制到仓库外临时模块，
  依赖公开的 `goui v0.1.0`，不 replace 到工作区。Linux amd64、Windows amd64、
  macOS arm64 的 `CGO_ENABLED=0` 构建均通过，未运行该 GUI 例程，不计为三端窗口验收。
- 下载名称、大小与 SHA-256 来自实际 GitHub release；教程链接的三个仓库文档
  路径已在 v0.1.0 tag 的 tree 中核对。原色 SVG 和三张 PNG 仅复制，未改动原仓库。
- 本轮未修改 GOUI 源码或 README；设计与网站源码在独立仓库中维护。
