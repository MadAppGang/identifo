import {
  FETCH_SERVER_SETTINGS,
  UPDATE_SERVER_SETTINGS,
  FETCH_JWT_KEYS,
  SET_VERIFICATION_STATUS,
} from './types';
import authTypes from '../auth/types';
import { verificationStatuses } from '~/enums';

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
  jwtKeys: null,
  adminPanel: null,
};

const INITIAL_STATE = {
  original: initialSettings,
  current: initialSettings,
  verificationStatus: verificationStatuses.required,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case FETCH_SERVER_SETTINGS:
      return { ...state, original: payload, current: payload };
    case UPDATE_SERVER_SETTINGS:
      return { ...state, current: { ...state.current, ...payload } };
    case FETCH_JWT_KEYS:
      return {
        ...state,
        original: { ...state.original, jwtKeys: payload },
        current: { ...state.current, jwtKeys: payload },
      };
    case SET_VERIFICATION_STATUS:
      return { ...state, verificationStatus: payload };
    case authTypes.LOGOUT_ATTEMPT:
      return INITIAL_STATE;
    default:
      return state;
  }
};

export default reducer;
