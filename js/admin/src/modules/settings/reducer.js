import update from '@madappgang/update-by-path';
import {
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
      return { ...state, current: { ...state.current, ...payload } };
    case authTypes.LOGOUT_ATTEMPT:
      return update(state, 'changed', false);
    default:
      return state;
  }
};

export default reducer;
