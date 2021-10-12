import createHttpClient from './HttpClient';
import createAuthService from './Auth';
import createDatabaseService from './Database';
import createAccountService from './Account';
import createUserService from './Users';
import createApplicationService from './Applications';
import createSettingsService from './Settings';
import createAppleService from './Apple';
import createStaticService from './Static';

import createAuthServiceMock from './Auth/mock';
import createDatabaseServiceMock from './Database/mock';
import createAccountServiceMock from './Account/mock';
import createUserServiceMock from './Users/mock';
import createApplicationServiceMock from './Applications/mock';
import createSettingsServiceMock from './Settings/mock';
import createAppleServiceMock from './Apple/mock';
import createStaticServiceMock from './Static/mock';

const httpClient = createHttpClient();

const authService = createAuthService({ httpClient });
const databaseService = createDatabaseService({ httpClient });
const accountService = createAccountService({ httpClient });
const userService = createUserService({ httpClient });
const applicationService = createApplicationService({ httpClient });
const settingsService = createSettingsService({ httpClient });
const appleService = createAppleService({ httpClient });
const staticService = createStaticService({ httpClient });

const useMock = process.env.MOCK_API === 'true';

const services = {
  http: httpClient,
  auth: useMock ? createAuthServiceMock() : authService,
  database: useMock ? createDatabaseServiceMock() : databaseService,
  account: useMock ? createAccountServiceMock() : accountService,
  users: useMock ? createUserServiceMock() : userService,
  applications: useMock ? createApplicationServiceMock() : applicationService,
  settings: useMock ? createSettingsServiceMock() : settingsService,
  apple: useMock ? createAppleServiceMock() : appleService,
  static: useMock ? createStaticServiceMock() : staticService,
};

export default services;
