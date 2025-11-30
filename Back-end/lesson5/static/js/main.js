// 首页统计信息
document.addEventListener('DOMContentLoaded', function() {
    loadStats();
});

async function loadStats() {
    try {
        const [studentsResponse, lessonsResponse] = await Promise.all([
            fetch('/api/students'),
            fetch('/api/lessons')
        ]);

        const students = await studentsResponse.json();
        const lessons = await lessonsResponse.json();

        document.getElementById('studentCount').textContent = students.length;
        document.getElementById('lessonCount').textContent = lessons.length;
    } catch (error) {
        console.error('加载统计信息失败:', error);
    }
}