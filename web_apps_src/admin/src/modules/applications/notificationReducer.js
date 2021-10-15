import update from '@madappgang/update-by-path';
import types from './types';
import { notificationStatuses } from '~/enums';

const INITIAL_STATE = {
  loading: false,
  settingsDialog: {
    show: false,
    config: null,
  },

  settingsSnack: {
    show: false,
  },

  notificationSnack: {
    message: '',
    status: notificationStatuses.idle,
  },
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.SHOW_SETTINGS_DIALOG:
      return update(state, 'settingsDialog', { show: true, config: payload });
    case types.HIDE_SETTINGS_DIALOG:
      return update(state, 'settingsDialog', { show: false, config: null });
    case types.SHOW_SETTINGS_SNACK:
      return update(state, 'settingsSnack', { show: true });
    case types.HIDE_SETTINGS_SNACK:
      return update(state, 'settingsSnack', { show: false });
    case types.SHOW_SUCESS_NOTIFICATION_SNACK:
      return update(state, 'notificationSnack', { status: notificationStatuses.success, message: payload });
    case types.SHOW_ERROR_NOTIFICATION_SNACK:
      return update(state, 'notificationSnack', { status: notificationStatuses.error, message: payload });
    case types.HIDE_NOTIFICATION_SNACK:
      return update(state, 'notificationSnack', { status: notificationStatuses.idle, message: '' });
    case types.SET_LOADING_STATUS:
      return update(state, 'loading', payload);
    default:
      return state;
  }
};

export default reducer;
