import StorageManager from './storage-manager';
declare class SessionStorage extends StorageManager {
    constructor(accessKey?: string, refreshKey?: string);
}
export default SessionStorage;
