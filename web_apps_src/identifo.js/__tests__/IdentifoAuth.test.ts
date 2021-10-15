import { INVALID_TOKEN_ERROR } from '../src/constants';
import IdentifoAuth from '../src/IdentifoAuth';
import Iframe from '../src/iframe';

const jwt = require('jwt-simple');

describe('IdentifoAuth: ', () => {
  const config = {
    url: 'http://localhost:8081',
    appId: '59fd884d8f6b180001f5b4e2',
    scopes: [],
    issuer: 'http://localhost:8081',
    redirectUri: 'http://localhost:3000/callback',
  };

  const payload = {
    type: 'access',
    aud: config.appId,
    iss: config.issuer,
    exp: new Date().getTime() / 1000 + 3600,
  };
  const generatedToken = jwt.encode(payload, 'secret');
  const hash = `#${generatedToken}`;
  const identifo = new IdentifoAuth(config);

  test('should be defined', () => {
    expect(identifo).toBeDefined();
    expect(identifo.isAuth).toBe(false);
    expect(identifo.getToken()).toBeNull();
  });

  test('handleAuthentication should return false when has no token', () => {
    expect(identifo.handleAuthentication()).resolves.toBe(false);
  });

  test('handleAuthentication should return true if has correct url and token', async () => {
    Object.defineProperty(window, 'location', {
      value: { href: config.redirectUri, hash },
    });
    const handleAuthenticationStatus = await identifo.handleAuthentication();
    expect(handleAuthenticationStatus).toBe(true);
  });

  test('getToken should return token object', () => {
    const tokenData = identifo.getToken();
    if (tokenData) {
      expect(Object.keys(tokenData)).toEqual(['token', 'payload']);
      expect(typeof tokenData.token === 'string').toBe(true);
      expect(tokenData.token).toBe(generatedToken);
    }
  });

  test('renewSession should return generated token', async () => {
    jest.spyOn(Iframe, 'captureMessage')
      .mockImplementation(() => Promise.resolve(generatedToken));
    const createMock = jest.spyOn(Iframe, 'create');
    expect(await identifo.renewSession()).toBe(generatedToken);
    expect(createMock).toBeCalledTimes(1);
    createMock.mockClear();
  });

  describe('Falsy scenario:', () => {
    const falsyPayload = { ...payload, exp: payload.exp - 7200 }; // EXP = current time - 1 hour
    const falsyGeneratedToken = jwt.encode(falsyPayload, 'secret');
    const falsyHash = `#${falsyGeneratedToken}`;
    beforeAll(() => {
      Object.defineProperty(window, 'location', {
        value: { href: config.redirectUri, hash: 'falsyHash' },
      });
      window.localStorage.removeItem('identifo_access_token');
    });
    test('handleAuthentication should be falsy', async () => {
      const status = await identifo.handleAuthentication();
      expect(status).toBe(false);
    });

    test('getToken should return null if token is invalid', () => {
      const tokenData = identifo.getToken();
      expect(tokenData).toBeNull();
    });

    test('renewSession should be rejected', () => {
      const createMock = jest.spyOn(Iframe, 'create');
      jest.spyOn(Iframe, 'captureMessage')
        .mockImplementation(() => Promise.resolve('token'));
      expect(identifo.renewSession()).rejects.toStrictEqual(new Error(INVALID_TOKEN_ERROR));

      jest.spyOn(Iframe, 'captureMessage')
        .mockImplementation(() => Promise.reject(new Error('Some Error')));
      expect(identifo.renewSession()).rejects.toBeInstanceOf(Error);
      expect(createMock).toBeCalledTimes(2);
    });
  });
});
