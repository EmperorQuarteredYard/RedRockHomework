document.addEventListener('DOMContentLoaded', function() {
    loadLessons();

    document.getElementById('addLessonForm').addEventListener('submit', function(e) {
        e.preventDefault();
        addLesson();
    });
});

async function loadLessons() {
    try {
        const response = await fetch('/api/lessons');
        const lessons = await response.json();

        const tbody = document.querySelector('#lessonsTable tbody');
        tbody.innerHTML = '';

        lessons.forEach(lesson => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${lesson.ID}</td>
                <td>${lesson.Name}</td>
                <td>${lesson.Code}</td>
                <td>${lesson.Credit}</td>
                <td>${lesson.Capacity}</td>
                <td>
                    <button class="delete-btn" onclick="deleteLesson(${lesson.ID})">删除</button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error('加载课程列表失败:', error);
    }
}

async function addLesson() {
    const name = document.getElementById('lessonName').value;
    const code = document.getElementById('lessonCode').value;
    const credit = document.getElementById('lessonCredit').value;
    const capacity = document.getElementById('lessonCapacity').value;
    const description = document.getElementById('lessonDescription').value;

    try {
        const response = await fetch('/api/lessons', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Name: name,
                Code: code,
                Credit: parseInt(credit),
                Capacity: parseInt(capacity),
                Description: description
            })
        });

        if (response.ok) {
            document.getElementById('addLessonForm').reset();
            loadLessons();
            alert('课程添加成功！');
        } else {
            alert('添加失败！');
        }
    } catch (error) {
        console.error('添加课程失败:', error);
        alert('添加失败！');
    }
}

async function deleteLesson(id) {
    if (!confirm('确定要删除这个课程吗？')) {
        return;
    }

    try {
        const response = await fetch(`/api/lessons/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            loadLessons();
            alert('课程删除成功！');
        } else {
            alert('删除失败！');
        }
    } catch (error) {
        console.error('删除课程失败:', error);
        alert('删除失败！');
    }
}