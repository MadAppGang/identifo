import {
  RECEIVE_LOGIN_SETTINGS,
  FETCH_SERVER_SETTINGS,
  UPDATE_SERVER_SETTINGS,
} from './types';
import { logout } from '../auth/actions';
import actionCreator from '@madappgang/action-creator';
import { pause } from '~/utils';

const setServerSettings = actionCreator(FETCH_SERVER_SETTINGS);
export const updateServerSettings = actionCreator(UPDATE_SERVER_SETTINGS);

export const fetchServerSetings = () => async (dispatch, _, services) => {
  const settings = await services.settings.fetchServerSettings();
  dispatch(setServerSettings(settings));
};

export const updateLoginSettings = settings => async (dispatch, _, services) => {
  await services.settings.updateLoginSettings(settings);
  dispatch({
    type: RECEIVE_LOGIN_SETTINGS,
    payload: settings,
  });
};

export const uploadJWTKeys = (pubKey, privKey) => async (dispatch, _, services) => {
  await services.settings.uploadJWTKeys(pubKey, privKey);
};

export const restartServer = () => async (dispatch, _, services) => {
  await services.settings.requestServerRestart();
  await pause(1000);

  dispatch(logout());
};
