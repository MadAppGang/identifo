import { Api } from './api/api';
import { jwtRegex, REFRESH_TOKEN_QUERY_KEY, TOKEN_QUERY_KEY } from './constants';
import TokenService from './tokenService';
import { ClientToken, IdentifoConfig, UrlBuilderInit } from './types/types';
import { UrlBuilder } from './UrlBuilder';

class IdentifoAuth {
  public api: Api;

  public tokenService: TokenService;

  public config: IdentifoConfig;

  public urlBuilder: UrlBuilderInit;

  private token: ClientToken | null = null;

  isAuth = false;

  constructor(config: IdentifoConfig) {
    this.config = { ...config, autoRenew: config.autoRenew ?? true };
    this.tokenService = new TokenService(config.tokenManager);
    this.urlBuilder = new UrlBuilder(this.config);
    this.api = new Api(config, this.tokenService);
    this.handleToken(this.tokenService.getToken()?.token || '', 'access');
  }

  private handleToken(token: string, tokenType: 'access' | 'refresh') {
    if (token) {
      if (tokenType === 'access') {
        const payload = this.tokenService.parseJWT(token);
        this.token = { token, payload };
        this.isAuth = true;
        this.tokenService.saveToken(token);
      } else {
        this.tokenService.saveToken(token, 'refresh');
      }
    }
  }

  private resetAuthValues() {
    this.token = null;
    this.isAuth = false;
    this.tokenService.removeToken();
    this.tokenService.removeToken('refresh');
  }

  signup(): void {
    window.location.href = this.urlBuilder.createSignupUrl();
  }

  signin(): void {
    window.location.href = this.urlBuilder.createSigninUrl();
  }

  logout(): void {
    this.resetAuthValues();
    window.location.href = this.urlBuilder.createLogoutUrl();
  }

  async handleAuthentication(): Promise<boolean> {
    const { access, refresh } = this.getTokenFromUrl();
    if (!access) {
      this.resetAuthValues();
      return Promise.reject();
    }
    try {
      await this.tokenService.handleVerification(access, this.config.appId, this.config.issuer);
      this.handleToken(access, 'access');
      if (refresh) {
        this.handleToken(refresh, 'refresh');
      }
      return await Promise.resolve(true);
    } catch (err) {
      this.resetAuthValues();
      return await Promise.reject();
    } finally {
      // TODO: Nikita K cahnge correct window key
      window.location.hash = '';
    }
  }

  private getTokenFromUrl(): { access: string; refresh: string } {
    const urlParams = new URLSearchParams(window.location.search);
    const tokens = { access: '', refresh: '' };
    const accessToken = urlParams.get(TOKEN_QUERY_KEY);
    const refreshToken = urlParams.get(REFRESH_TOKEN_QUERY_KEY);
    if (refreshToken && jwtRegex.test(refreshToken)) {
      tokens.refresh = refreshToken;
    }
    if (accessToken && jwtRegex.test(accessToken)) {
      tokens.access = accessToken;
    }
    return tokens;
  }

  async getToken(): Promise<ClientToken | null> {
    const token = this.tokenService.getToken();
    const refreshToken = this.tokenService.getToken('refresh');
    if (token) {
      const isExpired = this.tokenService.isJWTExpired(token.payload);
      if (isExpired && refreshToken) {
        try {
          await this.renewSession();
          return await Promise.resolve(this.token);
        } catch (err) {
          this.resetAuthValues();
          throw new Error('No token');
        }
      }
      return Promise.resolve(token);
    }
    return Promise.resolve(null);
  }

  async renewSession(): Promise<string> {
    try {
      const { access, refresh } = await this.renewSessionWithToken();
      this.handleToken(access, 'access');
      this.handleToken(refresh, 'refresh');
      return await Promise.resolve(access);
    } catch (err) {
      return Promise.reject();
    }
  }

  private async renewSessionWithToken(): Promise<{ access: string, refresh: string }> {
    try {
      const tokens = await this.api.renewToken()
        .then((l) => ({ access: l.access_token || '', refresh: l.refresh_token || '' }));
      return tokens;
    } catch (err) {
      return Promise.reject(err);
    }
  }
}
export default IdentifoAuth;
