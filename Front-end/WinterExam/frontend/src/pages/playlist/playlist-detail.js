import api from '../../services/api.js'
import store from '../../store/store.js'

class PlaylistDetail {
  constructor() {
    this.playlist = null
    this.songs = []
  }

  render() {
    return `
      <div class="playlist-detail" id="playlistDetail">
        <div class="detail-loading">加载中...</div>
      </div>
    `
  }

  renderContent() {
    if (!this.playlist) return ''
    
    const playlist = this.playlist
    const songs = this.songs
    
    return `
      <div class="playlist-header">
        <div class="playlist-cover-large">
          <img src="${playlist.coverImgUrl}" alt="${playlist.name}">
        </div>
        <div class="playlist-meta">
          <h1 class="playlist-title">${playlist.name}</h1>
          <div class="playlist-creator">
            <img src="${playlist.creator?.avatarUrl || ''}" alt="" class="creator-avatar">
            <span class="creator-name">${playlist.creator?.nickname || '未知'}</span>
          </div>
          <div class="playlist-stats">
            <span class="stat-item">
              <span class="stat-label">歌曲数</span>
              <span class="stat-value">${playlist.trackCount || songs.length}</span>
            </span>
            <span class="stat-item">
              <span class="stat-label">播放量</span>
              <span class="stat-value">${this.formatCount(playlist.playCount)}</span>
            </span>
          </div>
          <div class="playlist-actions">
            <button class="action-btn primary" onclick="window.playlistDetail.playAll()">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M8 5v14l11-7z"/>
              </svg>
              播放全部
            </button>
          </div>
          ${playlist.description ? `
            <div class="playlist-desc">
              <p>${playlist.description}</p>
            </div>
          ` : ''}
        </div>
      </div>
      <div class="song-list">
        <div class="song-list-header">
          <span class="col-index">#</span>
          <span class="col-title">歌曲</span>
          <span class="col-artist">歌手</span>
          <span class="col-album">专辑</span>
          <span class="col-duration">时长</span>
        </div>
        <div class="song-list-body">
          ${songs.map((song, index) => this.renderSong(song, index)).join('')}
        </div>
      </div>
    `
  }

  renderSong(song, index) {
    const artists = song.ar || song.artists || []
    const album = song.al || song.album || {}
    const duration = song.dt || song.duration || 0
    
    return `
      <div class="song-item" data-id="${song.id}" onclick="window.playlistDetail.playSong(${index})">
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

  formatCount(count) {
    if (!count) return '0'
    if (count >= 100000000) {
      return (count / 100000000).toFixed(1) + '亿'
    } else if (count >= 10000) {
      return (count / 10000).toFixed(1) + '万'
    }
    return count.toString()
  }

  formatDuration(ms) {
    if (!ms) return '00:00'
    const seconds = Math.floor(ms / 1000)
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  async loadData(id) {
    const container = document.getElementById('playlistDetail')
    if (!container) return
    
    try {
      const result = await api.getPlaylistDetail(id)
      
      if (result.code === 200 && result.playlist) {
        this.playlist = result.playlist
        
        const trackIds = result.playlist.trackIds?.map(t => t.id) || []
        if (trackIds.length > 0) {
          const idsToFetch = trackIds.slice(0, 100)
          const songsResult = await api.getSongDetail(idsToFetch)
          
          if (songsResult.code === 200 && songsResult.songs) {
            this.songs = songsResult.songs
          }
        }
        
        container.innerHTML = this.renderContent()
      } else {
        container.innerHTML = '<div class="error-message">加载失败</div>'
      }
    } catch (error) {
      console.error('Failed to load playlist:', error)
      container.innerHTML = '<div class="error-message">加载失败</div>'
    }
  }

  playAll() {
    if (this.songs.length > 0) {
      store.setPlaylist(this.songs, 0)
    }
  }

  playSong(index) {
    if (index >= 0 && index < this.songs.length) {
      store.setPlaylist(this.songs, index)
    }
  }
}

const playlistDetail = new PlaylistDetail()
window.playlistDetail = playlistDetail

export default playlistDetail
