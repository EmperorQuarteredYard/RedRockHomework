# 迭代1：基础框架搭建 - API接口文档

## 概述

本文档记录了迭代1（基础框架搭建）阶段需要使用的API接口，包括接口调用方法、参数说明和对应的api-enhanced接口信息。

## API基础配置

- **基础URL**: `http://localhost:3000` (api-enhanced默认端口)
- **请求方式**: GET/POST
- **数据格式**: JSON

---

## 1. 用户登录相关接口

### 1.1 发送验证码

**功能说明**: 向指定手机号发送登录验证码

**接口调用**:
```javascript
// 方法: GET 或 POST
// 路由: /captcha/sent
// 示例请求
fetch('http://localhost:3000/captcha/sent?phone=13800138000')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| phone | String | 是 | 手机号码 |
| ctcode | String | 否 | 国家代码，默认为86 |

**返回示例**:
```json
{
  "code": 200,
  "message": "发送成功"
}
```

**api-enhanced接口信息**:
- 接口名称: `captcha_sent`
- 原代码出处: `api-enhanced/module/captcha_sent.js`
- 内部API: `/api/sms/captcha/sent`

---

### 1.2 校验验证码

**功能说明**: 校验手机验证码是否正确

**接口调用**:
```javascript
// 方法: GET 或 POST
// 路由: /captcha/verify
// 示例请求
fetch('http://localhost:3000/captcha/verify?phone=13800138000&captcha=123456')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| phone | String | 是 | 手机号码 |
| captcha | String | 是 | 验证码 |
| ctcode | String | 否 | 国家代码，默认为86 |

**返回示例**:
```json
{
  "code": 200,
  "message": "验证成功"
}
```

**api-enhanced接口信息**:
- 接口名称: `captcha_verify`
- 原代码出处: `api-enhanced/module/captcha_verify.js`
- 内部API: `/api/sms/captcha/verify`

---

### 1.3 手机号登录

**功能说明**: 使用手机号和验证码或密码登录

**接口调用**:
```javascript
// 方法: GET 或 POST
// 路由: /login/cellphone
// 验证码登录示例
fetch('http://localhost:3000/login/cellphone?phone=13800138000&captcha=123456')
  .then(res => res.json())
  .then(data => console.log(data))

// 密码登录示例
fetch('http://localhost:3000/login/cellphone?phone=13800138000&password=yourpassword')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| phone | String | 是 | 手机号码 |
| password | String | 否 | 密码（与captcha二选一） |
| captcha | String | 否 | 验证码（与password二选一） |
| countrycode | String | 否 | 国家代码，默认为86 |

**返回示例**:
```json
{
  "code": 200,
  "account": {
    "id": 123456789,
    "userName": "用户名",
    "type": 1,
    "status": 0,
    "whitelistAuthority": 0
  },
  "profile": {
    "userId": 123456789,
    "nickname": "昵称",
    "avatarUrl": "http://...",
    "backgroundUrl": "http://...",
    "gender": 0,
    "city": 0,
    "province": 0,
    "country": "中国"
  },
  "cookie": "..."
}
```

**api-enhanced接口信息**:
- 接口名称: `login_cellphone`
- 原代码出处: `api-enhanced/module/login_cellphone.js`
- 内部API: `/api/w/login/cellphone`

---

### 1.4 获取登录状态

**功能说明**: 获取当前用户的登录状态和用户信息

**接口调用**:
```javascript
// 方法: GET
// 路由: /login/status
// 示例请求
fetch('http://localhost:3000/login/status', {
  credentials: 'include' // 携带cookie
})
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**: 无需参数，需要携带登录后的cookie

**返回示例**:
```json
{
  "data": {
    "code": 200,
    "account": {
      "id": 123456789,
      "userName": "用户名"
    },
    "profile": {
      "userId": 123456789,
      "nickname": "昵称",
      "avatarUrl": "http://..."
    }
  }
}
```

**api-enhanced接口信息**:
- 接口名称: `login_status`
- 原代码出处: `api-enhanced/module/login_status.js`
- 内部API: `/api/w/nuser/account/get`

---

### 1.5 退出登录

**功能说明**: 退出当前登录状态

**接口调用**:
```javascript
// 方法: GET
// 路由: /logout
// 示例请求
fetch('http://localhost:3000/logout', {
  credentials: 'include'
})
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**: 无需参数，需要携带登录后的cookie

**返回示例**:
```json
{
  "code": 200,
  "msg": "退出成功"
}
```

**api-enhanced接口信息**:
- 接口名称: `logout`
- 原代码出处: `api-enhanced/module/logout.js`

---

## 2. 首页相关接口

### 2.1 获取轮播图

**功能说明**: 获取首页轮播图数据

**接口调用**:
```javascript
// 方法: GET
// 路由: /banner
// 示例请求
fetch('http://localhost:3000/banner?type=0')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| type | Number | 否 | 设备类型：0-pc, 1-android, 2-iphone, 3-ipad，默认为0 |

**返回示例**:
```json
{
  "banners": [
    {
      "pic": "http://...",
      "targetId": 123456,
      "adid": null,
      "targetType": 10,
      "titleColor": "red",
      "typeTitle": "独家",
      "url": "http://...",
      "encodeId": "123456",
      "bannerId": "123456",
      "video": null
    }
  ]
}
```

**api-enhanced接口信息**:
- 接口名称: `banner`
- 原代码出处: `api-enhanced/module/banner.js`
- 内部API: `/api/v2/banner/get`

---

### 2.2 获取推荐歌单

**功能说明**: 获取首页推荐歌单列表

**接口调用**:
```javascript
// 方法: GET
// 路由: /personalized
// 示例请求
fetch('http://localhost:3000/personalized?limit=10')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| limit | Number | 否 | 返回数量，默认为30 |

**返回示例**:
```json
{
  "code": 200,
  "result": [
    {
      "id": 123456,
      "name": "歌单名称",
      "picUrl": "http://...",
      "playCount": 1000000,
      "copywriter": "推荐语",
      "type": 0
    }
  ]
}
```

**api-enhanced接口信息**:
- 接口名称: `personalized`
- 原代码出处: `api-enhanced/module/personalized.js`
- 内部API: `/api/personalized/playlist`

---

## 3. 歌单相关接口

### 3.1 获取歌单分类

**功能说明**: 获取所有歌单分类信息

**接口调用**:
```javascript
// 方法: GET
// 路由: /playlist/catlist
// 示例请求
fetch('http://localhost:3000/playlist/catlist')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**: 无

**返回示例**:
```json
{
  "code": 200,
  "all": {
    "name": "全部",
    "category": 0
  },
  "sub": [
    {
      "name": "华语",
      "category": 0,
      "hot": true
    }
  ],
  "categories": {
    "0": "语种",
    "1": "风格",
    "2": "场景",
    "3": "情感",
    "4": "主题"
  }
}
```

**api-enhanced接口信息**:
- 接口名称: `playlist_catlist`
- 原代码出处: `api-enhanced/module/playlist_catlist.js`
- 内部API: `/api/playlist/catalogue`

---

### 3.2 获取分类歌单

**功能说明**: 根据分类获取歌单列表

**接口调用**:
```javascript
// 方法: GET
// 路由: /top/playlist
// 示例请求
fetch('http://localhost:3000/top/playlist?cat=华语&limit=10&offset=0')
  .then(res => res.json())
  .then(data => console.log(data))
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| cat | String | 否 | 分类名称，默认为"全部" |
| order | String | 否 | 排序方式：hot-热门, new-最新，默认为hot |
| limit | Number | 否 | 返回数量，默认为50 |
| offset | Number | 否 | 偏移量，用于分页，默认为0 |

**返回示例**:
```json
{
  "code": 200,
  "playlists": [
    {
      "id": 123456,
      "name": "歌单名称",
      "coverImgUrl": "http://...",
      "creator": {
        "nickname": "创建者昵称",
        "avatarUrl": "http://..."
      },
      "playCount": 1000000,
      "trackCount": 100
    }
  ],
  "total": 1000,
  "more": true
}
```

**api-enhanced接口信息**:
- 接口名称: `top_playlist`
- 原代码出处: `api-enhanced/module/top_playlist.js`
- 内部API: `/api/playlist/list`

---

## 4. 接口调用封装示例

### 4.1 API服务封装

```javascript
// src/services/api.js
const BASE_URL = 'http://localhost:3000'

class ApiService {
  constructor() {
    this.baseURL = BASE_URL
  }

  async request(url, options = {}) {
    const defaultOptions = {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      }
    }
    
    const response = await fetch(`${this.baseURL}${url}`, {
      ...defaultOptions,
      ...options
    })
    
    return response.json()
  }

  async get(url, params = {}) {
    const queryString = new URLSearchParams(params).toString()
    const fullUrl = queryString ? `${url}?${queryString}` : url
    return this.request(fullUrl)
  }

  async post(url, data = {}) {
    return this.request(url, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  // 用户相关接口
  async sendCaptcha(phone) {
    return this.get('/captcha/sent', { phone })
  }

  async verifyCaptcha(phone, captcha) {
    return this.get('/captcha/verify', { phone, captcha })
  }

  async loginWithPhone(phone, captcha) {
    return this.get('/login/cellphone', { phone, captcha })
  }

  async getLoginStatus() {
    return this.get('/login/status')
  }

  async logout() {
    return this.get('/logout')
  }

  // 首页相关接口
  async getBanner(type = 0) {
    return this.get('/banner', { type })
  }

  async getPersonalized(limit = 30) {
    return this.get('/personalized', { limit })
  }

  // 歌单相关接口
  async getPlaylistCategories() {
    return this.get('/playlist/catlist')
  }

  async getPlaylistsByCategory(cat = '全部', limit = 50, offset = 0) {
    return this.get('/top/playlist', { cat, limit, offset })
  }
}

export default new ApiService()
```

### 4.2 使用示例

```javascript
// 发送验证码
const sendCode = async (phone) => {
  const result = await api.sendCaptcha(phone)
  if (result.code === 200) {
    console.log('验证码发送成功')
  }
}

// 验证码登录
const login = async (phone, captcha) => {
  const result = await api.loginWithPhone(phone, captcha)
  if (result.code === 200) {
    console.log('登录成功', result.profile)
  }
}

// 获取用户信息
const getUserInfo = async () => {
  const result = await api.getLoginStatus()
  if (result.data.code === 200 && result.data.profile) {
    console.log('用户信息', result.data.profile)
    return result.data.profile
  }
  return null
}

// 获取轮播图
const getBanners = async () => {
  const result = await api.getBanner()
  if (result.banners) {
    return result.banners
  }
  return []
}

// 获取推荐歌单
const getRecommendPlaylists = async (limit = 10) => {
  const result = await api.getPersonalized(limit)
  if (result.code === 200) {
    return result.result
  }
  return []
}
```

---

## 5. 注意事项

1. **跨域问题**: api-enhanced已配置CORS，允许跨域请求
2. **Cookie管理**: 登录后需要保存cookie，后续请求需要携带cookie
3. **错误处理**: 所有接口都有可能返回错误，需要进行错误处理
4. **请求频率**: 避免频繁请求，部分接口有频率限制
5. **数据缓存**: 建议对不常变化的数据进行本地缓存

---

## 6. 接口汇总表

| 功能 | 接口名称 | 路由 | api-enhanced文件 |
|------|----------|------|------------------|
| 发送验证码 | captcha_sent | /captcha/sent | module/captcha_sent.js |
| 校验验证码 | captcha_verify | /captcha/verify | module/captcha_verify.js |
| 手机号登录 | login_cellphone | /login/cellphone | module/login_cellphone.js |
| 获取登录状态 | login_status | /login/status | module/login_status.js |
| 退出登录 | logout | /logout | module/logout.js |
| 获取轮播图 | banner | /banner | module/banner.js |
| 获取推荐歌单 | personalized | /personalized | module/personalized.js |
| 获取歌单分类 | playlist_catlist | /playlist/catlist | module/playlist_catlist.js |
| 获取分类歌单 | top_playlist | /top/playlist | module/top_playlist.js |
