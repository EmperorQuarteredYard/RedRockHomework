import api from '../../services/api.js'
import loginModal from '../login/login.js'

class Header {
  constructor() {
    this.userInfo = null
    this.isLoggedIn = false
  }

  render() {
    return `
      <header class="header">
        <div class="header-left">
          <div class="logo">
            <img src="https://s1.music.126.net/style/favicon.ico" alt="网易云音乐">
            <span>网易云音乐</span>
          </div>
        </div>
        <div class="header-center">
          <div class="search-box">
            <input type="text" placeholder="搜索音乐、歌手、歌词" id="searchInput">
            <button class="search-btn" onclick="window.header.search()">搜索</button>
          </div>
        </div>
        <div class="header-right" id="userInfoContainer">
          ${this.renderUserInfo()}
        </div>
      </header>
    `
  }

  renderUserInfo() {
    if (this.isLoggedIn && this.userInfo) {
      return `
        <div class="user-info">
          <img class="user-avatar" src="${this.userInfo.avatarUrl || 'https://s1.music.126.net/style/favicon.ico'}" alt="头像">
          <span class="user-name">${this.userInfo.nickname || '用户'}</span>
          <button class="logout-btn" onclick="window.header.logout()">退出</button>
        </div>
      `
    }
    return `
      <button class="login-btn" onclick="window.header.showLogin()">登录</button>
    `
  }

  async checkLoginStatus() {
    try {
      const result = await api.getLoginStatus()
      if (result.data && result.data.code === 200 && result.data.profile) {
        this.userInfo = result.data.profile
        this.isLoggedIn = true
        this.updateUserInfo()
        return true
      }
    } catch (error) {
      console.error('检查登录状态失败:', error)
    }
    this.isLoggedIn = false
    this.userInfo = null
    this.updateUserInfo()
    return false
  }

  updateUserInfo() {
    const container = document.getElementById('userInfoContainer')
    if (container) {
      container.innerHTML = this.renderUserInfo()
    }
  }

  showLogin() {
    loginModal.show()
  }

  async logout() {
    try {
      await api.logout()
      this.isLoggedIn = false
      this.userInfo = null
      this.updateUserInfo()
      alert('已退出登录')
    } catch (error) {
      console.error('退出登录失败:', error)
    }
  }

  search() {
    const keyword = document.getElementById('searchInput')?.value?.trim()
    if (keyword) {
      const event = new CustomEvent('navigate', { 
        detail: { page: 'search', keyword } 
      })
      window.dispatchEvent(event)
    }
  }

  init() {
    loginModal.onLoginSuccess = (profile) => {
      this.userInfo = profile
      this.isLoggedIn = true
      this.updateUserInfo()
    }
    
    this.checkLoginStatus()
    
    document.getElementById('searchInput')?.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') {
        this.search()
      }
    })
  }
}

const header = new Header()
window.header = header

export default header
