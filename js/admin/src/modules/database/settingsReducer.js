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
  const { type, payload } = action;

  switch (type) {
    case types.FETCH_DB_SETTINGS_ATTEMPT:
      return update(state, 'fetching', true);
    case types.FETCH_DB_SETTINGS_SUCCESS:
      return update(state, {
        fetching: false,
        config: payload,
        error: null,
      });
    case types.POST_DB_SETTINGS_ATTEMPT:
      return update(state, 'posting', true);
    case types.POST_DB_SETTINGS_SUCCESS:
      return update(state, {
        error: null,
        posting: false,
        config: payload,
      });
    case types.POST_DB_SETTINGS_FAILURE:
    case types.FETCH_DB_SETTINGS_FAILURE:
      return update(state, {
        posting: false,
        fetching: false,
        error: payload,
      });
    case types.RESET_DB_SETTINGS_ERROR:
      return update(state, 'error', null);
    default:
      return state;
  }
};

export default reducer;
