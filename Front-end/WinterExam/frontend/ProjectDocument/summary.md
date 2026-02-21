# 音乐播放器项目开发总结

## 项目概述

本项目是一个仿网易云音乐风格的音乐播放器Web应用，采用快速原型迭代开发模式，共完成4个迭代周期。

## 技术栈

- **前端框架**: 原生JavaScript (ES6+)
- **构建工具**: Vite
- **样式**: 原生CSS
- **API**: api-enhanced (网易云音乐API)
- **模块化**: ES Modules

## 迭代一：基础架构搭建

### 目标
- 搭建项目基础结构
- 实现用户登录功能
- 完成基本页面布局
- API连接配置

### 完成的工作

#### 1. 项目结构创建
```
frontend/
├── src/
│   ├── components/        # 组件目录
│   │   ├── carousel/      # 轮播图组件
│   │   ├── header/        # 头部组件
│   │   ├── login/         # 登录组件
│   │   ├── playlist/      # 歌单卡片组件
│   │   ├── player/        # 播放器组件
│   │   └── sidebar/       # 侧边栏组件
│   ├── pages/             # 页面目录
│   │   ├── home/          # 首页
│   │   ├── playlist/      # 歌单相关页面
│   │   └── search/        # 搜索页面
│   ├── services/          # API服务
│   ├── store/             # 状态管理
│   └── styles/            # 样式文件
└── ProjectDocument/       # 项目文档
```

#### 2. API服务封装
创建了 `api.js`，封装了所有API请求：
- 登录相关：验证码发送/验证、密码登录、二维码登录
- 用户相关：登录状态、登出、用户歌单
- 内容相关：轮播图、推荐歌单、歌单详情、歌曲URL、歌词
- 搜索功能

#### 3. 登录功能实现

**遇到的问题**：
1. **验证码未发送到手机** - 网易云音乐API的验证码服务有严格限制
2. **二维码显示问题** - 初始实现时二维码无法正确渲染

**解决方案**：
- 实现二维码登录作为主要登录方式
- 添加QR Key获取、二维码生成、状态检查的完整流程
- 处理API返回数据结构的兼容性（`data.unikey` vs `body.data.unikey`）

```javascript
// 二维码登录流程
async initQRCodeLogin() {
  // 1. 获取QR Key
  const keyResult = await api.getQRKey()
  // 2. 生成二维码
  const qrResult = await api.getQRCode(this.qrKey, true)
  // 3. 轮询检查状态
  this.startQRCodeCheck()
}
```

#### 4. 核心组件开发
- **Header**: Logo、搜索框、用户信息显示
- **Sidebar**: 导航菜单
- **Carousel**: 首页轮播图
- **PlaylistCard**: 歌单卡片展示

## 迭代二：核心功能原型

### 目标
- 实现播放器核心功能
- 实现首页轮播图和推荐歌单
- 实现歌单列表展示

### 完成的工作

#### 1. 播放器状态管理 (store.js)
创建了全局状态管理器，管理：
- 当前播放歌曲
- 播放列表
- 播放状态（播放/暂停）
- 播放进度
- 音量控制

```javascript
class MusicStore {
  constructor() {
    this.currentSong = null
    this.playlist = []
    this.currentIndex = -1
    this.isPlaying = false
    this.audio = new Audio()
  }
  
  // 播放控制方法
  async playSong(song, index)
  togglePlay()
  next()
  prev()
  seek(time)
  setVolume(volume)
}
```

#### 2. 播放器组件 (player.js)
底部播放器栏，包含：
- 歌曲信息显示（封面、歌名、歌手）
- 播放控制按钮（上一首、播放/暂停、下一首）
- 进度条（可拖拽）
- 音量控制

#### 3. 歌单详情页 (playlist-detail.js)
- 歌单信息展示（封面、名称、创建者、播放量）
- 歌曲列表展示
- 点击歌曲播放功能
- "播放全部"按钮

## 迭代三：功能完善

### 目标
- 实现搜索功能
- 完善播放器功能
- 实现歌曲列表播放控制

### 完成的工作

#### 1. 搜索页面 (search.js)
- 关键词搜索歌曲
- 搜索结果列表展示
- 点击搜索结果播放

#### 2. 路由系统完善
更新 `main.js` 支持多页面导航：
- 首页 (home)
- 歌单广场 (playlist-square)
- 歌单详情 (playlist-detail)
- 搜索页 (search)

使用自定义事件实现页面切换：
```javascript
window.addEventListener('navigate', (e) => {
  const { page, ...params } = e.detail
  this.navigateTo(page, params)
})
```

## 迭代四：优化与测试

### 目标
- 实现播放详情页和歌词显示
- 优化用户体验
- 进行功能测试和bug修复

### 完成的工作

#### 1. 播放详情页 (player-detail.js)
- 全屏播放界面
- 模糊背景效果（使用专辑封面）
- 旋转唱片动画
- 歌词同步显示

```javascript
// 歌词解析
parseLyric(lrc) {
  const lines = lrc.split('\n')
  const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2,3})\]/
  // 解析时间戳和歌词文本
}

// 歌词同步
updateLyricPosition(currentTime) {
  // 根据当前播放时间高亮对应歌词
}
```

#### 2. 交互优化
- 点击播放器封面打开详情页
- 歌词自动滚动
- 播放状态同步

#### 3. Bug修复

**问题1：轮播图点击跳转到桌面应用**
- 原因：轮播图使用 `<a>` 标签链接到外部URL
- 解决：移除 `<a>` 标签，使用JavaScript处理点击事件

**问题2：图片无法显示**
- 原因：API返回数据结构不一致
- 解决：添加多种数据格式兼容处理，添加图片加载错误处理

## 项目文件清单

### 核心文件

| 文件 | 功能 |
|------|------|
| `main.js` | 应用入口，路由管理 |
| `api.js` | API服务封装 |
| `store.js` | 全局状态管理 |

### 组件文件

| 文件 | 功能 |
|------|------|
| `headbar.js` | 顶部导航栏 |
| `sidebar.js` | 侧边栏导航 |
| `carousel.js` | 轮播图组件 |
| `playlist-card.js` | 歌单卡片 |
| `player.js` | 底部播放器 |
| `player-detail.js` | 播放详情页 |
| `login.js` | 登录弹窗 |

### 页面文件

| 文件 | 功能 |
|------|------|
| `home.js` | 首页 |
| `playlist-square.js` | 歌单广场 |
| `playlist-detail.js` | 歌单详情 |
| `search.js` | 搜索页面 |

## 遇到的主要问题及解决方案

### 1. 跨域问题
- **问题**：前端直接请求API会遇到CORS错误
- **解决**：在 `vite.config.js` 配置代理
```javascript
proxy: {
  '/api': {
    target: 'http://localhost:3000',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '')
  }
}
```

### 2. 登录验证码问题
- **问题**：验证码未发送到手机
- **解决**：实现二维码登录作为替代方案

### 3. 二维码显示问题
- **问题**：二维码图片无法正确渲染
- **解决**：
  - 处理API返回数据结构的兼容性
  - 添加加载状态和错误处理
  - 使用第三方服务生成二维码作为备用方案

### 4. 模块化设计
- **问题**：组件间通信复杂
- **解决**：
  - 使用发布-订阅模式实现状态管理
  - 使用自定义事件实现页面导航

## 项目亮点

1. **纯原生实现**：不依赖任何前端框架，使用原生JavaScript实现所有功能
2. **模块化设计**：组件、页面、服务、状态分离，代码结构清晰
3. **响应式状态管理**：使用发布-订阅模式，UI自动响应状态变化
4. **歌词同步**：实现歌词解析和时间同步高亮
5. **多种登录方式**：支持二维码、验证码、密码三种登录方式

## 后续优化建议

1. **性能优化**
   - 图片懒加载
   - 列表虚拟滚动
   - 代码分割

2. **功能完善**
   - 播放模式切换（顺序、随机、单曲循环）
   - 收藏功能
   - 评论功能

3. **用户体验**
   - 加载骨架屏
   - 错误边界处理
   - 离线缓存

## 总结

本项目通过4个迭代周期，完成了一个功能完整的音乐播放器Web应用。采用快速原型迭代方法，每个迭代都有明确的目标和交付物，逐步完善功能。在开发过程中遇到了登录验证码、二维码显示等问题，通过技术调研和方案调整最终解决。项目展示了原生JavaScript模块化开发的能力，以及状态管理、组件设计等前端核心概念的应用。
谁知道一看项目开发指导，是让我用web登录时的cookie数据(网易云不用JWT登录发吗？？)