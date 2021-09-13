import { toDeepCase } from '~/utils/apiMapper';

const CAMEL_CASE = 'camel';
const SNAKE_CASE = 'snake';

const createSettingsService = ({ httpClient }) => {
  const fetchServerSettings = async () => {
    const url = `${httpClient.getApiUrl()}/settings`;
    const { data } = await httpClient.get(url);
    return toDeepCase(data, CAMEL_CASE);
  };

  const postServerSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings`;
    const { data } = await httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
    return toDeepCase(data, CAMEL_CASE);
  };

  const getJWTKeys = async (withPrivate) => {
    const url = `${httpClient.getApiUrl()}/static/keys?include_private_key=${withPrivate}`;

    const { data } = await httpClient.get(url);
    return data;
  };

  const uploadJWTKeys = async (privKey) => {
    const url = `${httpClient.getApiUrl()}/static/uploads/keys`;
    const { data } = await httpClient.post(url, privKey);
    return toDeepCase(data, CAMEL_CASE);
  };

  const generateKeys = async (alg) => {
    const url = `${httpClient.getApiUrl()}/generate_new_secret`;
    const { data } = await httpClient.post(url, alg);
    return toDeepCase(data, CAMEL_CASE);
  };

  const requestServerRestart = async () => {
    const url = `${httpClient.getApiUrl()}/restart`; // TODO: not final
    await httpClient.post(url);
  };

  const verifySettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/test_connection`;
    const { data } = await httpClient.post(url, toDeepCase(settings, 'snake'));
    return data;
  };

  return {
    uploadJWTKeys,
    postServerSettings,
    requestServerRestart,
    fetchServerSettings,
    getJWTKeys,
    generateKeys,
    verifySettings,
  };
};

export default createSettingsService;
