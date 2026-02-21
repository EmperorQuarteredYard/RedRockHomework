import header from './components/header/headbar.js'
import sidebar from './components/sidebar/sidebar.js'
import loginModal from './components/login/login.js'
import player from './components/player/player.js'
import playerDetail from './components/player/player-detail.js'
import homePage from './pages/home/home.js'
import playlistSquare from './pages/playlist/playlist-square.js'
import playlistDetail from './pages/playlist/playlist-detail.js'
import searchPage from './pages/search/search.js'
import './styles/main.css'

class App {
  constructor() {
    this.currentPage = 'home'
    this.currentParams = {}
  }

  render() {
    return `
      <div class="app-container">
        ${header.render()}
        <main class="main-content">
          ${sidebar.render()}
          <section class="content-area" id="contentArea">
            ${this.renderPage()}
          </section>
        </main>
        ${loginModal.render()}
        <div id="playerContainer"></div>
        <div id="playerDetailContainer"></div>
      </div>
    `
  }

  renderPage() {
    switch (this.currentPage) {
      case 'home':
        return homePage.render()
      case 'playlist-square':
        return playlistSquare.render()
      case 'playlist-detail':
        return playlistDetail.render()
      case 'search':
        return searchPage.render()
      default:
        return homePage.render()
    }
  }

  async navigateTo(page, params = {}) {
    this.currentPage = page
    this.currentParams = params
    const contentArea = document.getElementById('contentArea')
    
    if (contentArea) {
      contentArea.innerHTML = this.renderPage()
      await this.initPage(page, params)
    }
  }

  async initPage(page, params = {}) {
    switch (page) {
      case 'home':
        await homePage.loadData()
        break
      case 'playlist-square':
        await playlistSquare.loadData()
        break
      case 'playlist-detail':
        await playlistDetail.loadData(params.id)
        break
      case 'search':
        await searchPage.search(params.keyword)
        break
    }
  }

  init() {
    document.getElementById('app').innerHTML = this.render()
    
    header.init()
    sidebar.init()
    loginModal.init()
    
    const playerContainer = document.getElementById('playerContainer')
    if (playerContainer) {
      player.mount(playerContainer)
    }
    
    const playerDetailContainer = document.getElementById('playerDetailContainer')
    if (playerDetailContainer) {
      playerDetail.mount(playerDetailContainer)
    }
    
    this.initPage(this.currentPage)
    
    window.addEventListener('navigate', (e) => {
      const { page, ...params } = e.detail
      this.navigateTo(page, params)
    })
  }
}

const app = new App()

document.addEventListener('DOMContentLoaded', () => {
  app.init()
})

export default app
