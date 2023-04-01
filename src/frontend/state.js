export const useState = (value = undefined) => {
  return [() => value, (x) => {
    value = x;
  }];
};
