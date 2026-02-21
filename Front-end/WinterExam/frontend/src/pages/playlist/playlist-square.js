import api from '../../services/api.js'
import playlistCard from '../../components/playlist/playlist-card.js'

class PlaylistSquare {
  constructor() {
    this.categories = []
    this.currentCategory = '全部'
    this.playlists = []
    this.offset = 0
    this.limit = 50
    this.hasMore = true
    this.isLoading = false
  }

  render() {
    return `
      <div class="playlist-square-page">
        <div class="category-section">
          <div class="category-header">
            <h3>全部歌单</h3>
          </div>
          <div class="category-list" id="categoryList">
            <div class="loading">加载分类中...</div>
          </div>
        </div>
        <div class="playlist-section">
          <div class="playlist-list" id="playlistList">
            <div class="loading">加载歌单中...</div>
          </div>
          <div class="load-more" id="loadMore">
            <button class="load-more-btn" onclick="window.playlistSquare.loadMore()">加载更多</button>
          </div>
        </div>
      </div>
    `
  }

  async loadData() {
    await this.loadCategories()
    await this.loadPlaylists()
  }

  async loadCategories() {
    try {
      const result = await api.getPlaylistCategories()
      if (result.code === 200) {
        this.categories = result.sub || []
        this.renderCategories()
      }
    } catch (error) {
      console.error('加载分类失败:', error)
    }
  }

  renderCategories() {
    const container = document.getElementById('categoryList')
    if (!container) return

    const allCategory = { name: '全部', category: 0 }
    const hotCategories = this.categories.filter(cat => cat.hot)
    
    let html = `
      <span class="category-item ${this.currentCategory === '全部' ? 'active' : ''}" 
            onclick="window.playlistSquare.selectCategory('全部')">全部</span>
    `
    
    hotCategories.slice(0, 15).forEach(cat => {
      html += `
        <span class="category-item ${this.currentCategory === cat.name ? 'active' : ''}" 
              onclick="window.playlistSquare.selectCategory('${cat.name}')">${cat.name}</span>
      `
    })
    
    container.innerHTML = html
  }

  async loadPlaylists(reset = false) {
    if (this.isLoading) return
    
    if (reset) {
      this.offset = 0
      this.playlists = []
      this.hasMore = true
    }

    this.isLoading = true
    
    try {
      const result = await api.getPlaylistsByCategory(
        this.currentCategory, 
        this.limit, 
        this.offset
      )
      
      if (result.code === 200 && result.playlists) {
        this.playlists = reset ? result.playlists : [...this.playlists, ...result.playlists]
        this.hasMore = result.more
        this.offset += this.limit
        this.renderPlaylists()
      }
    } catch (error) {
      console.error('加载歌单失败:', error)
      const container = document.getElementById('playlistList')
      if (container) {
        container.innerHTML = '<div class="error">加载失败，请稍后重试</div>'
      }
    } finally {
      this.isLoading = false
    }
  }

  renderPlaylists() {
    const container = document.getElementById('playlistList')
    if (container && this.playlists.length > 0) {
      container.innerHTML = playlistCard.renderList(this.playlists)
    }
    
    const loadMoreBtn = document.querySelector('.load-more-btn')
    if (loadMoreBtn) {
      loadMoreBtn.style.display = this.hasMore ? 'block' : 'none'
    }
  }

  selectCategory(category) {
    this.currentCategory = category
    this.renderCategories()
    this.loadPlaylists(true)
  }

  async loadMore() {
    if (this.hasMore && !this.isLoading) {
      await this.loadPlaylists()
    }
  }
}

const playlistSquare = new PlaylistSquare()
window.playlistSquare = playlistSquare

export default playlistSquare
