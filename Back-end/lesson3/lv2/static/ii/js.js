
function setExample1() {
    document.getElementById('jsonInput').value = '{"name":"小明","score":[85.5,92.0,78.5,88.0,91.5]}';
}

function setExample2() {
    document.getElementById('jsonInput').value = '{"name":"小红","score":[95.0,87.5,92.5,90.0,88.5]}';
}

function setExample3() {
    document.getElementById('jsonInput').value = '{"name":"小刚","score":[68.0,72.5,75.0,80.5,69.5]}';
}

// 清空输入
function clearInput() {
    document.getElementById('jsonInput').value = '';
    document.getElementById('result').style.display = 'none';
    document.getElementById('error').style.display = 'none';
}

// 计算平均成绩
function calculateAverage() {
    // 隐藏之前的错误和结果
    document.getElementById('result').style.display = 'none';
    document.getElementById('error').style.display = 'none';

    // 获取输入的JSON
    const jsonInput = document.getElementById('jsonInput').value.trim();

    if (!jsonInput) {
        showError('请输入JSON数据');
        return;
    }

    // 验证JSON格式
    let jsonData;
    try {
        jsonData = JSON.parse(jsonInput);
    } catch (e) {
        showError('JSON格式错误: ' + e.message);
        return;
    }

    // 验证必需字段
    if (!jsonData.name || !jsonData.score || !Array.isArray(jsonData.score)) {
        showError('JSON数据格式不正确，需要包含name和score数组字段');
        return;
    }

    // 发送请求到后端
    fetch('/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: jsonInput
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('服务器响应错误: ' + response.status);
            }
            return response.json();
        })
        .then(data => {
            // 显示结果
            document.getElementById('resultContent').textContent = JSON.stringify(data, null, 2);
            document.getElementById('result').style.display = 'block';

            // 滚动到结果位置
            document.getElementById('result').scrollIntoView({ behavior: 'smooth' });
        })
        .catch(error => {
            showError('请求失败: ' + error.message);
        });
}

// 显示错误信息
function showError(message) {
    document.getElementById('errorContent').textContent = message;
    document.getElementById('error').style.display = 'block';

    // 滚动到错误位置
    document.getElementById('error').scrollIntoView({ behavior: 'smooth' });
}

// 页面加载时设置示例1
window.onload = function() {
    setExample1();
};