class MusicStore {
  constructor() {
    this.currentSong = null
    this.playlist = []
    this.currentIndex = -1
    this.isPlaying = false
    this.currentTime = 0
    this.duration = 0
    this.volume = 0.8
    this.listeners = new Set()
    this.audio = new Audio()
    this.audio.volume = this.volume
    
    this.setupAudioEvents()
  }

  setupAudioEvents() {
    this.audio.addEventListener('timeupdate', () => {
      this.currentTime = this.audio.currentTime
      this.duration = this.audio.duration || 0
      this.notify()
    })

    this.audio.addEventListener('ended', () => {
      this.next()
    })

    this.audio.addEventListener('error', (e) => {
      console.error('Audio error:', e)
      this.isPlaying = false
      this.notify()
    })

    this.audio.addEventListener('play', () => {
      this.isPlaying = true
      this.notify()
    })

    this.audio.addEventListener('pause', () => {
      this.isPlaying = false
      this.notify()
    })
  }

  subscribe(listener) {
    this.listeners.add(listener)
    return () => this.listeners.delete(listener)
  }

  notify() {
    this.listeners.forEach(listener => listener(this.getState()))
  }

  getState() {
    return {
      currentSong: this.currentSong,
      playlist: this.playlist,
      currentIndex: this.currentIndex,
      isPlaying: this.isPlaying,
      currentTime: this.currentTime,
      duration: this.duration,
      volume: this.volume
    }
  }

  async setPlaylist(songs, index = 0) {
    this.playlist = songs
    this.currentIndex = index
    if (songs.length > 0) {
      await this.playSong(songs[index], index)
    }
  }

  async playSong(song, index) {
    if (!song) return
    
    this.currentSong = song
    this.currentIndex = index
    
    try {
      const api = (await import('../services/api.js')).default
      const result = await api.getSongUrl(song.id)
      
      if (result.code === 200 && result.data && result.data[0]) {
        const url = result.data[0].url
        if (url) {
          this.audio.src = url
          this.audio.play()
        } else {
          console.error('No URL for song:', song.name)
        }
      }
    } catch (error) {
      console.error('Failed to get song URL:', error)
    }
    
    this.notify()
  }

  play() {
    if (this.audio.src) {
      this.audio.play()
    }
  }

  pause() {
    this.audio.pause()
  }

  togglePlay() {
    if (this.isPlaying) {
      this.pause()
    } else {
      this.play()
    }
  }

  async next() {
    if (this.playlist.length === 0) return
    
    let nextIndex = this.currentIndex + 1
    if (nextIndex >= this.playlist.length) {
      nextIndex = 0
    }
    
    await this.playSong(this.playlist[nextIndex], nextIndex)
  }

  async prev() {
    if (this.playlist.length === 0) return
    
    let prevIndex = this.currentIndex - 1
    if (prevIndex < 0) {
      prevIndex = this.playlist.length - 1
    }
    
    await this.playSong(this.playlist[prevIndex], prevIndex)
  }

  seek(time) {
    if (this.audio.duration) {
      this.audio.currentTime = time
    }
  }

  setVolume(volume) {
    this.volume = Math.max(0, Math.min(1, volume))
    this.audio.volume = this.volume
    this.notify()
  }

  formatTime(seconds) {
    if (!seconds || isNaN(seconds)) return '00:00'
    const mins = Math.floor(seconds / 60)
    const secs = Math.floor(seconds % 60)
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
}

const store = new MusicStore()
export default store
