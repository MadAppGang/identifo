import { TokenManager, TokenType } from '../types/types';
declare class StorageManager implements TokenManager {
    preffix: string;
    storageType: 'localStorage' | 'sessionStorage';
    access: string;
    refresh: string;
    isAccessible: boolean;
    constructor(storageType: 'localStorage' | 'sessionStorage', accessKey?: string, refreshKey?: string);
    saveToken(token: string, tokenType: TokenType): boolean;
    getToken(tokenType: TokenType): string;
    deleteToken(tokenType: TokenType): void;
}
export default StorageManager;
