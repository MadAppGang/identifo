import { UrlBuilder } from '../src/UrlBuilder';

describe('UrlBuilder: ', () => {
  const config = {
    url: 'http://localhost:8081',
    appId: '59fd884d8f6b180001f5b4e2',
    scopes: [],
    postLogoutRedirectUri: 'http://localhost:8081/returnTo',
    redirectUri: 'http://localhost:8081/callbackUrl',
  };
  const urlBuilder = new UrlBuilder(config);
  test('should be defined and has methods', () => {
    expect(urlBuilder).toBeDefined();
    expect(Object.keys(urlBuilder)).toEqual([
      'createSignupUrl',
      'createSigninUrl',
      'createLogoutUrl',
      'createRenewSessionUrl',
    ]);
  });

  test('should return correct url (all params is defined)', () => {
    const baseParams = `appId=${config.appId}&scopes=${JSON.stringify(config.scopes)}`;
    const baseSuffixParam = `${baseParams}&callbackUrl=${config.redirectUri}`;

    expect(urlBuilder.createSignupUrl()).toBe(`${config.url}/web/register?${baseSuffixParam}`);
    expect(urlBuilder.createSigninUrl()).toBe(`${config.url}/web/login?${baseSuffixParam}`);
    expect(urlBuilder.createLogoutUrl()).toBe(
      `${config.url}/web/logout?${baseParams}&callbackUrl=${config.postLogoutRedirectUri}`,
    );
    expect(urlBuilder.createRenewSessionUrl()).toBe(
      `${config.url}/web/token/renew?${baseParams}&redirectUri=${config.redirectUri}`,
    );
  });

  test('should return correct url (only app id & authUrl are defined)', () => {
    const urlConfig = {
      url: 'http://localhost:8081',
      appId: '59fd884d8f6b180001f5b4e2',
    };
    const urls = new UrlBuilder(urlConfig);
    const baseParams = `appId=${urlConfig.appId}&scopes=${JSON.stringify([])}`;
    const baseSuffixParam = `${baseParams}&callbackUrl=${window.location.href}`;

    expect(urls.createSignupUrl()).toBe(`${urlConfig.url}/web/register?${baseSuffixParam}`);
    expect(urls.createSigninUrl()).toBe(`${urlConfig.url}/web/login?${baseSuffixParam}`);
    expect(urls.createLogoutUrl()).toBe(`${urlConfig.url}/web/logout?${baseParams}&callbackUrl=`);
    expect(urls.createRenewSessionUrl()).toBe(
      `${urlConfig.url}/web/token/renew?${baseParams}&redirectUri=${window.location.href}`,
    );
  });
});
