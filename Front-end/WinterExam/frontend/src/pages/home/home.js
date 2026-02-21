import api from '../../services/api.js'
import carousel from '../../components/carousel/carousel.js'
import playlistCard from '../../components/playlist/playlist-card.js'

class HomePage {
  constructor() {
    this.recommendPlaylists = []
  }

  render() {
    return `
      <div class="home-page">
        <div class="carousel-section" id="carouselSection">
          ${carousel.render()}
        </div>
        <div class="recommend-section">
          <div class="section-header">
            <h3>推荐歌单</h3>
            <button class="more-btn" onclick="window.homePage.goToPlaylistSquare()">更多 ></button>
          </div>
          <div class="recommend-list" id="recommendList">
            <div class="loading">加载中...</div>
          </div>
        </div>
      </div>
    `
  }

  async loadData() {
    await this.loadCarousel()
    await this.loadRecommendPlaylists()
  }

  async loadCarousel() {
    const loaded = await carousel.loadBanners()
    if (loaded) {
      const section = document.getElementById('carouselSection')
      if (section) {
        section.innerHTML = carousel.render()
        carousel.init()
      }
    }
  }

  async loadRecommendPlaylists() {
    try {
      const result = await api.getPersonalized(10)
      if (result.code === 200 && result.result) {
        this.recommendPlaylists = result.result
        this.renderRecommendList()
      }
    } catch (error) {
      console.error('加载推荐歌单失败:', error)
      const container = document.getElementById('recommendList')
      if (container) {
        container.innerHTML = '<div class="error">加载失败，请稍后重试</div>'
      }
    }
  }

  renderRecommendList() {
    const container = document.getElementById('recommendList')
    if (container && this.recommendPlaylists.length > 0) {
      container.innerHTML = playlistCard.renderList(this.recommendPlaylists)
    }
  }

  goToPlaylistSquare() {
    const event = new CustomEvent('navigate', { 
      detail: { page: 'playlist-square' } 
    })
    window.dispatchEvent(event)
  }
}

const homePage = new HomePage()
window.homePage = homePage

export default homePage
