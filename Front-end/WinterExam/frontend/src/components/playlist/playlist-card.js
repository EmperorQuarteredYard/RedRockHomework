class PlaylistCard {
  constructor() {
    this.playlists = []
  }

  render(playlist, index) {
    const playCount = this.formatPlayCount(playlist.playCount)
    
    return `
      <div class="playlist-card" data-id="${playlist.id}" onclick="window.playlistCard.openPlaylist(${playlist.id})">
        <div class="playlist-cover">
          <img src="${playlist.picUrl || playlist.coverImgUrl}" alt="${playlist.name}" loading="lazy">
          <div class="play-count">
            <span>▶</span>
            <span>${playCount}</span>
          </div>
          <div class="play-overlay">
            <button class="play-btn">▶</button>
          </div>
        </div>
        <div class="playlist-info">
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
