class PlaylistCard {
  constructor() {
    this.playlists = []
  }

  render(playlist, index) {
    const playCount = this.formatPlayCount(playlist.playCount)
    const picUrl = playlist.picUrl || playlist.coverImgUrl || ''
    
    return `
      <div class="playlist-card" data-id="${playlist.id}">
        <div class="playlist-cover">
          <img src="${picUrl}" alt="${playlist.name}" loading="lazy" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><rect fill=%22%23ddd%22 width=%22100%22 height=%22100%22/><text x=%2250%22 y=%2250%22 text-anchor=%22middle%22 dy=%22.3em%22 fill=%22%23999%22>暂无封面</text></svg>'">
          <div class="play-count">
            <span>▶</span>
            <span>${playCount}</span>
          </div>
          <div class="play-overlay" onclick="window.playlistCard.openPlaylist(${playlist.id})">
            <button class="play-btn">▶</button>
          </div>
        </div>
        <div class="playlist-info" onclick="window.playlistCard.openPlaylist(${playlist.id})">
          <h4 class="playlist-name">${playlist.name}</h4>
        </div>
      </div>
    `
  }

  renderList(playlists) {
    this.playlists = playlists
    return `
      <div class="playlist-grid">
        ${playlists.map((playlist, index) => this.render(playlist, index)).join('')}
      </div>
    `
  }

  formatPlayCount(count) {
    if (count >= 100000000) {
      return (count / 100000000).toFixed(1) + '亿'
    } else if (count >= 10000) {
      return (count / 10000).toFixed(1) + '万'
    }
    return count.toString()
  }

  openPlaylist(id) {
    const event = new CustomEvent('navigate', { 
      detail: { page: 'playlist-detail', id } 
    })
    window.dispatchEvent(event)
  }
}

const playlistCard = new PlaylistCard()
window.playlistCard = playlistCard

export default playlistCard
