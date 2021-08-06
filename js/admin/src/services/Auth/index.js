import { pause, getError } from '~/utils';

const createAuthService = ({ httpClient }) => {
  const login = async (email, password) => {
    const url = `${httpClient.getApiUrl()}/login`;

    try {
      const response = await httpClient.post(url, { email, password });

      return response.data;
    } catch (err) {
      throw getError(err);
    }
  };

  const logout = () => {
    const url = `${httpClient.getApiUrl()}/logout`;
    return httpClient.post(url);
  };

  const checkAuthState = () => {
    const url = `${httpClient.getApiUrl()}/me`;

    return new Promise((resolve) => {
      httpClient.get(url)
        .then(() => pause(500))
        .then(() => resolve(true))
        .catch(() => resolve(false));
    });
  };

  return Object.freeze({
    login, logout, checkAuthState,
  });
};

export default createAuthService;
