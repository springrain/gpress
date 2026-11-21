# FixIt gpress Theme

一个为 gpress 新增的 FixIt-like 博客主题, 参考 [hugo-fixit/FixIt](https://github.com/hugo-fixit/FixIt) 的阅读取向与技术博客体验, 按 gpress 现有模板系统重新实现.

## 特性

- 响应式顶栏、搜索框、移动端菜单
- 亮色/暗色模式切换
- 首页文章卡片、分类页、标签页、404 页面
- 文章页封面图、目录侧栏、文章信息侧栏
- 代码块增强：语言标题、行号、换行、复制、下载、全屏

## 使用

把本目录放在：

```text
gpressdatadir/template/theme/gpress-FixIt
```

然后在 gpress 后台选择 `gpress-FixIt` 作为当前主题.

## 文件结构

模板页均提供 `.html`（页面渲染）与同名 `.md`（Markdown 输出, 用于 Feed/`.md` 路由）两份, 按 gpress 模板系统成对使用.

- `bodystart.html`: head、导航、搜索、移动菜单、main 开始
- `bodyend.html`: footer、页脚水波动画、返回顶部、脚本引用
- `index.html` / `index.md`: 首页和搜索结果页
- `content.html` / `content.md`: 文章详情页
- `category.html` / `category.md`: 分类文章列表
- `category-single.html` / `category-single.md`: 单篇独立页（单页/专题）
- `tag.html` / `tag.md`: 标签文章列表
- `error.html` / `error.md`: 404 页面
- `FIXIT-NOTICE.txt`: 代码块样式来源的 FixIt MIT 许可声明
- `css/theme.css`: 主题样式
- `js/theme.js`: 交互与代码块增强（暗色切换、目录、代码块、Mermaid、灯箱）
- `js/footer-water.js`: 页脚 Canvas 水波动画
- `js/jquery.min.js`、`js/jquery.fancybox.min.js`/`.css`: 灯箱图库依赖
- `js/mermaid.min.js`: Mermaid 图表渲染