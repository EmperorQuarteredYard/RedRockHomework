import api from '../../services/api.js'
import store from '../../store/store.js'

class SearchPage {
  constructor() {
    this.keyword = ''
    this.results = []
    this.type = 1
    this.loading = false
  }

  render() {
    return `
      <div class="search-page" id="searchPage">
        <div class="search-header">
          <h2>搜索结果</h2>
          <span class="search-keyword" id="searchKeyword"></span>
        </div>
        <div class="search-content" id="searchContent">
          <div class="search-placeholder">输入关键词搜索歌曲</div>
        </div>
      </div>
    `
  }

  renderResults() {
    if (this.results.length === 0) {
      return `
        <div class="search-empty">
          <p>未找到相关结果</p>
        </div>
      `
    }

    return `
      <div class="song-list">
        <div class="song-list-header">
          <span class="col-index">#</span>
          <span class="col-title">歌曲</span>
          <span class="col-artist">歌手</span>
          <span class="col-album">专辑</span>
          <span class="col-duration">时长</span>
        </div>
        <div class="song-list-body">
          ${this.results.map((song, index) => this.renderSong(song, index)).join('')}
        </div>
      </div>
    `
  }

  renderSong(song, index) {
    const artists = song.ar || song.artists || []
    const album = song.al || song.album || {}
    const duration = song.dt || song.duration || 0
    
    return `
      <div class="song-item" data-id="${song.id}" onclick="window.searchPage.playSong(${index})">
        <span class="col-index">${index + 1}</span>
        <span class="col-title">
          <span class="song-name">${song.name}</span>
        </span>
        <span class="col-artist">${artists.map(a => a.name).join('/')}</span>
        <span class="col-album">${album.name || '-'}</span>
        <span class="col-duration">${this.formatDuration(duration)}</span>
      </div>
    `
  }

  formatDuration(ms) {
    if (!ms) return '00:00'
    const seconds = Math.floor(ms / 1000)
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  async search(keyword) {
    if (!keyword || keyword.trim() === '') {
      return
    }

    this.keyword = keyword.trim()
    this.loading = true
    
    const keywordEl = document.getElementById('searchKeyword')
    const contentEl = document.getElementById('searchContent')
    
    if (keywordEl) {
      keywordEl.textContent = `"${this.keyword}"`
    }
    
    if (contentEl) {
      contentEl.innerHTML = '<div class="search-loading">搜索中...</div>'
    }

    try {
      const result = await api.search(this.keyword, this.type, 50)
      
      if (result.code === 200 && result.result) {
        if (this.type === 1) {
          this.results = result.result.songs || []
        }
        
        if (contentEl) {
          contentEl.innerHTML = this.renderResults()
        }
      } else {
        if (contentEl) {
          contentEl.innerHTML = '<div class="search-empty"><p>搜索失败，请重试</p></div>'
        }
      }
    } catch (error) {
      console.error('Search failed:', error)
      if (contentEl) {
        contentEl.innerHTML = '<div class="search-empty"><p>搜索失败，请重试</p></div>'
      }
    } finally {
      this.loading = false
    }
  }

  playSong(index) {
    if (index >= 0 && index < this.results.length) {
      store.setPlaylist(this.results, index)
    }
  }

  playAll() {
    if (this.results.length > 0) {
      store.setPlaylist(this.results, 0)
    }
  }
}

const searchPage = new SearchPage()
window.searchPage = searchPage

export default searchPage
