document.addEventListener('DOMContentLoaded', function() {
    loadStudentsAndLessons();

    document.getElementById('enrollForm').addEventListener('submit', function(e) {
        e.preventDefault();
        enrollStudent();
    });

    document.getElementById('unenrollForm').addEventListener('submit', function(e) {
        e.preventDefault();
        unenrollStudent();
    });
});

async function loadStudentsAndLessons() {
    try {
        const [studentsResponse, lessonsResponse] = await Promise.all([
            fetch('/api/students'),
            fetch('/api/lessons')
        ]);

        const students = await studentsResponse.json();
        const lessons = await lessonsResponse.json();

        // 填充所有下拉框
        populateSelect('enrollStudent', students);
        populateSelect('unenrollStudent', students);
        populateSelect('viewStudent', students);
        populateSelect('enrollLesson', lessons);
        populateSelect('unenrollLesson', lessons);
        populateSelect('viewLesson', lessons);
    } catch (error) {
        console.error('加载数据失败:', error);
    }
}

function populateSelect(selectId, items) {
    const select = document.getElementById(selectId);
    // 保留第一个选项
    select.innerHTML = select.options[0] ? select.options[0].outerHTML : '<option value="">选择</option>';

    items.forEach(item => {
        const option = document.createElement('option');
        option.value = item.ID;
        option.textContent = `${item.Name} (${item.Email || item.Code})`;
        select.appendChild(option);
    });
}

async function enrollStudent() {
    const studentId = document.getElementById('enrollStudent').value;
    const lessonId = document.getElementById('enrollLesson').value;

    if (!studentId || !lessonId) {
        alert('请选择学生和课程！');
        return;
    }

    try {
        const response = await fetch('/api/enroll', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                student_id: parseInt(studentId),
                lesson_id: parseInt(lessonId)
            })
        });

        const result = await response.json();

        if (response.ok) {
            alert('选课成功！');
            document.getElementById('enrollForm').reset();
        } else {
            alert('选课失败: ' + result.error);
        }
    } catch (error) {
        console.error('选课失败:', error);
        alert('选课失败！');
    }
}

async function unenrollStudent() {
    const studentId = document.getElementById('unenrollStudent').value;
    const lessonId = document.getElementById('unenrollLesson').value;

    if (!studentId || !lessonId) {
        alert('请选择学生和课程！');
        return;
    }

    try {
        const response = await fetch('/api/enroll', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                student_id: parseInt(studentId),
                lesson_id: parseInt(lessonId)
            })
        });

        const result = await response.json();

        if (response.ok) {
            alert('退选成功！');
            document.getElementById('unenrollForm').reset();
        } else {
            alert('退选失败: ' + result.error);
        }
    } catch (error) {
        console.error('退选失败:', error);
        alert('退选失败！');
    }
}

async function loadStudentLessons() {
    const studentId = document.getElementById('viewStudent').value;
    if (!studentId) return;

    try {
        const response = await fetch(`/api/students/${studentId}/lessons`);
        const lessons = await response.json();

        const container = document.getElementById('studentLessons');
        if (lessons.length === 0) {
            container.innerHTML = '<p>该学生没有选择任何课程</p>';
        } else {
            let html = '<div class="lesson-list">';
            lessons.forEach(lesson => {
                html += `<div class="lesson-item">${lesson.Name} (${lesson.Code}) - 学分: ${lesson.Credit}</div>`;
            });
            html += '</div>';
            container.innerHTML = html;
        }
    } catch (error) {
        console.error('加载学生课程失败:', error);
    }
}

async function loadLessonStudents() {
    const lessonId = document.getElementById('viewLesson').value;
    if (!lessonId) return;

    try {
        const response = await fetch(`/api/lessons/${lessonId}/students`);
        const students = await response.json();

        const container = document.getElementById('lessonStudents');
        if (students.length === 0) {
            container.innerHTML = '<p>该课程没有学生选择</p>';
        } else {
            let html = '<div class="student-list">';
            students.forEach(student => {
                html += `<div class="student-item">${student.Name} (${student.Email})</div>`;
            });
            html += '</div>';
            container.innerHTML = html;
        }
    } catch (error) {
        console.error('加载课程学生失败:', error);
    }
}