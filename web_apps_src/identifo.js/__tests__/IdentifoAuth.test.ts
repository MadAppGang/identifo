import { INVALID_TOKEN_ERROR } from '../src/constants';
import IdentifoAuth from '../src/IdentifoAuth';

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
    exp: new Date().getTime() / 1000 + 7200,
  };
  const generatedToken = jwt.encode(payload, 'secret');
  const identifo = new IdentifoAuth(config);

  const falsyPayload = {
    type: 'access',
    aud: config.appId,
    iss: config.issuer,
    exp: new Date().getTime() / 1000 - 7200,
  };
  const falsyGeneratedToken = jwt.encode(falsyPayload, 'secret');

  test('should be defined', async () => {
    expect(identifo).toBeDefined();
    expect(identifo.isAuth).toBe(false);
    expect(identifo.getToken()).resolves.toBeNull();
  });

  test('handleAuthentication should return false when has no token', () => {
    expect(identifo.handleAuthentication()).rejects.toBe(undefined);
  });

  test('handleAuthentication should return true if has correct url and token', async () => {
    Object.defineProperty(window, 'location', {
      value: { href: config.redirectUri, search: `?token=${generatedToken}` },
      writable: true,
    });
    expect(identifo.handleAuthentication()).resolves.toBe(true);
  });

  test('getToken should return token object', async () => {
    const tokenData = await identifo.getToken();
    if (tokenData) {
      expect(Object.keys(tokenData)).toEqual(['token', 'payload']);
      expect(typeof tokenData.token === 'string').toBe(true);
      expect(tokenData.token).toBe(generatedToken);
    }
  });

  describe('Falsy scenario:', () => {
    test('handleAuthentication should be falsy', async () => {
      Object.defineProperty(window, 'location', {
        value: { href: config.redirectUri, search: `?token=${falsyGeneratedToken}` },
        writable: true,
      });
      expect(identifo.handleAuthentication()).rejects.toBe(undefined);
    });

    test('getToken should return null if token is invalid', async () => {
      expect(identifo.getToken()).resolves.toBeNull();
    });
  });
});
