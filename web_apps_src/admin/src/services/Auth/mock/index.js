import { pause } from '~/utils';

const AUTH_STORAGE_KEY = 'auth_storage_key';

const createAuthServiceMock = () => {
  const login = async (email, password) => {
    await pause(1000);

    if (email === 'email' && password === 'password') {
      localStorage.setItem(AUTH_STORAGE_KEY, 'helloworld');
      return;
    }

    throw new Error('Email or password is incorrect');
  };

  const logout = async () => {
    localStorage.removeItem(AUTH_STORAGE_KEY);
  };

  const checkAuthState = async () => {
    await pause(500);

    return !!localStorage.getItem(AUTH_STORAGE_KEY);
  };

  return Object.freeze({
    login,
    logout,
    checkAuthState,
  });
};

export default createAuthServiceMock;
