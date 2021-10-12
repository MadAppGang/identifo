import { pause } from '~/utils';

const staticFiles = {
  invite: {
    html: '<div>invitation email markup</div>',
    css: '.invite { display: flex; }',
    js: 'function (a, b) { return a + b; }',
  },
  welcome: {
    html: '<div>welcome email markup</div>',
    css: '.invite { display: flex; }',
    js: 'function (a, b) { return a + b; }',
  },
  reset: {
    html: '<div>reset password markup</div>',
    css: '.invite { display: flex; }',
    js: 'function (a, b) { return a + b; }',
  },
  verify: {
    html: '<div>verify email markup</div>',
    css: '.invite { display: flex; }',
    js: 'function (a, b) { return a + b; }',
  },
  tfa: {
    html: '<div>2fa email markup</div>',
    css: '.invite { display: flex; }',
    js: 'function (a, b) { return a + b; }',
  },
};

const createStaticServiceMock = () => {
  const fetchStaticFile = async (name, ext) => {
    await pause(400);

    return staticFiles[name][ext];
  };

  const updateStaticFile = async (name, ext, contents) => {
    await pause(400);

    staticFiles[name][ext] = contents;
  };

  return {
    fetchStaticFile,
    updateStaticFile,
  };
};

export default createStaticServiceMock;
