import actionCreator from '@madappgang/action-creator';
import {
  FETCH_SERVER_SETTINGS,
  UPDATE_SERVER_SETTINGS,
  FETCH_JWT_KEYS,
  SET_VERIFICATION_STATUS,
} from './types';
import { pause } from '~/utils';
import { verificationStatuses } from '~/enums';
import { logout, checkAuthState } from '~/modules/auth/actions';
import { showNotificationSnack } from '~/modules/applications/actions';
import { notificationStates } from '~/modules/applications/notificationsStates';

const setServerSettings = actionCreator(FETCH_SERVER_SETTINGS);
export const updateServerSettings = actionCreator(UPDATE_SERVER_SETTINGS);
export const setJWTKeys = actionCreator(FETCH_JWT_KEYS);
export const setVerificationStatus = actionCreator(SET_VERIFICATION_STATUS);

export const getJWTKeys = withPrivate => async (dispatch, _, services) => {
  try {
    const keys = await services.settings.getJWTKeys(withPrivate);
    dispatch(setJWTKeys(keys));
  } catch (error) {
    await dispatch(showNotificationSnack(notificationStates.error.status));
    throw new Error(error);
  }
};

export const uploadJWTKeys = payload => async (dispatch, getState, services) => {
  try {
    const keys = await services.settings.uploadJWTKeys(payload);
    await dispatch(setJWTKeys(keys));
  } catch (error) {
    throw new Error(error);
  }
};

export const generateKeys = alg => async (dispatch, _, services) => {
  try {
    const res = await services.settings.generateKeys({ alg });
    dispatch(setJWTKeys(res));
  } catch (error) {
    await dispatch(showNotificationSnack(notificationStates.error.status));
    throw new Error(error);
  }
};

export const fetchServerSetings = () => async (dispatch, _, services) => {
  try {
    const settings = await services.settings.fetchServerSettings();
    dispatch(setServerSettings(settings));
    await dispatch(getJWTKeys(false));
  } catch (error) {
    await dispatch(showNotificationSnack(notificationStates.error.status));
    throw new Error(error);
  }
};

export const postServerSettings = () => async (dispatch, getState, services) => {
  try {
    const { jwtKeys, ...settings } = getState().settings.current;
    const res = await services.settings.postServerSettings(settings);
    dispatch(setServerSettings(res));
    // TODO: Nikita K fix logout
    // await dispatch(showNotificationSnack(notificationStates.success.status));
    await dispatch(checkAuthState());
    // await dispatch(logout());
  } catch (error) {
    await dispatch(showNotificationSnack(notificationStates.error.status));
    throw new Error(error);
  }
};

export const restartServer = () => async (dispatch, _, services) => {
  await services.settings.requestServerRestart();
  await pause(1000);

  dispatch(logout());
};

export const verifyConnection = settings => async (dispatch, _, services) => {
  dispatch(setVerificationStatus(verificationStatuses.loading));
  try {
    await services.settings.verifySettings(settings);
    dispatch(setVerificationStatus(verificationStatuses.success));
  } catch (err) {
    dispatch(setVerificationStatus(verificationStatuses.fail));
    throw new Error(err);
  }
};
