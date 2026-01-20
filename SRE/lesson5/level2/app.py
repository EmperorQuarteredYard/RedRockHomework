import objects
from flask import Flask, request, jsonify, session

# 创建Flask应用实例
app = Flask(__name__)
app.secret_key = 'your_secret_key_here'  # 用于会话加密

# 简单的用户存储（实际应用中应使用数据库）
users_db = {}

# 创建留言板实例
board = objects.commitboard()

@app.route('/register', methods=['POST'])
def register():
    """用户注册"""
    if not request.is_json:
        return jsonify({'error': '请求格式必须为JSON'}), 400
    
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')
    
    if not username or not password:
        return jsonify({'error': '用户名和密码不能为空'}), 400
    
    if username in users_db:
        return jsonify({'error': '用户名已存在'}), 400
    
    # 创建新用户
    user_obj = objects.user(username)
    users_db[username] = {
        'password': password,
        'user_obj': user_obj
    }
    
    return jsonify({'message': '注册成功'}), 201

@app.route('/login', methods=['POST'])
def login():
    """用户登录"""
    if not request.is_json:
        return jsonify({'error': '请求格式必须为JSON'}), 400
    
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')
    
    if not username or not password:
        return jsonify({'error': '用户名和密码不能为空'}), 400
    
    user_data = users_db.get(username)
    if not user_data or user_data['password'] != password:
        return jsonify({'error': '用户名或密码错误'}), 401
    
    # 设置会话
    session['username'] = username
    return jsonify({'message': '登录成功'}), 200

@app.route('/message', methods=['POST'])
def add_message():
    """提交留言"""
    if 'username' not in session:
        return jsonify({'error': '请先登录'}), 401
    
    if not request.is_json:
        return jsonify({'error': '请求格式必须为JSON'}), 400
    
    data = request.get_json()
    content = data.get('content')
    
    if not content:
        return jsonify({'error': '留言内容不能为空'}), 400
    
    # 获取当前登录用户
    username = session['username']
    user_obj = users_db[username]['user_obj']
    
    # 创建并添加留言
    new_commit = objects.commit(user_obj, content)
    board.add_commit(new_commit)
    
    return jsonify({'message': '留言成功'}), 201

@app.route('/messages', methods=['GET'])
def get_messages():
    """查看所有留言"""
    # 准备留言数据
    messages = []
    for commit in board.commits:
        messages.append({
            'username': commit.user.name,
            'content': commit.content,
            'timestamp': commit.timestamp.strftime('%Y-%m-%d %H:%M:%S'),
            'id': commit.id
        })
    
    return jsonify({'messages': messages}), 200

# 确保应用可以通过gunicorn正常启动
if __name__ == '__main__':
    app.run(debug=True)


