const createDatabaseService = ({ httpClient }) => {
  const testConnection = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/storage/test`;
    const { data } = await httpClient.post(url, settings);

    return data;
  };

  const verifySettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/test_connection`;
    const { data } = await httpClient.put(url, settings);
    return data;
  };

  return Object.freeze({
    testConnection,
    verifySettings,
  });
};

export default createDatabaseService;
