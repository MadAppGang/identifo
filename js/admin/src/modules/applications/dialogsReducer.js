import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  settings: {
    show: false,
    config: null,
  },
  settingsSnack: {
    show: false,
  },
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.SHOW_SETTINGS_DIALOG:
      return update(state, 'settings', { show: true, config: payload });
    case types.HIDE_SETTINGS_DIALOG:
      return update(state, 'settings', { show: false, config: null });
    case types.SHOW_SETTINGS_SNACK:
      return update(state, 'settingsSnack', { show: true });
    case types.HIDE_SETTINGS_SNACK:
      return update(state, 'settingsSnack', { show: false });
    default:
      return state;
  }
};

export default reducer;
