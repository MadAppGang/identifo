const createDatabaseService = ({ httpClient }) => {
    const testConnection = async (settings) => {
        const url = `${httpClient.getApiUrl()}/settings/storage/test`;
        const { data } = await httpClient.post(url, settings);

        return data;
    };

    const fetchSettings = async () => {
        const url = `${httpClient.getApiUrl()}/settings`;
        const { data } = await httpClient.get(url);

        return data.storage;
    };

    const postSettings = async (storage) => {
        const url = `${httpClient.getApiUrl()}/settings/storage`;
        const { data } = await httpClient.put(url, storage);

        return data;
    };

    return Object.freeze({
        testConnection,
        fetchSettings,
        postSettings,
    });
};

export default createDatabaseService;
