// 控制台页面功能
document.addEventListener('DOMContentLoaded', function() {
    loadDashboardStats();
});

async function loadDashboardStats() {
    if (!checkAuth()) return;

    try {
        const headers = getAuthHeader();

        const [lessonsResponse, profileResponse] = await Promise.all([
            fetch('/api/lessons', { headers }),
            fetch('/api/profile', { headers })
        ]);

        if (lessonsResponse.ok) {
            const lessons = await lessonsResponse.json();
            // 可以在这里更新课程统计信息
        }

        if (profileResponse.ok) {
            const profile = await profileResponse.json();
            // 可以在这里更新用户信息显示
        }
    } catch (error) {
        console.error('加载控制台数据失败:', error);
    }
}