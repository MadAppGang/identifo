import update from '@madappgang/update-by-path';
import {
  RECEIVE_LOGIN_SETTINGS,
  RECEIVE_EXTERNAL_SETTINGS,
  RECEIVE_SESSION_STORAGE_SETTINGS,
  RECEIVE_STATIC_FILES_SETTINGS,
  RECEIVE_GENERAL_SETTINGS,
  RECEIVE_CONFIGURATION_STORAGE_SETTINGS,
  SETTINGS_CHANGED,
} from './types';
import authTypes from '../auth/types';

const INITIAL_STATE = {
  login: {
    loginWith: {
      username: false,
      federated: false,
      phone: false,
    },
    tfaType: 'app',
  },
  configurationStorage: null,
  externalServices: null,
  sessionStorage: null,
  staticFiles: null,
  general: null,
  changed: false,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case RECEIVE_LOGIN_SETTINGS:
      return update(state, 'login', payload);
    case RECEIVE_EXTERNAL_SETTINGS:
      return update(state, 'externalServices', payload);
    case RECEIVE_SESSION_STORAGE_SETTINGS:
      return update(state, 'sessionStorage', payload);
    case RECEIVE_STATIC_FILES_SETTINGS:
      return update(state, 'staticFiles', payload);
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
