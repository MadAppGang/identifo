import { IdentifoConfig, UrlFlows } from './types/types';

export class UrlBuilder {
  constructor(private config: IdentifoConfig) {}

  getUrl(flow: UrlFlows): string {
    const scopes = this.config.scopes?.join() || '';
    const redirectUri = this.config.redirectUri ?? window.location.href;
    const baseParams = `appId=${this.config.appId}&scopes=${scopes}`;
    const urlParams = `${baseParams}&callbackUrl=${encodeURIComponent(redirectUri)}`;
    // if postLogoutRedirectUri is empty, login url will be instead
    const postLogoutRedirectUri = this.config.postLogoutRedirectUri
      ? `${this.config.postLogoutRedirectUri}`
      : `${redirectUri}&redirectUri=${this.config.url}/web/login?${encodeURIComponent(baseParams)}`;

    const urls = {
      signup: `${this.config.url}/web/register?${urlParams}`,
      signin: `${this.config.url}/web/login?${urlParams}`,
      logout: `${this.config.url}/web/logout?${baseParams}&callbackUrl=${encodeURIComponent(postLogoutRedirectUri)}`,
      renew: `${this.config.url}/web/token/renew?${baseParams}&redirectUri=${encodeURIComponent(redirectUri)}`,
      default: 'default',
    };

    return urls[flow] || urls.default;
  }

  createSignupUrl(): string {
    return this.getUrl('signup');
  }

  createSigninUrl(): string {
    return this.getUrl('signin');
  }

  createLogoutUrl(): string {
    return this.getUrl('logout');
  }

  createRenewSessionUrl(): string {
    return this.getUrl('renew');
  }
}
