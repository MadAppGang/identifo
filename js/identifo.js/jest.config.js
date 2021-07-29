module.exports = {
  preset: 'ts-jest',
  transform: {
    '^.+\\.(t|j)sx?$': 'ts-jest',
  },
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  moduleNameMapper: {
    '^jose/(.*)$': '<rootDir>/node_modules/jose/dist/node/cjs/$1',
  },
  testPathIgnorePatterns: ['<rootDir>/dist/', '<rootDir>/node_modules/', '<rootDir>/.history/'],
};
