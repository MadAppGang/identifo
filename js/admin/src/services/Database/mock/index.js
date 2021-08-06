import { pause } from '~/utils';

const data = {
  storage: {
    app_storage: {
      type: 'boltdb',
      name: 'apps',
      endpoint: 'localhost:27017',
      region: 'us-west-2',
      path: './db.db',
    },
    user_storage: {
      type: 'boltdb',
      name: 'users',
      endpoint: 'localhost:27017',
      region: 'us-west-2',
      path: './db.db',
    },
    token_storage: {
      type: 'boltdb',
      name: 'tokens',
      endpoint: 'localhost:27017',
      region: 'us-west-2',
      path: './db.db',
    },
    verification_code_storage: {
      type: 'boltdb',
      name: 'verification_codes',
      endpoint: 'localhost:27017',
      region: 'us-west-2',
      path: './db.db',
    },
    token_blacklist: {
      type: 'boltdb',
      name: 'blacklist',
      endpoint: 'localhost:27017',
      region: 'us-west-2',
      path: './db.db',
    },
  },
};

const createDatabaseServiceMock = () => {
  const testConnection = async (dbSettings) => {
    await pause(1000);

    if (dbSettings.type === 'mongodb') {
      return { result: 'ok' };
    }

    throw new Error('Unable to connect to specified endpoint');
  };

  const fetchSettings = async () => {
    await pause(1000);

    return data.storage;
  };

  const postSettings = async (settings) => {
    await pause(1000);

    data.storage = settings;
  };

  return Object.freeze({
    testConnection,
    fetchSettings,
    postSettings,
  });
};

export default createDatabaseServiceMock;
