const createDatabaseService = ({ httpClient }) => {
  const testConnection = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/storage/test`;
    const { data } = await httpClient.post(url, settings);

    return data;
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
    postSettings,
    verifySettings,
  });
};

export default createDatabaseService;
