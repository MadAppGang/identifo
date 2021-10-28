export const getCurentSettings = (state) => {
  return state.settings.current;
};

export const getOriginalSettings = (state) => {
  return state.settings.original;
};

export const getServerSectionSettings = (state) => {
  return getCurentSettings(state).configurationStorage;
};

export const getAdminPanelSettings = (state) => {
  return getCurentSettings(state).adminPanel;
};

export const getLoginWebAppSettings = (state) => {
  return getCurentSettings(state).loginWebApp;
};

export const getEmailTemplatesSettigns = (state) => {
  return getCurentSettings(state).emailTemplaits;
};

export const getGeneralServerSettings = (state) => {
  return getCurentSettings(state).general;
};

export const getKeyStorageSettings = (state) => {
  return getCurentSettings(state).keyStorage;
};

export const getStorageSettings = (state) => {
  return getCurentSettings(state).storage;
};

export const getLoginSettings = (state) => {
  return getCurentSettings(state).login;
};

export const getExternalServicesSettings = (state) => {
  return getCurentSettings(state).externalServices;
};

export const getStaticFilesSettings = (state) => {
  return getCurentSettings(state).staticFilesStorage;
};

export const getSettingsConfig = (state) => {
  return getCurentSettings(state).config;
};

export const getAdminAccountSettings = (state) => {
  return getCurentSettings(state).adminAccount;
};

export const getSessionStorageSettings = (state) => {
  return getCurentSettings(state).sessionStorage;
};

export const selectJWTKeys = (state) => {
  return getCurentSettings(state).jwtKeys;
};

export const getVerificationStatus = (state) => {
  return state.settings.verificationStatus;
};
