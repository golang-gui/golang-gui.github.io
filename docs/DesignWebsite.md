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

| 路径 | 作用 |
| --- | --- |
| `/` | golang-gui 组织入口，链接 GOUI、nanovgo 和 oksvg |
| `/goui/` | 英文展示首页，主要搜索入口 |
| `/goui/zh-cn/` | 中文展示首页 |
| `/goui/getting-started/` | 英文入门教程 |
| `/goui/zh-cn/getting-started/` | 中文入门教程 |
| `/404.html` | 明确的找不到页面说明，不伪装成成功页面 |

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
- [ ] 实现五个正式页面、404、素材、sitemap 和 robots。
- [ ] 生成器／结构校验测试、Go 入门代码对发布 tag 编译验证。
- [ ] 桌面与移动端浏览器检查，禁用 JS、截图切换和复制代码验证。
- [ ] 推送并启用 Pages，确认部署结果与线上返回内容。
- [ ] Google／Bing 所有权验证和 sitemap 提交，需要所有者账号。

具体结果在交付时补充，区分本地验证、线上部署、搜索工具接入和实际收录。
