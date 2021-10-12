import { compose } from '../fn';

describe('compose utility function', () => {
  test('returns a function that works as expected', () => {
    const add = addition => str => `${str}${addition}`;
    const greet = compose(add('Hello'), add(', '), add('World!'));

    expect(greet('')).toBe('Hello, World!');
  });

  test('calls each function with the result of previous one', () => {
    const input = '*';
    const aResult = 'A';
    const bResult = 'B';
    const cResult = 'C';

    const a = jest.fn(() => aResult);
    const b = jest.fn(() => bResult);
    const c = jest.fn(() => cResult);

    const output = compose(a, b, c)(input);

    expect(a).toHaveBeenCalledWith(input);
    expect(b).toHaveBeenCalledWith(aResult);
    expect(c).toHaveBeenCalledWith(bResult);
    expect(output).toBe(cResult);
  });
});
