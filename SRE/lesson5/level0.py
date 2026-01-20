class book:
    def __init__(self, name, author, price):
        self.name = name
        self.author = author
        self.price = price
        self.borrowed = False
    def borrow(self):
        if self.borrowed:
            print("这本书已经被借走了")
        else:
            self.borrowed = True
            print("你成功借走了这本书")
    def return_book(self):
        if self.borrowed:
            self.borrowed = False
            print("你成功归还了这本书")
        else:
            print("这本书没有被借走")

class user:
    def __init__(self, name, age):
        self.name = name
        self.age = age
        self.borrowed_books = []
    def borrow_book(self, book):
        if book.borrowed:
            print(f"{book.name}已经被借走了")
        else:
            self.borrowed_books.append(book.name)
            book.borrow()
    def return_book(self, book):
        if book in self.borrowed_books:
            self.borrowed_books.remove(book)
            book.return_book()
        else:
            print(f"{self.name}没有借走这本书{book.name}")
    def show_borrowed_books(self):
        if self.borrowed_books:
            print(f"{self.name}借走的书籍有：")
            for book in self.borrowed_books:
                print(book)
        else:
            print(f"{self.name}没有借走任何书籍")

class library:
    def __init__(self):
        self.books = []
    def add_book(self, book):
        self.books.append(book)
        print(f"成功添加{book.name}")
    def remove_book(self, book):
        if book in self.books:
            self.books.remove(book)
            print(f"成功删除{book.name}")
        else:
            print("这本书不在图书馆中")
    def display_books(self):
        if self.books:
            print("图书馆中的书籍有：")
            for book in self.books:
                print(f"{book.name} by {book.author},借出状态：{book.borrowed}")
        else:
            print("图书馆中没有书籍")
    def search_book(self, name):
        for book in self.books:
            if name.lower() in book.name.lower():
                return book

        print("没有符合条件的书籍")
        return None

lib=library()
lib.add_book(book("python基础","张三",100))
lib.add_book(book("python高级","李四",200))
lib.add_book(book("python网络编程","王五",300))
lib.display_books()

user1=user("张三",18)

user1.borrow_book(lib.books[0])
user1.show_borrowed_books()
lib.display_books()

user1.borrow_book(lib.books[0])
user1.show_borrowed_books()
lib.display_books()

user1.return_book(lib.books[0])
user1.show_borrowed_books()
lib.display_books()
