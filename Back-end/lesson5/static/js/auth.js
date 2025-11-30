// 登录功能
document.getElementById('loginForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });

        const result = await response.json();

        if (response.ok) {
            // 保存Token到localStorage
            localStorage.setItem('token', result.token);
            localStorage.setItem('user', JSON.stringify(result.user));

            alert('登录成功！');
            window.location.href = '/dashboard';
        } else {
            alert('登录失败: ' + result.error);
        }
    } catch (error) {
        console.error('登录失败:', error);
        alert('登录失败！');
    }
});

// 注册功能
document.getElementById('registerForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (password !== confirmPassword) {
        alert('两次输入的密码不一致！');
        return;
    }

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, email, password })
        });

        const result = await response.json();

        if (response.ok) {
            alert('注册成功！请登录');
            window.location.href = '/login';
        } else {
            alert('注册失败: ' + result.error);
        }
    } catch (error) {
        console.error('注册失败:', error);
        alert('注册失败！');
    }
});

// 获取认证头
function getAuthHeader() {
    const token = localStorage.getItem('token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
}

// 检查登录状态
function checkAuth() {
    return localStorage.getItem('token') !== null;
}

// 获取当前用户
function getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

// 退出登录
function logout() {
    if (confirm('确定要退出登录吗？')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/';
    }
}

// 页面加载时检查认证状态
document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    const user = getCurrentUser();

    // 更新导航栏显示
    updateNavigation(user);
});

function updateNavigation(user) {
    // 这里可以根据用户状态更新页面导航
    // 具体实现取决于页面结构
}