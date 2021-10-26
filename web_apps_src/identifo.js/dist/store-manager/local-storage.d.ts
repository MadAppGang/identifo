import StorageManager from './storage-manager';
declare class LocalStorage extends StorageManager {
    constructor(accessKey?: string, refreshKey?: string);
}
export default LocalStorage;
