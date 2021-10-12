import { curry } from '../fn';

const add = (a, b, c) => a + b + c;

describe('Curry utility function', () => {
  test('returns a function', () => {
    expect(curry(add)).toBeInstanceOf(Function);
  });

  test('returns the proper result if called with original number of args', () => {
    expect(curry(add)(1, 2, 3)).toBe(6);
  });

  test('returns a func when args count is less than expected', () => {
    expect(curry(add)(1, 2)).toBeInstanceOf(Function);
  });

  test('returns the result whenever the total number of args is greater than or equal to the original number of args', () => {
    const curriedAdd = curry(add);

    expect(curriedAdd(1)(2)).toBeInstanceOf(Function);
    expect(curriedAdd(1)(2)(3)).toBe(6);
    expect(curriedAdd(1, 2)(3)).toBe(6);
    expect(curriedAdd(1)(2, 3)).toBe(6);
    expect(curriedAdd(1, 2)(3, 4, 5, 6, 7)).toBe(6);
  });
});
