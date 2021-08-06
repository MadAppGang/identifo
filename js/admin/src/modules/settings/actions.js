import {
  RECEIVE_LOGIN_SETTINGS,
  RECEIVE_EXTERNAL_SETTINGS,
  RECEIVE_SESSION_STORAGE_SETTINGS,
  RECEIVE_STATIC_FILES_SETTINGS,
  RECEIVE_GENERAL_SETTINGS,
  RECEIVE_CONFIGURATION_STORAGE_SETTINGS,
  SETTINGS_CHANGED,
} from './types';

import { logout } from '../auth/actions';
import { pause } from '~/utils';

export const fetchLoginSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchLoginSettings();
  dispatch({
    type: RECEIVE_LOGIN_SETTINGS,
    payload: settings,
  });
};

export const updateLoginSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateLoginSettings(settings);
  dispatch({
    type: RECEIVE_LOGIN_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const fetchExternalServicesSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchExternalServicesSettings();
  dispatch({
    type: RECEIVE_EXTERNAL_SETTINGS,
    payload: settings,
  });
};

export const updateExternalServicesSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateExternalServicesSettings(settings);
  dispatch({
    type: RECEIVE_EXTERNAL_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const fetchSessionStorageSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchSessionStorageSettings();
  dispatch({
    type: RECEIVE_SESSION_STORAGE_SETTINGS,
    payload: settings,
  });
};

export const updateSessionStorageSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateSessionStorageSettings(settings);
  dispatch({
    type: RECEIVE_SESSION_STORAGE_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const fetchStaticFilesSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchStaticFilesSettings();
  dispatch({
    type: RECEIVE_STATIC_FILES_SETTINGS,
    payload: settings,
  });
};

export const updateStaticFilesSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateStaticFilesSettings(settings);
  dispatch({
    type: RECEIVE_STATIC_FILES_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const fetchGeneralSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchGeneralSettings();
  dispatch({
    type: RECEIVE_GENERAL_SETTINGS,
    payload: settings,
  });
};

export const updateGeneralSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateGeneralSettings(settings);
  dispatch({
    type: RECEIVE_GENERAL_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const fetchConfigurationStorageSettings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchConfigurationStorageSettings();
  dispatch({
    type: RECEIVE_CONFIGURATION_STORAGE_SETTINGS,
    payload: settings,
  });
};

export const updateConfigurationStorageSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateConfigurationStorageSettings(settings);
  dispatch({
    type: RECEIVE_CONFIGURATION_STORAGE_SETTINGS,
    payload: settings,
  });

  dispatch({ type: SETTINGS_CHANGED });
};

export const uploadJWTKeys = (pubKey, privKey) => async (dispatch, _, services) => {
  await services.settings.uploadJWTKeys(pubKey, privKey);
};

export const restartServer = () => async (dispatch, _, services) => {
  await services.settings.requestServerRestart();
  await pause(1000);

  dispatch(logout());
};
