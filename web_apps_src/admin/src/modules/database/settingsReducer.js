import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  config: null,
  fetching: false,
  posting: false,
  error: null,
  changed: false,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type } = action;

  switch (type) {
    case types.RESET_DB_SETTINGS_ERROR:
      return update(state, 'error', null);
    default:
      return state;
  }
};

export default reducer;
