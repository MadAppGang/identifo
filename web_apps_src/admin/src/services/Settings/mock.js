import { pause } from '~/utils';
import { toDeepCase } from '~/utils/apiMapper';

const data = {
  general: {
    host: 'http://localhost:8081',
    issuer: 'http://localhost:8081',
    algorithm: 'auto',
  },
  login: {
    loginWith: {
      username: true,
      phone: false,
      federated: false,
    },
    tfaType: 'app',
  },
  externalServices: {
    emailService: {
      type: 'mailgun',
      domain: 'example.com',
      privateKey: 'PRIVATE',
      publicKey: 'PUBLIC',
      sender: '',
      region: '',
    },
    smsService: {
      type: 'twilio',
      accountSid: 'asid',
      authToken: 'token',
      serviceSid: 'ssid',
    },
  },
  sessionStorage: {
    type: 'memory',
    sessionDuration: 300,
    address: '',
    password: '',
    db: '',
    region: '',
    endpoint: '',
  },
  staticFilesStorage: {
    type: 'local',
    serverConfigPath: 'server-config.yaml',
    region: '',
    bucket: '',
    endpoint: '',
    folder: '',
  },
  configurationStorage: {
    type: 'file',
    settingsKey: 'server-config.yaml',
    endpoints: [],
    region: '',
    bucket: 'bucket',
    keyStorage: {
      type: 'file',
      privateKey: 'jwt/private.pem',
      publicKey: 'jwt/public.pem',
      region: '',
      bucket: 'jwt bucket',
    },
  },
};

const createSettingsServiceMock = () => {
  const fetchLoginSettings = async () => {
    await pause(400);
    return data.login;
  };

  const updateLoginSettings = async (settings) => {
    await pause(400);
    data.login = settings;
  };

  const fetchExternalServicesSettings = async () => {
    await pause(400);
    return data.externalServices;
  };

  const updateExternalServicesSettings = async (settings) => {
    await pause(400);
    data.externalServices = settings;
  };

  const fetchSessionStorageSettings = async () => {
    await pause(400);
    return data.sessionStorage;
  };

  const updateSessionStorageSettings = async (settings) => {
    await pause(400);
    data.sessionStorage = settings;
  };

  const fetchStaticFilesSettings = async () => {
    await pause(400);
    return data.staticFiles;
  };

  const updateStaticFilesSettings = async (settings) => {
    await pause(400);
    data.staticFiles = settings;
  };

  const fetchGeneralSettings = async () => {
    await pause(400);
    return data.general;
  };

  const updateGeneralSettings = async (settings) => {
    await pause(400);
    data.general = settings;
  };

  const fetchConfigurationStorageSettings = async () => {
    await pause(400);
    return data.configurationStorage;
  };

  const updateConfigurationStorageSettings = async (settings) => {
    console.log(toDeepCase(settings, 'snake'));
    await pause(400);
    data.configurationStorage = settings;
  };

  const uploadJWTKeys = async () => {
    await pause(400);
  };

  const requestServerRestart = async () => {
    await pause(100);
  };

  return {
    fetchLoginSettings,
    updateLoginSettings,
    fetchExternalServicesSettings,
    updateExternalServicesSettings,
    fetchSessionStorageSettings,
    updateSessionStorageSettings,
    fetchStaticFilesSettings,
    updateStaticFilesSettings,
    fetchGeneralSettings,
    updateGeneralSettings,
    fetchConfigurationStorageSettings,
    updateConfigurationStorageSettings,
    uploadJWTKeys,
    requestServerRestart,
  };
};

export default createSettingsServiceMock;
