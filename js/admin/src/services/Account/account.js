import { toDeepCase } from '~/utils/apiMapper';

const createAccountService = ({ httpClient }) => {
  const fetchSettings = async () => {
    const url = `${httpClient.getApiUrl()}/settings`;
    const { data } = await httpClient.get(url);
    return data.admin_account;
  };

  const postSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/account`;
    const { data } = httpClient.patch(url, toDeepCase(settings, 'snake'));

    return data;
  };

  return Object.freeze({
    fetchSettings,
    postSettings,
  });
};

export default createAccountService;
