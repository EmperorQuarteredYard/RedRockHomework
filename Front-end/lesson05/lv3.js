export function IsObject(value) {
  if (typeof value === 'object' && value !== null) return true;
  return false;
}
function DeepCopy(orginalValue, map = new WeakMap()) {
  if (!IsObject(orginalValue)) return orginalValue;
  // 处理循环引用
  if (map.has(orginalValue)) return map.get(orginalValue);
  // 处理Map
  if (orginalValue instanceof Map) {
    const newmap = new Map();
    for (const item of originalValue) {
      newmap.set(item[0], DeepCopy(item[1], map));
    }
    return newmap;
  }
  //处理Set
  if (orginalValue instanceof Set) {
    const newset = new Set();
    for (const item of originalValue) {
      newset.add(DeepCopy(item, map));
    }
    return newset;
  }
  //处理普通对象
  if (IsObject(orginalValue)) {
    const newObj = {};
    for (const key in orginalValue) {
      newObj[key] = DeepCopy(orginalValue[key], map)
    }
    const symbolKeys = Object.getOwnPropertySymbols(orginalValue)
    for (const symbolKey of symbolKeys) {
      newObj[Symbol(symbolKey.description)] = DeepCopy(orginalValue[symbolKey], map)
    }
    return newObj;
  }
}