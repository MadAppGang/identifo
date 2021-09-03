import { toDeepCase } from '~/utils/apiMapper';

const createAccountService = ({ httpClient }) => {
  const postSettings = async (settings) => {
    const url = `${httpClient.getApiUrl()}/settings/account`;
    const { data } = httpClient.patch(url, toDeepCase(settings, 'snake'));

    return data;
  };

  return Object.freeze({
    postSettings,
  });
};

export default createAccountService;
