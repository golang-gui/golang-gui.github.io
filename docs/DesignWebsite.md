# GOUI 展示网站设计

## 目标与范围

以 v0.1.0 为首个展示版本，帮助搜索 Go／Golang GUI 的开发者了解 GOUI、
体验真实 Demo、编写第一个应用。网站独立于框架仓库，不修改框架实现。

公开仓库：`golang-gui/golang-gui.github.io`。
正式地址：`https://golang-gui.github.io/`。

展示内容遵循网站设计技能的建议：以真实截图、具体能力、可运行示例和可执行的
下一步形成完整阅读路径，不使用虚构评价、用户数量或未经验证的性能比较。
托管只使用用户指定的 GitHub Pages，不使用 Sites 托管。

## 页面与维护

| 路径                           | 作用                                             |
| ------------------------------ | ------------------------------------------------ |
| `/`                            | 英文展示首页                                     |
| `/zh-cn/`                     | 中文展示首页                                     |
| `/getting-started/`           | 英文入门教程                                     |
| `/zh-cn/getting-started/`     | 中文入门教程                                     |
| `/404.html`                    | 明确的找不到页面说明，不伪装成成功页面           |

网页以仓库中直接维护的静态 HTML 为主，CSS 负责布局和主题，少量原生 JavaScript
增强交互。英文、中文首页分别为根目录 `index.html` 与 `zh-cn/index.html`；英文教程
为 `getting-started/index.html`，中文教程为 `zh-cn/getting-started/index.html`。
所有页面从站点根目录组织，移除 `goui/` 网页目录。
GitHub Pages 直接发布已提交的 HTML、CSS、JS 和图片，不需要页面生成或前端构建。

Go 标准库生成器保留为可选的多语言维护工具：模板和翻译变动后手动运行，默认
更新仓库中的 HTML，审查后一起提交。直接修改 HTML 时也同步对应模板／翻译，
避免下次运行工具覆盖修改。生成器不参与预览或部署；更多语言需要追加翻译、
路由和语言配对，再生成并提交独立 HTML 页面。

正文、截图、导航、下载入口、FAQ 和教程代码都直接写入 HTML。浏览器不需要
执行 JavaScript 或请求内容接口才能阅读页面。404 页不加载 JavaScript；展示首页
（含根目录）和教程按需加载原生脚本，仅增强截图轮播和复制代码。禁用脚本后仍能
阅读完整内容，三平台截图逐张展示，FAQ 使用原生 details／summary。

静态 HTML 文件直接加入版本控制。Pages 配置为 Deploy from a branch，来源为
`main / (root)`，仓库根目录就是网站根目录；`.nojekyll` 保持原文件静态发布。
不再维护发布文件复制清单或自定义打包部署工作流，也不运行 Go 生成器。
独立的 Check website 工作流校验 main 推送及 PR，校验不作为分支发布的前置条件。

## 内容与视觉

- 沿用 GOUI 原色 Logo 与三平台截图，不生成或重绘品牌图形。
- Hero 只保留主标题、正常字号的副标题与两个入口；No cgo 和 MIT 放在副标题，
  中文按钮为“快速开始”和“体验Demo”。版本信息放在 Demo 下载区，不使用徽章、
  细字说明、统计条或各区块标题前的小标题；内容列表不增加装饰编号。
- 视觉参照本地 `goui/theme/modern/modern.go` 与 `doc/DesignModernPresets.md`：
  亮色页面使用 `#F6F6F8` 背景、白色内容表面、`#202024` 正文与中性细描边；
  按钮 6px 圆角，普通内容卡片 10px 圆角，不以阴影表达层次。主要操作使用
  Modern 默认蓝色主题色的对比度调整值，选中平台使用中性表面。
- 扩大桌面内容宽度至 1320px，把更多首页空间分配给截图；1080px 以下上下排列。
  三平台统一使用 824:779 的布局区域，以 contain 保留完整原图，macOS 的较矮截图
  居中展示，不裁切、不拉伸。不绘制截图外层卡片、边框、背景或额外圆角，只保留
  截图原有的窗口轮廓；平台切换与轮播控件集中在图下，区域尺寸不随平台变化。
- 默认 Windows，每 5.5 秒轮播至 macOS、Linux，再回到 Windows。提供平台切换、
  上一张／下一张、暂停／播放、方向键／Home／End 和触摸横划；手动操作才播报
  截图变化。悬停、键盘焦点、页面隐藏和滚出视口时停止自动切换，reduced-motion
  默认不自动播放且不执行切换动画。禁用 JS 仍可阅读全部三张截图与原图链接。
- 用实际功能组织能力介绍，避免只列抽象营销口号。
- 展示已发布资源的实际字节大小；macOS ZIP 大小不是可执行文件大小。
- 若使用已有开发测量，必须明确 Go 版本、编译参数、OS、后端和测量范围，不能冒充发布产物。
- 明确 Linux X11／运行库要求、macOS 签名现状、支持的架构和 v0.1.0 的 API 演进状态。
- 使用系统字体、响应式布局、键盘可操作链接与控件、可见焦点和 reduced-motion。
- Logo 保留 Renee French／CC BY 4.0 的来源归属，不把全部素材简单声明为 MIT。
- nanovgo 和 oksvg 作为自维护的绘图依赖在 GOUI 页尾简要说明，不设独立组织入口。

## SEO

- 标题和自然正文覆盖 Go GUI、Golang GUI、desktop applications、no cgo 等真实搜索意图。
- 正式页面有独立 title／description、自引用 canonical、清晰标题和正常的 a 链接；
  英文首页 canonical 为 `/`，中文首页为 `/zh-cn/`，四个正式 URL 均列入 sitemap。
- 中英文使用独立 URL；对应页面互相声明 en／zh-CN／x-default，不按浏览器语言强制跳转。
- 提供根目录 robots.txt 与只包含正式页面的 sitemap.xml；不填充虚构 lastmod。
- 提供真实的 SoftwareSourceCode JSON-LD 和所属组织信息，不虚构评分或承诺搜索富结果。
- 配置 Open Graph、分享图片、描述性 alt、图片尺寸；避免外部字体和第三方追踪脚本。
- 网站完整生成后检查所有站内路径、锚点、元数据、结构化数据、语言配对及图片。

搜索收录和排名不是部署验收条件，也不能保证。Google Search Console 与 Bing
Webmaster Tools 需要站点所有者账号；没有真实验证码时不输出假验证标签。
部署完成后由所有者提供 HTML 文件或 meta 验证值，再提交 sitemap 和检查 URL。
README／About／发布说明反链单独授权，不在本次改动框架仓库。

参考：

- [GitHub Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/what-is-github-pages)
- [Pages 发布来源配置](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
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
- [x] 推送并启用 Pages，确认部署结果与线上返回内容。
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

### 2026-10-14 发布与线上检查

- 网站设计 `e601d56` 与源码 `61628e7` 已推送；首次源码推送因 OAuth 缺少
  workflow scope 被拒绝，没有绕过权限。所有者明确授权后完成权限刷新并正常推送。
- 授权刷新前两次因 CLI 的 TLS 请求超时未保存；第三次用仅作用于该进程的
  连接兼容设置完成，没有修改系统网络或全局 Git 配置。这是工具执行记录，
  不是网站运行要求。
- Pages 发布来源为 workflow，HTTPS 已启用。源码对应的
  [Actions 运行](https://github.com/golang-gui/golang-gui.github.io/actions/runs/38527060393)
  中，测试、构建、上传和部署均成功。
- 五个正式页面、robots.txt、sitemap.xml 的线上 GET 均返回 200，响应正文
  与本地生成文件逐字一致；不存在的路径返回真正的 404，并使用自定义说明页。
- 线上检查不等于搜索引擎收录。Google／Bing 账号验证及 sitemap 提交仍未执行；
  具体接入见 [SearchSetup](SearchSetup.md)，不保证泛词 golang gui 的排名。
- 本地浏览器验证只覆盖 Chromium；未在 Safari／Firefox 或远端移动设备复验，
  没有将跨平台截图或示例编译算作这些浏览器通过。

### 2026-10-14 Modern 样式与截图轮播验收

- `go test ./...`、静态网站生成、`node --check assets/site.js` 通过。
- Chromium 检查五个正式页面，各使用 1920／1440／1080／1024／768／760／390／320px
  宽度，共 40 组，无页面横向溢出、脚本错误或资源错误。调整英文首屏标题字号后，
  追加检查中英文首页上述宽度及 1100px，共 18 组，结果通过。
- 中英文首页在各宽度下逐一切换 Windows、macOS、Linux，外卡片与图片展示框
  宽高均一致，图片以 contain 完整展示且没有超出原生尺寸。1440px 宽度下，
  卡片约 802×860px，展示框约 768×726px；macOS 截图在同一个框内垂直居中。
- 自动播放三平台循环、平台按钮、前后切换、暂停／播放、方向键／Home／End、
  悬停／聚焦／页面隐藏／滚出视口暂停、模拟触摸横划及防止横划误开原图均通过。
  reduced-motion 默认暂停并禁用动画；用户主动播放与手动切换正常。
- 禁用 JavaScript 后中英文全部三张截图仍可见且展示框尺寸一致；跳过导航及复制
  代码仍正常。测试工具与截图保存在仓库外，不增加网站运行或构建依赖。
- 本轮为本地修改与 Chromium 验证，未推送发布；Safari／Firefox 和真实触摸设备
  未验证。

### 2026-10-14 静态 HTML 验收

- 页面以预生成 HTML 和共享 CSS 承载完整内容，不使用前端框架；Go 仅参与构建。
  组织首页与 404 页不加载交互脚本，展示首页和教程只加载原生 JavaScript。
- `go test ./...` 与静态网站生成通过。Chromium 禁用 JavaScript 后检查六个页面，
  分别使用 1440px 与 390px 宽度，共 12 组：正文、导航、三平台截图、下载链接、
  教程代码均可读，无横向溢出，原生 FAQ 可展开。
- 启用 JavaScript 后核对各页面脚本请求：首页轮播和教程复制按钮正常；组织首页
  与 404 页没有 JavaScript 资源请求。本轮未推送发布。

### 2026-10-14 HTML 文件直接维护与发布

- 六个完整 HTML 文件已放到仓库根目录及对应语言／教程目录，可以直接编辑并
  从仓库根目录预览。`.nojekyll`、robots 和 sitemap 与静态页面一起维护。
- 保留多语言生成器作为可选工具，默认输出到仓库根目录；`go test ./...` 通过，
  原地刷新验证通过，页面及图片／样式内容均保持一致。
- Actions 不安装 Go、不生成 HTML，只运行静态文件校验并复制发布文件。发布包
  为 16 个静态文件，不包含模板、翻译源、Go 模块或生成器代码。
- `python3 scripts/check_site.py`、JavaScript 语法与空白检查通过；隔离目录中
  人为制造的错误链接被静态校验拒绝。
- Chromium 分别从仓库根目录和静态发布包检查六个页面，在 1440px 与 390px
  下验证禁用 JavaScript 的阅读和导航，共 24 组；FAQ、三平台截图、下载链接、
  教程代码正常。启用脚本后的轮播和教程增强正常，无资源或脚本错误。
- 本轮修改和验收在本地完成，尚未推送，未宣称新的工作流已在线执行。

### 2026-10-14 页面元素简化

- 根目录直接展示完整 GOUI 英文介绍，不跳转、不再展示组织门户；nanovgo 与
  oksvg 在英文、中文介绍页的页尾简要说明。根目录与 `/goui/` HTML 内容一致，
  canonical 保持 `/goui/`，sitemap 去重后包含四个正式语言／教程 URL。
- 截图区域去掉两层表面、边框、背景和所有额外圆角，只使用原图的窗口轮廓；
  原有三平台等尺寸布局与轮播保留，控制条放在截图下方。
- Hero 仅保留标题、副标题与两个按钮。No cgo 和 MIT 放入正常字号的副标题，
  中文按钮为“快速开始”和“体验Demo”；移除版本徽章、细字说明、横线统计条。
  版本说明放在下载区，不在各下载卡片重复展示版本。
- 删除各大标题前的小标题，以及内容列表的装饰编号；删除旧组织页模板、未使用
  的文案与样式。静态 HTML 和可选生成器／中英文翻译同步更新。
- `go test ./...`、静态校验、JavaScript 语法与空白检查通过。Chromium 在六个
  页面、1440／1024／768／390／320px 宽度下检查，共 30 组，无横向溢出；各平台
  截图区域尺寸相同、未超过原生尺寸、外层无描边／背景／圆角。
- 自动轮播、暂停／播放的可访问名称、前后切换、键盘操作、复制、reduced-motion
  和禁用 JavaScript 的阅读通过。中文副标题调整后复查 1440／390／360／320px，
  修正手机末尾孤字换行。已复核桌面、手机和页尾截图，本轮未推送发布。

### 2026-10-14 语言目录移到站点根目录

- 英文首页为 `/`，中文首页为 `/zh-cn/`；教程为 `/getting-started/` 和
  `/zh-cn/getting-started/`。移除旧 `goui/` 网页目录，更新静态 HTML、可选生成器、
  发布清单、导航、语言切换、canonical、hreflang、JSON-LD 与 sitemap。
- `go test ./...` 和静态校验通过，回归检查确保生成器不再创建 `goui/` 目录。
  本地旧预览保存在仓库外，当前 `_site` 仅复制新的根目录布局。
- Chromium 检查五个页面的桌面／手机布局，共 10 组；语言、首页、教程、Demo
  入口、轮播及禁用 JavaScript 后的导航通过，没有旧网页路径或资源错误。
  本轮未推送发布。

### 2026-10-14 网页替换透明截图

- 网页三平台图片替换为已预览的透明版本：Linux、Windows 仅清理四角背景，
  macOS 裁掉边缘壁纸并清理四角，原窗口内容保持不变。macOS 原生尺寸更新为
  820×634，中英文静态 HTML 与本地 `_site` 已同步。
- 保留同尺寸展示区域与原生轮播。静态文件校验、`go test ./...` 和空白检查通过；
  Chromium 检查中英文首页在 1440px、390px 下的三平台切换，展示框尺寸一致，
  无横向溢出或资源错误，自动轮播正常。
- 页面预览保存在 `preview/screenshot-cleanup/page-*-desktop.png` 和
  `page-mobile.png`。仅替换网页资源，GOUI 框架仓库截图与 README 未修改，
  本轮未推送发布。

### 2026-10-14 修正 macOS 窗口细边缘

- 上一版裁切贴近窗口内沿，移除了与壁纸混色的外侧细边缘。依据原图外沿的
  明暗比例估计透明度，恢复 1px 中性灰色轮廓，保留内侧高光与窗口实体像素。
  修正版为 822×636，四角及外侧壁纸仍透明，无额外网页卡片或 CSS 边框。
- 网页资源、中英文 HTML 尺寸与 `_site` 同步更新。静态校验通过；Chromium
  检查桌面、手机上的三平台等尺寸展示区域与轮播，未发现溢出或资源错误。
- 修正版与处理报告保存在 `demo-macos-transparent-v2.png` 和
  `report-macos-v2.json`，实际页面效果为 `page-macos-desktop-v2.png`、
  `page-macos-mobile-v2.png`。GOUI 原仓库与 README 未修改，未推送发布。

### 2026-10-27 更新 GOUI v0.4.0

- 中英文首页与教程跟随正式发布的 v0.4.0；安装命令、API／文档链接、FAQ、
  JSON-LD 版本及 404 页页尾同步更新，保留现有布局、原色 Logo 与三平台截图。
- `content/release.json` 的发布日期为 2026-10-26；三平台下载文件名、URL、
  字节大小和 SHA-256 与 GitHub 正式 Release 元数据及公开 SHA256SUMS 核对一致。
  页面显示 Windows EXE 9.50 MiB、macOS .app ZIP 3.33 MiB、Linux 可执行文件
  9.34 MiB，仍区分压缩包大小与可执行文件大小，不重写历史开发内存测量。
- 功能介绍补充 Tooltip、菜单按钮组、长段落有界排版和声明式内容组合；没有将
  规划中的子菜单、原生系统菜单栏或其他未验收目标写为已支持能力。
- Linux amd64／Go 1.25.8：`CGO_ENABLED=0 go test ./... -count=1`、
  `go vet ./...`、`python3 scripts/check_site.py`、`node --check assets/site.js`
  及空白检查通过。翻译键一致，四个正式页面和 404 的本地链接／锚点／元数据通过。
- 首页与教程两份原始 Go 示例复制到仓库外模块，通过公共 proxy 获取 v0.4.0，
  不使用本地 replace；Linux amd64、Windows amd64、macOS arm64 的 no-cgo
  构建全部通过。未启动 GUI 示例，不将交叉编译算作原生窗口验收。
- 本机 Chrome 156.0.8078.47 无头浏览器检查五页的 1440／390px 布局，并追加
  四个正式页面禁用 JavaScript 的检查，共 18 组通过；没有页面级横向溢出、
  脚本异常或资源错误。下载入口、大小和日期正确，首页三平台截图切换正常，
  无脚本时三张截图均可读；浏览器结果及截图保存在仓库外临时目录。
- 本轮只更新本地网站仓库，尚未提交、推送或触发 Pages 发布；未修改 GOUI 框架
  仓库，也未将本机 Chromium 检查当作 Safari／Firefox 或线上发布验收。
