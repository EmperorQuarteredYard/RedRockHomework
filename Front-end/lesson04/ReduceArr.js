function ReduceArr(arr, depth) {
  // console.table("depth", depth);
  // console.table("arr", arr);
  depth++;
  let result = [];
  for (let i = 0; i < arr.length; i++) {

    if (Array.isArray(arr[i])) {
      oarr = ReduceArr(arr[i], depth);
      for (let j of oarr) {
        result.push(j);
      }
    } else {
      result.push(arr[i]);
    }
  }
  // console.table("result", result);
  return result;
}
let dep = 0;
let arr = [[1, [2, 3], 4], [5, [6, 7], 8, 9]];
console.log(ReduceArr(arr, dep));
