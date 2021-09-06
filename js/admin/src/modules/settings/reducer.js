import update from '@madappgang/update-by-path';
import {
  RECEIVE_LOGIN_SETTINGS,
  RECEIVE_EXTERNAL_SETTINGS,
  RECEIVE_STATIC_FILES_SETTINGS,
  RECEIVE_GENERAL_SETTINGS,
  RECEIVE_CONFIGURATION_STORAGE_SETTINGS,
  SETTINGS_CHANGED,
  FETCH_SERVER_SETTINGS,
  UPDATE_SERVER_SETTINGS,
} from './types';
import authTypes from '../auth/types';

const initialSettings = {
  login: {
    loginWith: {
      username: false,
      federated: false,
      phone: false,
    },
    tfaType: 'app',
  },
  configurationStorage: null,
  sessionStorage: null,
  staticFilesStorage: null,
  general: null,
  changed: false,
  adminAccount: null,
  config: null,
  externalServices: null,
  keyStorage: null,
  logger: null,
  storage: null,
};

const INITIAL_STATE = {
  original: initialSettings,
  current: initialSettings,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case FETCH_SERVER_SETTINGS:
      return { ...state, original: payload, current: payload };
    case UPDATE_SERVER_SETTINGS:
      return update(state, 'current', payload);
    case RECEIVE_LOGIN_SETTINGS:
      return update(state, 'login', payload);
    case RECEIVE_EXTERNAL_SETTINGS:
      return update(state, 'externalServices', payload);
    case RECEIVE_STATIC_FILES_SETTINGS:
      return update(state, 'staticFilesStorage', payload);
    case RECEIVE_GENERAL_SETTINGS:
      return update(state, 'general', payload);
    case RECEIVE_CONFIGURATION_STORAGE_SETTINGS:
      return update(state, 'configurationStorage', payload);
    case SETTINGS_CHANGED:
      return update(state, 'changed', true);
    case authTypes.LOGOUT_ATTEMPT:
      return update(state, 'changed', false);
    default:
      return state;
  }
};

export default reducer;
