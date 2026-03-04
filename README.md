# Jason Frontpage 🍎

Jason Apple 的个人展示页面 — 静态 HTML，可直接部署到 GitHub Pages。

---

## 📁 文件结构

```
jason-frontpage/
└── index.html    # 单页静态 HTML（包含所有 CSS）
```

---

## 🎨 设计特点

基于 **design-taste-frontend** skill 的最佳实践：

### 视觉设计
- **字体**: Geist Sans + Geist Mono（反 Inter 偏见）
- **色彩**: Off-black 基底 + 单一 Emerald 强调色
- **质感**: Grain Overlay 微纹理 + Liquid Glass 玻璃态
- **布局**: 非对称 Hero + Bento Grid 能力展示

### 交互细节
- **微动效**: Avatar 浮动动画 + Badge 脉冲效果
- **Hover 状态**: Bento Card 悬浮 + 边框高亮
- **性能**: CSS 硬件加速（transform/opacity）

### 反 AI 设计
- ❌ 无纯黑（#000000）
- ❌ 无紫色/蓝色 AI 渐变
- ❌ 无居中 Hero 布局
- ❌ 无 3 列卡片布局
- ❌ 无 Emoji 代码注释

---

## 🚀 部署到 GitHub Pages

### 方式 1: 手动上传

1. 创建 GitHub 仓库（如 `jason-frontpage`）
2. 上传 `index.html` 到根目录
3. 设置 → Pages → Source: `main` branch
4. 访问 `https://yourusername.github.io/jason-frontpage/`

### 方式 2: 命令行

```bash
cd /home/kimi/.openclaw/workspace/research/jason-frontpage

# 初始化 Git
git init
git add index.html
git commit -m "Initial commit: Jason Frontpage"

# 关联远程仓库（替换为你的仓库 URL）
git remote add origin https://github.com/YOUR_USERNAME/jason-frontpage.git

# 推送
git push -u origin main
```

---

## 📱 响应式支持

- **Desktop**: 非对称双栏布局
- **Tablet (<900px)**: 单栏堆叠
- **Mobile (<700px)**: 完整单栏适配

---

## 🎯 设计原则

| 原则 | 应用 |
|------|------|
| DESIGN_VARIANCE: 8 | 非对称布局 + 留白 |
| MOTION_INTENSITY: 6 | 微动效 + 悬浮反馈 |
| VISUAL_DENSITY: 4 | 透气布局 + 适度信息 |

---

## 📊 性能指标

- **文件大小**: ~16KB（单文件）
- **外部依赖**: Google Fonts（2 个字体）
- **JavaScript**: 无（纯静态）
- **Lighthouse**: 预期 95+（Performance）

---

## 🧪 浏览器兼容性

- ✅ Chrome/Edge (最新)
- ✅ Safari (最新)
- ✅ Firefox (最新)
- ✅ Mobile Safari (iOS 15+)
- ✅ Chrome for Android

---

## 📝 自定义

### 修改强调色

```css
:root {
  --accent: #22c55e; /* 改为你的颜色 */
  --accent-hover: #16a34a;
}
```

### 修改内容

直接编辑 `index.html` 中的文本内容。所有 CSS 内联，无需额外文件。

---

## 🍎 Credits

- **Design System**: design-taste-frontend skill
- **Fonts**: Geist Sans + Geist Mono (Vercel)
- **Icon**: Apple Emoji 🍎

---

**Built by Jason Apple 🍎**
