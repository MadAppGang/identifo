import { INVALID_TOKEN_ERROR } from './constants';
import { LocalStorageManager } from './store-manager';
import {
  ClientToken, JWTPayload, TokenManager, TokenType,
} from './types/types';

class TokenService {
  private tokenManager: TokenManager;

  constructor(tokenManager?: TokenManager) {
    this.tokenManager = tokenManager || new LocalStorageManager();
    // TODO: implement cookie as default
    // this.tokenManager = tokenManager || new CoockieStorage();
  }

  async handleVerification(token: string, audience: string, issuer?: string): Promise<boolean> {
    if (!this.tokenManager.isAccessible) return true;
    try {
      await this.validateToken(token, audience, issuer);
      this.saveToken(token);
      return true;
    } catch (err) {
      this.removeToken();
      return Promise.reject(err);
    }
  }

  async validateToken(token: string, audience: string, issuer?: string): Promise<boolean> {
    if (!token) throw new Error(INVALID_TOKEN_ERROR);
    const jwtPayload = this.parseJWT(token);
    const isJwtExpired = this.isJWTExpired(jwtPayload);
    if (jwtPayload.aud?.includes(audience) && (!issuer || jwtPayload.iss === issuer) && !isJwtExpired) {
      return Promise.resolve(true);
    }
    throw new Error(INVALID_TOKEN_ERROR);
  }

  parseJWT(token: string): JWTPayload {
    const base64Url = token.split('.')[1];
    if (!base64Url) return { aud: [], iss: '', exp: 10 };
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`)
        .join(''),
    );
    return JSON.parse(jsonPayload) as JWTPayload;
  }

  isJWTExpired(token: JWTPayload): boolean {
    const now = new Date().getTime() / 1000;
    if (token.exp && now > token.exp) {
      return true;
    }
    return false;
  }

  isAuthenticated(audience: string, issuer?: string): Promise<boolean> {
    if (!this.tokenManager.isAccessible) return Promise.resolve(true);
    const token = this.tokenManager.getToken('access');
    // TODO: may be change to handleAuth instead validateToken
    return this.validateToken(token, audience, issuer);
  }

  saveToken(token: string, type: TokenType = 'access'): boolean {
    return this.tokenManager.saveToken(token, type);
  }

  removeToken(type: TokenType = 'access'): void {
    this.tokenManager.deleteToken(type);
  }

  getToken(type: TokenType = 'access'): ClientToken | null {
    const token = this.tokenManager.getToken(type);
    if (!token) return null;
    const jwtPayload = this.parseJWT(token);
    return { token, payload: jwtPayload };
  }
}

export default TokenService;
