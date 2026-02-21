import store from '../../store/store.js'
import playerDetail from './player-detail.js'

class Player {
  constructor() {
    this.container = null
    this.unsubscribe = null
  }

  render() {
    return `
      <div class="player-bar" id="playerBar">
        <div class="player-left">
          <div class="player-cover" id="playerCover" title="点击查看详情">
            <img src="" alt="封面" id="playerCoverImg">
          </div>
          <div class="player-info">
            <div class="player-song-name" id="playerSongName">未播放歌曲</div>
            <div class="player-artist" id="playerArtist">-</div>
          </div>
        </div>
        <div class="player-center">
          <div class="player-controls">
            <button class="control-btn" id="prevBtn" title="上一首">
              <svg viewBox="0 0 24 24" width="20" height="20">
                <path fill="currentColor" d="M6 6h2v12H6zm3.5 6l8.5 6V6z"/>
              </svg>
            </button>
            <button class="control-btn play-btn" id="playBtn" title="播放/暂停">
              <svg viewBox="0 0 24 24" width="28" height="28" id="playIcon">
                <path fill="currentColor" d="M8 5v14l11-7z"/>
              </svg>
            </button>
            <button class="control-btn" id="nextBtn" title="下一首">
              <svg viewBox="0 0 24 24" width="20" height="20">
                <path fill="currentColor" d="M6 18l8.5-6L6 6v12zM16 6v12h2V6h-2z"/>
              </svg>
            </button>
          </div>
          <div class="player-progress">
            <span class="time" id="currentTime">00:00</span>
            <div class="progress-bar" id="progressBar">
              <div class="progress-bg"></div>
              <div class="progress-current" id="progressCurrent"></div>
              <div class="progress-thumb" id="progressThumb"></div>
            </div>
            <span class="time" id="totalTime">00:00</span>
          </div>
        </div>
        <div class="player-right">
          <div class="volume-control">
            <button class="control-btn" id="volumeBtn" title="音量">
              <svg viewBox="0 0 24 24" width="18" height="18">
                <path fill="currentColor" d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
              </svg>
            </button>
            <div class="volume-bar" id="volumeBar">
              <div class="volume-bg"></div>
              <div class="volume-current" id="volumeCurrent"></div>
              <div class="volume-thumb" id="volumeThumb"></div>
            </div>
          </div>
          <button class="control-btn" id="playlistBtn" title="播放列表">
            <svg viewBox="0 0 24 24" width="18" height="18">
              <path fill="currentColor" d="M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z"/>
            </svg>
          </button>
        </div>
      </div>
    `
  }

  mount(container) {
    this.container = container
    this.container.innerHTML = this.render()
    this.bindEvents()
    this.unsubscribe = store.subscribe((state) => this.updateUI(state))
    this.updateUI(store.getState())
    this.initVolumeBar()
  }

  bindEvents() {
    const playBtn = document.getElementById('playBtn')
    const prevBtn = document.getElementById('prevBtn')
    const nextBtn = document.getElementById('nextBtn')
    const progressBar = document.getElementById('progressBar')
    const volumeBar = document.getElementById('volumeBar')
    const playerCover = document.getElementById('playerCover')

    playBtn?.addEventListener('click', () => store.togglePlay())
    prevBtn?.addEventListener('click', () => store.prev())
    nextBtn?.addEventListener('click', () => store.next())

    playerCover?.addEventListener('click', () => {
      playerDetail.toggle()
    })

    progressBar?.addEventListener('click', (e) => {
      const rect = progressBar.getBoundingClientRect()
      const percent = (e.clientX - rect.left) / rect.width
      const time = percent * store.duration
      store.seek(time)
    })

    volumeBar?.addEventListener('click', (e) => {
      const rect = volumeBar.getBoundingClientRect()
      const percent = (e.clientX - rect.left) / rect.width
      store.setVolume(percent)
    })

    const volumeBtn = document.getElementById('volumeBtn')
    volumeBtn?.addEventListener('click', () => {
      if (store.volume > 0) {
        this.prevVolume = store.volume
        store.setVolume(0)
      } else {
        store.setVolume(this.prevVolume || 0.8)
      }
    })
  }

  initVolumeBar() {
    const volumeCurrent = document.getElementById('volumeCurrent')
    const volumeThumb = document.getElementById('volumeThumb')
    if (volumeCurrent) volumeCurrent.style.width = `${store.volume * 100}%`
    if (volumeThumb) volumeThumb.style.left = `${store.volume * 100}%`
  }

  updateUI(state) {
    const { currentSong, isPlaying, currentTime, duration, volume } = state

    const songName = document.getElementById('playerSongName')
    const artist = document.getElementById('playerArtist')
    const coverImg = document.getElementById('playerCoverImg')
    const playIcon = document.getElementById('playIcon')
    const currentTimeEl = document.getElementById('currentTime')
    const totalTimeEl = document.getElementById('totalTime')
    const progressCurrent = document.getElementById('progressCurrent')
    const progressThumb = document.getElementById('progressThumb')
    const volumeCurrent = document.getElementById('volumeCurrent')
    const volumeThumb = document.getElementById('volumeThumb')

    if (currentSong) {
      if (songName) songName.textContent = currentSong.name || '未知歌曲'
      if (artist) {
        const artists = currentSong.ar || currentSong.artists || []
        artist.textContent = artists.map(a => a.name).join('/') || '-'
      }
      if (coverImg) {
        const picUrl = currentSong.al?.picUrl || currentSong.album?.picUrl || currentSong.picUrl || ''
        coverImg.src = picUrl || ''
      }
    } else {
      if (songName) songName.textContent = '未播放歌曲'
      if (artist) artist.textContent = '-'
      if (coverImg) coverImg.src = ''
    }

    if (playIcon) {
      if (isPlaying) {
        playIcon.innerHTML = '<path fill="currentColor" d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>'
      } else {
        playIcon.innerHTML = '<path fill="currentColor" d="M8 5v14l11-7z"/>'
      }
    }

    if (currentTimeEl) currentTimeEl.textContent = store.formatTime(currentTime)
    if (totalTimeEl) totalTimeEl.textContent = store.formatTime(duration)

    const progress = duration ? (currentTime / duration) * 100 : 0
    if (progressCurrent) progressCurrent.style.width = `${progress}%`
    if (progressThumb) progressThumb.style.left = `${progress}%`

    if (volumeCurrent) volumeCurrent.style.width = `${volume * 100}%`
    if (volumeThumb) volumeThumb.style.left = `${volume * 100}%`
  }

  destroy() {
    if (this.unsubscribe) {
      this.unsubscribe()
    }
  }
}

const player = new Player()
export default player
