import actionCreator from '@madappgang/action-creator';
import {
  FETCH_SERVER_SETTINGS,
  UPDATE_SERVER_SETTINGS,
  FETCH_JWT_KEYS,
  SET_VERIFICATION_STATUS,
} from './types';
import { pause } from '~/utils';
import { verificationStatuses } from '~/enums';
import { logout, authStateChange } from '~/modules/auth/actions';
import { commonAsyncHandler } from '~/utils/commonAsynсHandler';
import { showSuccessNotificationSnack } from '~/modules/applications/notification-actions';
import { successSnackMessages } from '~/modules/applications/constants';


const setServerSettings = actionCreator(FETCH_SERVER_SETTINGS);
export const updateServerSettings = actionCreator(UPDATE_SERVER_SETTINGS);
export const setJWTKeys = actionCreator(FETCH_JWT_KEYS);
export const setVerificationStatus = actionCreator(SET_VERIFICATION_STATUS);

export const getJWTKeys = withPrivate => async (dispatch, _, services) => {
  await commonAsyncHandler(async () => {
    const keys = await services.settings.getJWTKeys(withPrivate);
    dispatch(setJWTKeys(keys));
  }, dispatch);
};

export const uploadJWTKeys = payload => async (dispatch, getState, services) => {
  await commonAsyncHandler(async () => {
    const keys = await services.settings.uploadJWTKeys(payload);
    dispatch(setJWTKeys(keys));
    dispatch(showSuccessNotificationSnack(successSnackMessages.uploadKey));
  }, dispatch);
};

export const generateKeys = alg => async (dispatch, _, services) => {
  await commonAsyncHandler(async () => {
    const res = await services.settings.generateKeys({ alg });
    dispatch(setJWTKeys(res));
    dispatch(showSuccessNotificationSnack(successSnackMessages.generateKeys));
  }, dispatch);
};

export const fetchServerSetings = () => async (dispatch, _, services) => {
  commonAsyncHandler(async () => {
    const settings = await services.settings.fetchServerSettings();
    dispatch(setServerSettings(settings));
    await dispatch(getJWTKeys(false));
  }, dispatch);
};

export const postServerSettings = () => async (dispatch, getState, services) => {
  await commonAsyncHandler(async () => {
    const { jwtKeys, ...settings } = getState().settings.current;
    const res = await services.settings.postServerSettings(settings);
    dispatch(setServerSettings(res));
    dispatch(authStateChange(false));
  }, dispatch);
};

export const restartServer = () => async (dispatch, _, services) => {
  await services.settings.requestServerRestart();
  await pause(1000);

  dispatch(logout());
};

export const verifyConnection = settings => async (dispatch, _, services) => {
  dispatch(setVerificationStatus(verificationStatuses.loading));
  await commonAsyncHandler(async () => {
    await services.settings.verifySettings(settings);
    dispatch(setVerificationStatus(verificationStatuses.success));
  }, dispatch, true)
    .catch(() => {
      dispatch(setVerificationStatus(verificationStatuses.fail));
    });
};
