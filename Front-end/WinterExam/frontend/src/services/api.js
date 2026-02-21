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

  async sendCaptcha(phone) {
    return this.get('/captcha/sent', { phone })
  }

  async verifyCaptcha(phone, captcha) {
    return this.get('/captcha/verify', { phone, captcha })
  }

  async loginWithPhone(phone, captcha) {
    return this.get('/login/cellphone', { phone, captcha })
  }

  async loginWithPassword(phone, password) {
    return this.get('/login/cellphone', { phone, password })
  }

  async getLoginStatus() {
    return this.get('/login/status')
  }

  async logout() {
    return this.get('/logout')
  }

  async getBanner(type = 0) {
    return this.get('/banner', { type })
  }

  async getPersonalized(limit = 30) {
    return this.get('/personalized', { limit })
  }

  async getPlaylistCategories() {
    return this.get('/playlist/catlist')
  }

  async getPlaylistsByCategory(cat = '全部', limit = 50, offset = 0) {
    return this.get('/top/playlist', { cat, limit, offset })
  }

  async getPlaylistDetail(id) {
    return this.get('/playlist/detail', { id })
  }

  async getSongDetail(ids) {
    return this.get('/song/detail', { ids: Array.isArray(ids) ? ids.join(',') : ids })
  }

  async getSongUrl(id) {
    return this.get('/song/url', { id })
  }

  async getLyric(id) {
    return this.get('/lyric', { id })
  }

  async search(keywords, type = 1, limit = 30, offset = 0) {
    return this.get('/search', { keywords, type, limit, offset })
  }

  async getUserPlaylist(uid) {
    return this.get('/user/playlist', { uid })
  }

  async getLikelist(uid) {
    return this.get('/likelist', { uid })
  }

  async getQRKey() {
    return this.get('/login/qr/key')
  }

  async getQRCode(key, qrimg = true) {
    return this.get('/login/qr/create', { key, qrimg })
  }

  async checkQRCode(key) {
    return this.get('/login/qr/check', { key })
  }
}

export default new ApiService()
