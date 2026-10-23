# 搜索引擎接入与维护

页面具备可抓取的静态正文、canonical、中英文 hreflang、站点地图和 robots。
这些不是搜索收录或排名保证；部署成功不能标记为 SEO 验收完成。

## Google Search Console

1. 使用站点所有者的 Google 账号打开 Search Console。
2. 添加 URL 前缀资源 `https://golang-gui.github.io/`，不选择无法控制 DNS 的 github.io 域名资源。
3. 选择 HTML 文件或 HTML meta 验证方式；由所有者提供真实文件／验证值。
4. 将验证文件提交到仓库根目录并加入工作流的静态文件复制清单，或将 meta 加入对应 HTML；不要填入占位验证码。
5. 验证后提交 `https://golang-gui.github.io/sitemap.xml`。
6. 用 URL 检查工具查看 `/`、`/zh-cn/` 与中英文教程，确认抓取、canonical 和索引状态。

已加入所有者提供的 [Google 验证文件](../google7ae34c380db01faf.html)，
Pages 工作流将根目录的 `google*.html` 原样复制到发布目录。生成页面不会覆盖该文件。
提交文件不等于已完成账号验证；发布后仍需所有者在 Search Console 中确认。
若改用 meta 验证，同时更新可选生成器的模板，避免后续刷新 HTML 时丢失验证标签。

## Bing Webmaster Tools

1. 登录站点所有者账号，手动添加站点，或按 Bing 提示从已验证的 Search Console 导入。
2. 使用提供的 XML 文件或 meta 验证值完成所有权验证。
3. 提交同一 sitemap，检查 URL、抓取问题和搜索展示数据。
4. 后续需要 IndexNow 时另行增加真实 key 和提交步骤，首版不提前承诺自动提交。

已加入所有者提供的 [Bing 验证文件](../BingSiteAuth.xml)，Pages 工作流将其原样复制到
发布目录。发布后仍需所有者在 Webmaster Tools 中确认，文件不随 GOUI 版本更新而改变。
静态检查同时检查两种验证文件的格式，避免后续发布遗漏文件。

## 内容与外部链接

- 待网站部署确认后，可由所有者将官网链接加入 GOUI 仓库 About、README 和发布说明。
- 教程始终对应明确的发布版本；更新版本时先编译示例再改页面。
- 优先补充真实应用示例和实用教程，不批量生成近似的关键词页面。
- 不增加 meta keywords、隐藏关键词、虚构评价或不必要的第三方追踪脚本。
- 记录收录、展示、点击与搜索词，区分品牌词 GOUI、长尾词和泛词 golang gui。

参考：[Google SEO 指南](https://developers.google.com/search/docs/fundamentals/seo-starter-guide)、
[Google 多语言页面](https://developers.google.com/search/docs/specialty/international/localized-versions)、
[Bing 入门流程](https://www.bing.com/webmasters/help/getting-started-checklist-66a806de)。
