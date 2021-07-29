import { TokenManager, TokenType } from '../types/types';

class StorageManager implements TokenManager {
  preffix = 'identifo_';

  storageType: 'localStorage' | 'sessionStorage' = 'localStorage';

  access = `${this.preffix}access_token`;

  refresh = `${this.preffix}refresh_token`;

  isAccessible = true;

  constructor(storageType: 'localStorage' | 'sessionStorage', accessKey?: string, refreshKey?: string) {
    this.access = accessKey ? this.preffix + accessKey : this.access;
    this.refresh = refreshKey ? this.preffix + refreshKey : this.refresh;
    this.storageType = storageType;
  }

  saveToken(token: string, tokenType: TokenType): boolean {
    if (token) {
      window[this.storageType].setItem(this[tokenType], token);
      return true;
    }
    return false;
  }

  getToken(tokenType: TokenType): string {
    return window[this.storageType].getItem(this[tokenType]) ?? '';
  }

  deleteToken(tokenType: TokenType): void {
    window[this.storageType].removeItem(this[tokenType]);
  }
}

export default StorageManager;
