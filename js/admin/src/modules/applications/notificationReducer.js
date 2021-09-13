import update from '@madappgang/update-by-path';
import types from './types';
import { notificationStates } from '~/modules/applications/notificationsStates';

const INITIAL_STATE = {
  settingsDialog: {
    show: false,
    config: null,
  },

  settingsSnack: {
    show: false,
  },

  notificationSnack: {
    show: false,
    config: null,
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
    case types.SHOW_NOTIFICATION_SNACK:
      return update(state, 'notificationSnack', { show: true, config: notificationStates[payload] });
    case types.HIDE_NOTIFICATION_SNACK:
      return update(state, 'notificationSnack', { show: false, config: null });
    default:
      return state;
  }
};

export default reducer;
