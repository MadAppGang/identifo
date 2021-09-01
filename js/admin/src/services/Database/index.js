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

  const verifySettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/test_connection`;
    const { data } = await httpClient.put(url, settings);
    return data;
  };

  return Object.freeze({
    testConnection,
    fetchSettings,
    postSettings,
    verifySettings,
  });
};

export default createDatabaseService;
