import { toDeepCase } from '~/utils/apiMapper';

const CAMEL_CASE = 'camel';
const SNAKE_CASE = 'snake';

const createSettingsService = ({ httpClient }) => {
  const fetchServerSettings = async () => {
    const url = `${httpClient.getApiUrl()}/settings`;
    const { data } = await httpClient.get(url);
    return toDeepCase(data, CAMEL_CASE);
  };

  const updateLoginSettings = (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/login`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const fetchExternalServicesSettings = async () => {
    const url = `${httpClient.getApiUrl()}/settings/services`;
    const { data } = await httpClient.get(url);
    return toDeepCase(data, CAMEL_CASE);
  };

  const updateExternalServicesSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/services`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const updateSessionStorageSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/storage/session`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const updateStaticFilesSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/static`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const updateGeneralSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/general`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const updateConfigurationStorageSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/storage/configuration`;
    return httpClient.put(url, toDeepCase(settings, SNAKE_CASE));
  };

  const uploadJWTKeys = async (pubKey, privKey) => {
    const url = `${httpClient.getApiUrl()}/static/uploads/keys`;

    const formData = new FormData();
    formData.append('keys', pubKey, 'public.pem');
    formData.append('keys', privKey, 'private.pem');

    return httpClient.post(url, formData);
  };

  const requestServerRestart = async () => {
    const url = `${httpClient.getApiUrl()}/restart`; // TODO: not final
    await httpClient.post(url);
  };

  return {
    updateLoginSettings,
    fetchExternalServicesSettings,
    updateExternalServicesSettings,
    updateSessionStorageSettings,
    updateStaticFilesSettings,
    updateGeneralSettings,
    updateConfigurationStorageSettings,
    uploadJWTKeys,
    requestServerRestart,
    fetchServerSettings,
  };
};

export default createSettingsService;
