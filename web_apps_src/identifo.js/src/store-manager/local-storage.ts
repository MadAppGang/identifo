import StorageManager from './storage-manager';

class LocalStorage extends StorageManager {
  constructor(accessKey?: string, refreshKey?: string) {
    super('localStorage', accessKey, refreshKey);
  }
}

export default LocalStorage;
