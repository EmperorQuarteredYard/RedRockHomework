import store from '../../store/store.js'
import api from '../../services/api.js'

class PlayerDetail {
  constructor() {
    this.isVisible = false
    this.lyric = []
    this.currentLyricIndex = -1
    this.unsubscribe = null
  }

  render() {
    return `
      <div class="player-detail" id="playerDetail">
        <div class="detail-background" id="detailBg"></div>
        <div class="detail-content">
          <button class="detail-close" onclick="window.playerDetail.hide()">×</button>
          <div class="detail-main">
            <div class="detail-left">
              <div class="disc-container">
                <div class="disc" id="disc">
                  <div class="disc-cover">
                    <img src="" alt="" id="detailCover">
                  </div>
                </div>
              </div>
            </div>
            <div class="detail-right">
              <div class="song-info">
                <h2 class="detail-song-name" id="detailSongName">未播放歌曲</h2>
                <div class="detail-artist-album">
                  <span class="detail-artist" id="detailArtist">-</span>
                  <span class="separator">-</span>
                  <span class="detail-album" id="detailAlbum">-</span>
                </div>
              </div>
              <div class="lyric-container" id="lyricContainer">
                <div class="lyric-content" id="lyricContent">
                  <p class="lyric-placeholder">暂无歌词</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    `
  }

  mount(container) {
    this.container = container
    this.container.innerHTML = this.render()
    this.unsubscribe = store.subscribe((state) => this.updateUI(state))
  }

  show() {
    this.isVisible = true
    const detail = document.getElementById('playerDetail')
    if (detail) {
      detail.classList.add('visible')
    }
    this.updateUI(store.getState())
  }

  hide() {
    this.isVisible = false
    const detail = document.getElementById('playerDetail')
    if (detail) {
      detail.classList.remove('visible')
    }
  }

  toggle() {
    if (this.isVisible) {
      this.hide()
    } else {
      this.show()
    }
  }

  async updateUI(state) {
    const { currentSong, currentTime, isPlaying } = state
    
    const detailBg = document.getElementById('detailBg')
    const detailCover = document.getElementById('detailCover')
    const disc = document.getElementById('disc')
    const detailSongName = document.getElementById('detailSongName')
    const detailArtist = document.getElementById('detailArtist')
    const detailAlbum = document.getElementById('detailAlbum')

    if (currentSong) {
      const picUrl = currentSong.al?.picUrl || currentSong.album?.picUrl || currentSong.picUrl || ''
      
      if (detailBg) {
        detailBg.style.backgroundImage = `url(${picUrl})`
      }
      if (detailCover) {
        detailCover.src = picUrl
      }
      if (detailSongName) {
        detailSongName.textContent = currentSong.name || '未知歌曲'
      }
      if (detailArtist) {
        const artists = currentSong.ar || currentSong.artists || []
        detailArtist.textContent = artists.map(a => a.name).join('/') || '-'
      }
      if (detailAlbum) {
        const album = currentSong.al || currentSong.album || {}
        detailAlbum.textContent = album.name || '-'
      }

      if (!this.currentSongId || this.currentSongId !== currentSong.id) {
        this.currentSongId = currentSong.id
        this.loadLyric(currentSong.id)
      }
    }

    if (disc) {
      if (isPlaying) {
        disc.classList.add('playing')
      } else {
        disc.classList.remove('playing')
      }
    }

    this.updateLyricPosition(currentTime)
  }

  async loadLyric(songId) {
    const lyricContent = document.getElementById('lyricContent')
    
    try {
      const result = await api.getLyric(songId)
      
      if (result.code === 200 && result.lrc && result.lrc.lyric) {
        this.lyric = this.parseLyric(result.lrc.lyric)
        this.renderLyric()
      } else {
        this.lyric = []
        if (lyricContent) {
          lyricContent.innerHTML = '<p class="lyric-placeholder">暂无歌词</p>'
        }
      }
    } catch (error) {
      console.error('Failed to load lyric:', error)
      this.lyric = []
      if (lyricContent) {
        lyricContent.innerHTML = '<p class="lyric-placeholder">歌词加载失败</p>'
      }
    }
  }

  parseLyric(lrc) {
    const lines = lrc.split('\n')
    const result = []
    const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2,3})\]/

    for (const line of lines) {
      const match = line.match(timeRegex)
      if (match) {
        const minutes = parseInt(match[1])
        const seconds = parseInt(match[2])
        const milliseconds = parseInt(match[3].padEnd(3, '0'))
        const time = minutes * 60 + seconds + milliseconds / 1000
        const text = line.replace(timeRegex, '').trim()
        
        if (text) {
          result.push({ time, text })
        }
      }
    }

    return result
  }

  renderLyric() {
    const lyricContent = document.getElementById('lyricContent')
    if (!lyricContent || this.lyric.length === 0) return

    lyricContent.innerHTML = this.lyric.map((item, index) => 
      `<p class="lyric-line" data-index="${index}">${item.text}</p>`
    ).join('')
  }

  updateLyricPosition(currentTime) {
    if (this.lyric.length === 0) return

    let newIndex = -1
    for (let i = 0; i < this.lyric.length; i++) {
      if (this.lyric[i].time <= currentTime) {
        newIndex = i
      } else {
        break
      }
    }

    if (newIndex !== this.currentLyricIndex) {
      this.currentLyricIndex = newIndex
      this.highlightLyric(newIndex)
    }
  }

  highlightLyric(index) {
    const lines = document.querySelectorAll('.lyric-line')
    const container = document.getElementById('lyricContainer')
    
    lines.forEach((line, i) => {
      if (i === index) {
        line.classList.add('active')
      } else {
        line.classList.remove('active')
      }
    })

    if (index >= 0 && container) {
      const activeLine = lines[index]
      if (activeLine) {
        const containerHeight = container.clientHeight
        const lineTop = activeLine.offsetTop
        const lineHeight = activeLine.clientHeight
        const scrollTop = lineTop - containerHeight / 2 + lineHeight / 2
        
        container.scrollTo({
          top: scrollTop,
          behavior: 'smooth'
        })
      }
    }
  }
}

const playerDetail = new PlayerDetail()
window.playerDetail = playerDetail

export default playerDetail
