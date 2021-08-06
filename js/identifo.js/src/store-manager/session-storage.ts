import StorageManager from './storage-manager';

class SessionStorage extends StorageManager {
  constructor(accessKey?: string, refreshKey?: string) {
    super('sessionStorage', accessKey, refreshKey);
  }
}

export default SessionStorage;
