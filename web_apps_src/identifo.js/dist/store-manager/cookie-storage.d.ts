declare class CookieStorage {
    isAccessible: boolean;
    saveToken(): boolean;
    getToken(): string;
    deleteToken(): void;
}
export default CookieStorage;
