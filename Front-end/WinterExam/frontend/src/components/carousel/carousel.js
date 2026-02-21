import api from '../../services/api.js'

class Carousel {
  constructor() {
    this.banners = []
    this.currentIndex = 0
    this.timer = null
    this.autoPlayInterval = 5000
  }

  render() {
    if (this.banners.length === 0) {
      return `<div class="carousel loading">加载中...</div>`
    }

    return `
      <div class="carousel" id="carousel">
        <div class="carousel-wrapper">
          <div class="carousel-track" id="carouselTrack">
            ${this.banners.map((banner, index) => `
              <div class="carousel-slide ${index === 0 ? 'active' : ''}" data-index="${index}">
                <img src="${banner.pic || banner.imageUrl || ''}" alt="${banner.typeTitle || '轮播图'}" onclick="window.carousel.handleClick(${index})" onerror="this.style.display='none'">
              </div>
            `).join('')}
          </div>
        </div>
        <button class="carousel-btn prev" onclick="window.carousel.prev()">❮</button>
        <button class="carousel-btn next" onclick="window.carousel.next()">❯</button>
        <div class="carousel-indicators">
          ${this.banners.map((_, index) => `
            <span class="indicator ${index === 0 ? 'active' : ''}" 
                  data-index="${index}" 
                  onclick="window.carousel.goTo(${index})"></span>
          `).join('')}
        </div>
      </div>
    `
  }

  async loadBanners() {
    try {
      const result = await api.getBanner()
      console.log('Banner result:', result)
      if (result.banners && result.banners.length > 0) {
        this.banners = result.banners
        return true
      }
    } catch (error) {
      console.error('加载轮播图失败:', error)
    }
    return false
  }

  goTo(index) {
    if (index < 0) {
      index = this.banners.length - 1
    } else if (index >= this.banners.length) {
      index = 0
    }

    this.currentIndex = index
    this.updateCarousel()
  }

  prev() {
    this.goTo(this.currentIndex - 1)
  }

  next() {
    this.goTo(this.currentIndex + 1)
  }

  updateCarousel() {
    const slides = document.querySelectorAll('.carousel-slide')
    const indicators = document.querySelectorAll('.indicator')
    
    slides.forEach((slide, index) => {
      if (index === this.currentIndex) {
        slide.classList.add('active')
      } else {
        slide.classList.remove('active')
      }
    })
    
    indicators.forEach((indicator, index) => {
      if (index === this.currentIndex) {
        indicator.classList.add('active')
      } else {
        indicator.classList.remove('active')
      }
    })
  }

  startAutoPlay() {
    this.stopAutoPlay()
    this.timer = setInterval(() => {
      this.next()
    }, this.autoPlayInterval)
  }

  stopAutoPlay() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
  }

  init() {
    const carousel = document.getElementById('carousel')
    if (carousel) {
      carousel.addEventListener('mouseenter', () => this.stopAutoPlay())
      carousel.addEventListener('mouseleave', () => this.startAutoPlay())
    }
    this.startAutoPlay()
  }

  handleClick(index) {
    const banner = this.banners[index]
    if (!banner) return
    
    if (banner.targetId && banner.targetType === 1) {
      window.playlistCard?.openPlaylist(banner.targetId)
    } else if (banner.url) {
      window.open(banner.url, '_blank')
    }
  }
}

const carousel = new Carousel()
window.carousel = carousel

export default carousel
