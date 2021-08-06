const { resolve } = require('path');

module.exports = {
  rootDir: resolve(__dirname, '../../src'),
  snapshotSerializers: ['enzyme-to-json/serializer'],
  setupFiles: [resolve(__dirname, 'setup.js')],
  moduleNameMapper: {
    '~(.*)$': '<rootDir>/$1',
    '\\.(css|sass|scss)$': 'identity-obj-proxy',
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$': 'identity-obj-proxy',
  },
};
