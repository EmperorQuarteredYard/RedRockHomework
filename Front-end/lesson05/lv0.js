const person = {
  name: '小明',
  age: 18,
  address: {
    city: 'ChongQing',
    area: 'RedRock'
  },
  title: ['student', { year: 2025, title: 'GoodStudent' }]
}
// 你的代码
const { name, age, address: { city, area: mountain }, title: [title1, { year, title: title2 }] } = person;

console.log(name) // 小明
console.log(year) // 18
console.log(city) // ChongQing
console.log(mountain) // RedRock
console.log(title1) // student
console.log(title2) // GoodStudent