const createDatabaseService = ({ httpClient }) => {
  const testConnection = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/storage/test`;
    const { data } = await httpClient.post(url, settings);

    return data;
  };

  return Object.freeze({
    testConnection,
  });
};

export default createDatabaseService;
