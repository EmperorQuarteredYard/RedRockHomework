import api from '../../services/api.js'

class LoginModal {
  constructor() {
    this.isVisible = false
    this.countdown = 0
    this.timer = null
    this.qrTimer = null
    this.qrKey = null
    this.onLoginSuccess = null
  }

  render() {
    return `
      <div class="login-modal" id="loginModal">
        <div class="login-overlay" onclick="window.loginModal.hide()"></div>
        <div class="login-content">
          <div class="login-header">
            <h2>登录</h2>
            <button class="close-btn" onclick="window.loginModal.hide()">×</button>
          </div>
          <div class="login-body">
            <div class="login-tabs">
              <div class="login-tab active" data-tab="qrcode">二维码登录</div>
              <div class="login-tab" data-tab="captcha">验证码登录</div>
              <div class="login-tab" data-tab="password">密码登录</div>
            </div>
            <div class="login-form" id="qrcodeForm">
              <div class="qrcode-container">
                <div class="qrcode-tip">使用网易云音乐APP扫码登录</div>
                <div class="qrcode-wrapper" id="qrcodeWrapper">
                  <img id="qrcodeImg" src="" alt="登录二维码">
                  <div class="qrcode-loading" id="qrcodeLoading">
                    <p>加载中...</p>
                  </div>
                  <div class="qrcode-expired hidden" id="qrcodeExpired">
                    <p>二维码已过期</p>
                    <button onclick="window.loginModal.refreshQRCode()">点击刷新</button>
                  </div>
                </div>
                <div class="qrcode-status" id="qrcodeStatus">正在获取二维码...</div>
              </div>
            </div>
            <div class="login-form hidden" id="captchaForm">
              <div class="form-group">
                <input type="text" id="phoneInput" placeholder="请输入手机号" maxlength="11">
              </div>
              <div class="form-group captcha-group">
                <input type="text" id="captchaInput" placeholder="请输入验证码" maxlength="6">
                <button class="send-captcha-btn" id="sendCaptchaBtn" onclick="window.loginModal.sendCaptcha()">获取验证码</button>
              </div>
              <div class="form-tip">验证码登录需要网易云音乐服务器支持，如无法收到验证码请使用二维码登录</div>
              <button class="login-btn" onclick="window.loginModal.loginWithCaptcha()">登录</button>
            </div>
            <div class="login-form hidden" id="passwordForm">
              <div class="form-group">
                <input type="text" id="phoneInputPwd" placeholder="请输入手机号" maxlength="11">
              </div>
              <div class="form-group">
                <input type="password" id="passwordInput" placeholder="请输入密码">
              </div>
              <button class="login-btn" onclick="window.loginModal.loginWithPassword()">登录</button>
            </div>
          </div>
        </div>
      </div>
    `
  }

  show() {
    this.isVisible = true
    const modal = document.getElementById('loginModal')
    if (modal) {
      modal.classList.add('show')
    }
    this.initQRCodeLogin()
  }

  hide() {
    this.isVisible = false
    const modal = document.getElementById('loginModal')
    if (modal) {
      modal.classList.remove('show')
    }
    this.stopQRCodeCheck()
    this.resetForm()
  }

  resetForm() {
    const phoneInput = document.getElementById('phoneInput')
    const captchaInput = document.getElementById('captchaInput')
    const phoneInputPwd = document.getElementById('phoneInputPwd')
    const passwordInput = document.getElementById('passwordInput')
    
    if (phoneInput) phoneInput.value = ''
    if (captchaInput) captchaInput.value = ''
    if (phoneInputPwd) phoneInputPwd.value = ''
    if (passwordInput) passwordInput.value = ''
    
    this.stopCountdown()
    this.qrKey = null
  }

  async initQRCodeLogin() {
    const loading = document.getElementById('qrcodeLoading')
    const expired = document.getElementById('qrcodeExpired')
    const status = document.getElementById('qrcodeStatus')
    const qrImg = document.getElementById('qrcodeImg')
    
    if (loading) loading.classList.remove('hidden')
    if (expired) expired.classList.add('hidden')
    if (status) status.textContent = '正在获取二维码...'
    if (qrImg) qrImg.src = ''
    
    try {
      const keyResult = await api.getQRKey()
      console.log('QR Key Result:', keyResult)
      
      let unikey = null
      if (keyResult.code === 200) {
        unikey = keyResult.data?.unikey || keyResult.body?.data?.unikey || keyResult.unikey
      }
      
      if (!unikey) {
        throw new Error('获取二维码Key失败: ' + JSON.stringify(keyResult))
      }
      
      this.qrKey = unikey
      const qrResult = await api.getQRCode(this.qrKey, true)
      console.log('QR Code Result:', qrResult)
      
      let qrimg = null
      let qrurl = null
      
      if (qrResult.code === 200) {
        qrimg = qrResult.data?.qrimg || qrResult.body?.data?.qrimg || qrResult.qrimg
        qrurl = qrResult.data?.qrurl || qrResult.body?.data?.qrurl || qrResult.qrurl
      }
      
      console.log('qrimg:', qrimg ? 'exists' : 'null')
      console.log('qrurl:', qrurl)
      
      if (qrimg) {
        console.log('Using base64 QR image')
        if (qrImg) {
          qrImg.src = qrimg
          qrImg.style.display = 'block'
          qrImg.style.visibility = 'visible'
        }
      } else if (qrurl) {
        console.log('Using QR URL to generate image')
        if (qrImg) {
          qrImg.src = `https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(qrurl)}`
          qrImg.style.display = 'block'
          qrImg.style.visibility = 'visible'
        }
      } else {
        throw new Error('获取二维码图片失败: ' + JSON.stringify(qrResult))
      }
      
      if (loading) {
        loading.classList.add('hidden')
        loading.style.display = 'none'
      }
      if (expired) {
        expired.classList.add('hidden')
        expired.style.display = 'none'
      }
      if (status) status.textContent = '请使用网易云音乐APP扫码'
      this.startQRCodeCheck()
      
    } catch (error) {
      console.error('获取二维码失败:', error)
      if (loading) loading.classList.add('hidden')
      if (status) status.textContent = '获取二维码失败: ' + error.message
      if (expired) {
        expired.classList.remove('hidden')
        expired.querySelector('p').textContent = '获取失败，点击重试'
      }
    }
  }

  async refreshQRCode() {
    this.stopQRCodeCheck()
    await this.initQRCodeLogin()
  }

  startQRCodeCheck() {
    this.stopQRCodeCheck()
    
    this.qrTimer = setInterval(async () => {
      if (!this.qrKey) {
        this.stopQRCodeCheck()
        return
      }
      
      try {
        const result = await api.checkQRCode(this.qrKey)
        console.log('QR Check Result:', result)
        
        const status = document.getElementById('qrcodeStatus')
        const expired = document.getElementById('qrcodeExpired')
        
        if (result.code === 800 || result.code === 801 && result.message?.includes('过期')) {
          if (status) status.textContent = '二维码已过期'
          if (expired) expired.classList.remove('hidden')
          this.stopQRCodeCheck()
        } else if (result.code === 803) {
          if (status) status.textContent = '登录成功！'
          this.stopQRCodeCheck()
          
          setTimeout(async () => {
            this.hide()
            if (this.onLoginSuccess) {
              try {
                const profileResult = await api.getLoginStatus()
                console.log('Login Status:', profileResult)
                if (profileResult.data?.profile) {
                  this.onLoginSuccess(profileResult.data.profile)
                }
              } catch (e) {
                console.error('获取用户信息失败:', e)
              }
            }
          }, 500)
        } else if (result.code === 802) {
          if (status) status.textContent = '请在手机上确认登录'
        } else if (result.code === 801) {
          if (status) status.textContent = '等待扫码...'
        } else {
          console.log('Unknown QR code status:', result)
        }
      } catch (error) {
        console.error('检查二维码状态失败:', error)
      }
    }, 3000)
  }

  stopQRCodeCheck() {
    if (this.qrTimer) {
      clearInterval(this.qrTimer)
      this.qrTimer = null
    }
  }

  async sendCaptcha() {
    const phone = document.getElementById('phoneInput')?.value
    if (!phone || !/^1[3-9]\d{9}$/.test(phone)) {
      alert('请输入正确的手机号')
      return
    }

    if (this.countdown > 0) return

    try {
      const result = await api.sendCaptcha(phone)
      if (result.code === 200) {
        alert('验证码发送成功，请注意查收')
        this.startCountdown()
      } else {
        alert(result.message || '验证码发送失败，请尝试二维码登录')
      }
    } catch (error) {
      console.error('发送验证码失败:', error)
      alert('验证码发送失败，请使用二维码登录')
    }
  }

  startCountdown() {
    this.countdown = 60
    const btn = document.getElementById('sendCaptchaBtn')
    
    this.timer = setInterval(() => {
      this.countdown--
      if (btn) {
        btn.textContent = `${this.countdown}s后重发`
        btn.disabled = true
      }
      
      if (this.countdown <= 0) {
        this.stopCountdown()
      }
    }, 1000)
  }

  stopCountdown() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
    this.countdown = 0
    const btn = document.getElementById('sendCaptchaBtn')
    if (btn) {
      btn.textContent = '获取验证码'
      btn.disabled = false
    }
  }

  async loginWithCaptcha() {
    const phone = document.getElementById('phoneInput')?.value
    const captcha = document.getElementById('captchaInput')?.value

    if (!phone || !/^1[3-9]\d{9}$/.test(phone)) {
      alert('请输入正确的手机号')
      return
    }

    if (!captcha || !/^\d{4,6}$/.test(captcha)) {
      alert('请输入正确的验证码')
      return
    }

    try {
      const result = await api.loginWithPhone(phone, captcha)
      if (result.code === 200) {
        alert('登录成功')
        this.hide()
        if (this.onLoginSuccess) {
          this.onLoginSuccess(result.profile)
        }
      } else {
        alert(result.message || '登录失败')
      }
    } catch (error) {
      console.error('登录失败:', error)
      alert('登录失败，请稍后重试')
    }
  }

  async loginWithPassword() {
    const phone = document.getElementById('phoneInputPwd')?.value
    const password = document.getElementById('passwordInput')?.value

    if (!phone || !/^1[3-9]\d{9}$/.test(phone)) {
      alert('请输入正确的手机号')
      return
    }

    if (!password) {
      alert('请输入密码')
      return
    }

    try {
      const result = await api.loginWithPassword(phone, password)
      if (result.code === 200) {
        alert('登录成功')
        this.hide()
        if (this.onLoginSuccess) {
          this.onLoginSuccess(result.profile)
        }
      } else {
        alert(result.message || '登录失败')
      }
    } catch (error) {
      console.error('登录失败:', error)
      alert('登录失败，请稍后重试')
    }
  }

  switchTab(tab) {
    const tabs = document.querySelectorAll('.login-tab')
    const qrcodeForm = document.getElementById('qrcodeForm')
    const captchaForm = document.getElementById('captchaForm')
    const passwordForm = document.getElementById('passwordForm')
    
    tabs.forEach(t => t.classList.remove('active'))
    document.querySelector(`.login-tab[data-tab="${tab}"]`)?.classList.add('active')
    
    qrcodeForm?.classList.add('hidden')
    captchaForm?.classList.add('hidden')
    passwordForm?.classList.add('hidden')
    
    if (tab === 'qrcode') {
      qrcodeForm?.classList.remove('hidden')
      if (!this.qrKey) {
        this.initQRCodeLogin()
      } else {
        this.startQRCodeCheck()
      }
    } else if (tab === 'captcha') {
      captchaForm?.classList.remove('hidden')
      this.stopQRCodeCheck()
    } else {
      passwordForm?.classList.remove('hidden')
      this.stopQRCodeCheck()
    }
  }

  init() {
    document.addEventListener('click', (e) => {
      if (e.target.classList.contains('login-tab')) {
        const tab = e.target.dataset.tab
        this.switchTab(tab)
      }
    })
  }
}

const loginModal = new LoginModal()
window.loginModal = loginModal

export default loginModal
