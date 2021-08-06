import { pause } from '~/utils';

const data = {
  settings: {
    email: 'john@doe.com',
    loginEnvName: 'IDENTIFO_ADMIN_LOGIN',
    passwordEnvName: 'IDENTIFO_ADMIN_PASSWORD',
  },
};

const createAccountServiceMock = () => {
  const fetchSettings = async () => {
    await pause(1000);

    return data.settings;
  };

  const postSettings = async (settings) => {
    await pause(1000);

    if (settings.email === 'trigger@error.com') {
      throw new Error('Email is already in use');
    }

    data.settings = settings;
  };

  return Object.freeze({
    fetchSettings,
    postSettings,
  });
};

export default createAccountServiceMock;
