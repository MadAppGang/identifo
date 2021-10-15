class CookieStorage {
  isAccessible = false;

  saveToken(): boolean {
    return true;
  }

  getToken(): string {
    throw new Error('Can not get token from HttpOnly');
  }

  deleteToken(): void {
    // throw new Error('Can not get token from HttpOnly');
  }
}

export default CookieStorage;
