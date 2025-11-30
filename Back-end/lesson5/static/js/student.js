// 在所有fetch请求中添加headers
async function loadStudents() {
    try {
        const response = await fetch('/api/students', {
            headers: getAuthHeader()
        });
        // ... 其余代码不变
    } catch (error) {
        console.error('加载学生列表失败:', error);
    }
}

// 其他函数类似修改...