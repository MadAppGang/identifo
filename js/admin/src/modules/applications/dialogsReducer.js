import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  settings: {
    show: false,
    config: null,
  },
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.SHOW_SETTINGS_DIALOG:
      return update(state, 'settings', { show: true, config: payload });
    case types.HIDE_SETTINGS_DIALOG:
      return update(state, 'settings', { show: false, config: null });
    default:
      return state;
  }
};

export default reducer;
