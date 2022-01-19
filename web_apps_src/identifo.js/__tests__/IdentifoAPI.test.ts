import { INVALID_TOKEN_ERROR } from '../src/constants';
import IdentifoAuth from '../src/IdentifoAuth';
import { setupFetchStub } from './utils/fetch';

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

  test('getToken should return token object', async () => {
    Object.defineProperty(window, 'location', {
      value: { href: config.redirectUri, search: `?token=${generatedToken}` },
      writable: true,
    });
    expect(identifo.handleAuthentication()).resolves.toBe(true);

    let tokenData = await identifo.getToken();
    if (tokenData) {
      expect(tokenData.token).toBe(generatedToken);
    }

    expect(identifo.isAuth).toBe(true);

    global.fetch = jest.fn().mockImplementation(
      setupFetchStub({
        result: 'ok',
      }),
    );

    await expect(identifo.api.logout()).resolves.toStrictEqual({ result: 'ok' });

    expect(identifo.isAuth).toBe(false);
  });
});
