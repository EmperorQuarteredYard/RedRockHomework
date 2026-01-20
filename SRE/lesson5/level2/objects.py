import datetime

class user:
    id_counter = 0
    def __init__(self, name):
        self.name = name
        self.id = user.id_counter
        id_counter += 1

class commit:
    id_counter = 0
    def __init__(self, user, content):
        self.user = user.id
        self.content = content
        self.timestamp = datetime.datetime.now()
        self.id = id_counter
        id_counter += 1
        self.timestamp = datetime.now()
        

class commitboard:
    def __init__(self):
        self.commits = []
    def add_commit(self, commit):
        self.commits.append(commit)
        print(f"成功添加{commit.user.name}的评论")
    def display_commits(self):
        if self.commits:
            print("评论板中的评论有：")
            for commit in self.commits:
                print(f"{commit.user.name}说：{commit.content}")
        else:
            print("评论板中没有评论")
