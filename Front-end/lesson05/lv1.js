function factorial(n) {
  if (n === 1) return 1;
  return n * factorial(n - 1);
}
function factorialIterative(n) {
  let result = 1;
  for (let i = 2; i <= n; i++) {
    result *= i;
  }
  return result;
}