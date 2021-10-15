import { toDeepCase } from '../apiMapper';

const SNAKE = 'snake';
const CAMEL = 'camel';

describe('toDeepCase transformator', () => {
  test('transforms camel into snake properly', () => {
    const input = {
      helloWorld: 'Hello, world!',
    };
    const expectedResult = {
      hello_world: 'Hello, world!',
    };
    expect(toDeepCase(input, SNAKE)).toEqual(expectedResult);
  });

  test('transforms snake into camel properly', () => {
    const input = {
      hello_world: 'Hello, world!',
    };
    const expectedResult = {
      helloWorld: 'Hello, world!',
    };
    expect(toDeepCase(input, CAMEL)).toEqual(expectedResult);
  });

  test('array value stays array after snake to camel transformation', () => {
    const input = {
      hello_world: ['Hello', ',', 'World!'],
      foo_bar: [],
    };
    const expectedResult = {
      helloWorld: ['Hello', ',', 'World!'],
      fooBar: [],
    };
    expect(toDeepCase(input, CAMEL)).toEqual(expectedResult);
  });

  test('array value stays array after camel to snake transformation', () => {
    const input = {
      helloWorld: ['Hello', ',', 'World!'],
      fooBar: [],
    };
    const expectedResult = {
      hello_world: ['Hello', ',', 'World!'],
      foo_bar: [],
    };
    expect(toDeepCase(input, SNAKE)).toEqual(expectedResult);
  });
});
