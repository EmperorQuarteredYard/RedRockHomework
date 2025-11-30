// 使用英文变量名，确保更好的兼容性
function BuildObject_Food(name) {
  this.name = name;
  this.comment = []; // 使用更简洁的数组创建方式
}

function BuildObject_Comment(name, content, user, time, score) {
  this.name = name;
  this.content = content;
  this.user = user;
  this.time = time;
  this.score = score;
}
function CopyComment(comment) {
  for (let key in comment) {
    this[key] = comment[key];
  }
  return this;
}
// 使用英文变量名
let spicyChicken = new BuildObject_Food('辣子鸡');
let comment1 = new BuildObject_Comment('中国移动', '这道菜真好吃', '10086', '2023-01-01', 5);
spicyChicken.comment.push(comment1);

let spicyPork = new BuildObject_Food('辣椒炒肉');
let comment2 = new BuildObject_Comment('南方妹子', '爽', '13311', '2023-01-01', 100);
spicyPork.comment.push(comment1);
spicyPork.comment.push(comment2);

// 添加简单的测试输出，确保控制台正常工作
console.log("JavaScript代码开始执行");
console.log("辣子鸡对象:", spicyChicken);
console.log("辣椒炒肉对象:", spicyPork);

let arr = [];
arr.push(spicyChicken);
arr.push(spicyPork);
console.log("数组内容:", arr);
comment3 = comment1;
comment4 = CopyComment(comment2);
// 确保table输出正常
console.table(arr);