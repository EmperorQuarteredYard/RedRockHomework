class Sidebar {
  constructor() {
    this.currentNav = 'home'
    this.navItems = [
      { id: 'home', name: '发现音乐', icon: '🎵' },
      { id: 'playlist-square', name: '歌单广场', icon: '📋' },
      { id: 'liked', name: '我喜欢的音乐', icon: '❤️' }
    ]
  }

  render() {
    return `
      <aside class="sidebar">
        <nav class="nav-menu">
          ${this.navItems.map(item => `
            <div class="nav-item ${this.currentNav === item.id ? 'active' : ''}" 
                 data-nav="${item.id}" 
                 onclick="window.sidebar.navigate('${item.id}')">
              <span class="nav-icon">${item.icon}</span>
              <span class="nav-text">${item.name}</span>
            </div>
          `).join('')}
        </nav>
      </aside>
    `
  }

  navigate(navId) {
    this.currentNav = navId
    this.updateActiveNav()
    
    const event = new CustomEvent('navigate', { 
      detail: { page: navId } 
    })
    window.dispatchEvent(event)
  }

  updateActiveNav() {
    const navItems = document.querySelectorAll('.nav-item')
    navItems.forEach(item => {
      if (item.dataset.nav === this.currentNav) {
        item.classList.add('active')
      } else {
        item.classList.remove('active')
      }
    })
  }

  init() {
    this.updateActiveNav()
  }
}

const sidebar = new Sidebar()
window.sidebar = sidebar

export default sidebar
