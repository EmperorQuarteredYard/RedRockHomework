// 通用工具函数和API请求封装
class Auth {
    static getToken() {
        return localStorage.getItem('access_token');
    }

    static setToken(token) {
        localStorage.setItem('access_token', token);
    }

    static getRefreshToken() {
        return localStorage.getItem('refresh_token');
    }

    static setRefreshToken(token) {
        localStorage.setItem('refresh_token', token);
    }

    static clearTokens() {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_role');
    }

    static isAuthenticated() {
        return !!this.getToken();
    }

    static getUserRole() {
        return localStorage.getItem('user_role');
    }
}

class API {
    static async request(endpoint, options = {}) {
        const url = endpoint.startsWith('http') ? endpoint : `http://localhost:8080${endpoint}`;

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        const token = Auth.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers,
            });

            if (response.status === 401) {
                // Token expired, try to refresh
                const refreshed = await this.refreshToken();
                if (refreshed) {
                    return this.request(endpoint, options);
                } else {
                    window.location.href = '/login';
                    return null;
                }
            }

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || '请求失败');
            }

            return await response.json();
        } catch (error) {
            console.error('API请求错误:', error);
            UI.showError(error.message);
            throw error;
        }
    }

    static async refreshToken() {
        const refreshToken = Auth.getRefreshToken();
        if (!refreshToken) return false;

        try {
            const response = await fetch('http://localhost:8080/api/refresh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ refresh_token: refreshToken }),
            });

            if (response.ok) {
                const data = await response.json();
                Auth.setToken(data.access_token);
                Auth.setRefreshToken(data.refresh_token);
                return true;
            }
        } catch (error) {
            console.error('刷新Token失败:', error);
        }

        return false;
    }

    static async get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }

    static async post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    static async put(endpoint, data) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    static async delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    }
}

class UI {
    static showMessage(message, type = 'success') {
        // 移除现有的消息
        const existingMsg = document.querySelector('.alert-message');
        if (existingMsg) existingMsg.remove();

        const alert = document.createElement('div');
        alert.className = `alert alert-${type} alert-message`;
        alert.innerHTML = `
            <span>${message}</span>
            <button class="close-btn" onclick="this.parentElement.remove()">×</button>
        `;
        alert.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 9999;
            min-width: 300px;
            max-width: 500px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        `;

        document.body.appendChild(alert);

        // 自动消失
        setTimeout(() => {
            if (alert.parentElement) {
                alert.remove();
            }
        }, 5000);
    }

    static showError(message) {
        this.showMessage(message, 'error');
    }

    static showLoading() {
        const loading = document.createElement('div');
        loading.id = 'global-loading';
        loading.innerHTML = `
            <div class="loading-spinner"></div>
            <style>
                #global-loading {
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background: rgba(255,255,255,0.8);
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    z-index: 9999;
                }
                .loading-spinner {
                    border: 4px solid #f3f3f3;
                    border-top: 4px solid #3498db;
                    border-radius: 50%;
                    width: 40px;
                    height: 40px;
                    animation: spin 1s linear infinite;
                }
                @keyframes spin {
                    0% { transform: rotate(0deg); }
                    100% { transform: rotate(360deg); }
                }
            </style>
        `;
        document.body.appendChild(loading);
    }

    static hideLoading() {
        const loading = document.getElementById('global-loading');
        if (loading) loading.remove();
    }

    static confirm(message, callback) {
        const modal = document.createElement('div');
        modal.className = 'modal confirm-modal';
        modal.innerHTML = `
            <div class="modal-content" style="max-width: 400px;">
                <div class="modal-header">
                    <h3>确认操作</h3>
                </div>
                <div class="modal-body">
                    <p>${message}</p>
                </div>
                <div class="modal-footer" style="text-align: right; margin-top: 20px;">
                    <button class="btn btn-secondary cancel-btn">取消</button>
                    <button class="btn btn-danger confirm-btn">确认</button>
                </div>
            </div>
        `;

        document.body.appendChild(modal);
        modal.style.display = 'block';

        modal.querySelector('.cancel-btn').onclick = () => {
            modal.remove();
        };

        modal.querySelector('.confirm-btn').onclick = () => {
            callback();
            modal.remove();
        };

        modal.onclick = (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        };
    }

    static formatDate(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
        });
    }
}

class Modal {
    static open(title, content) {
        const modal = document.createElement('div');
        modal.className = 'modal';
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>${title}</h3>
                    <span class="close">&times;</span>
                </div>
                <div class="modal-body">
                    ${content}
                </div>
            </div>
        `;

        document.body.appendChild(modal);
        modal.style.display = 'block';

        modal.querySelector('.close').onclick = () => {
            modal.remove();
        };

        modal.onclick = (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        };

        return modal;
    }

    static close() {
        const modal = document.querySelector('.modal');
        if (modal) modal.remove();
    }
}

// 防抖函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 节流函数
function throttle(func, limit) {
    let inThrottle;
    return function() {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

// 导出到全局
window.Auth = Auth;
window.API = API;
window.UI = UI;
window.Modal = Modal;