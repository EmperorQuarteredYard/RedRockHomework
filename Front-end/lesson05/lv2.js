
var name = 'window'
function Person(name) {
  this.name = name
  this.foo1 = function () {
    console.log(this.name)
  }//输出名字:花式定义1
  this.foo2 = () => console.log(this.name)//输出名字:箭头函数2
  this.foo3 = function () {
    return function () {
      console.log(this.name)
    }
  }//输出名字:花式定义3
  this.foo4 = function () {
    return () => {
      console.log(this.name)
    }
  }
}//输出名字:花式定义4

var person1 = new Person('person1')
var person2 = new Person('person2')

person1.foo1()
person1.foo1.call(person2) //person2借person1尸还魂，输出结果:person2

person1.foo2()
person1.foo2.call(person2)

person1.foo3()()
person1.foo3.call(person2)()
person1.foo3().call(person2)

person1.foo4()()
person1.foo4.call(person2)()
person1.foo4().call(person2) 