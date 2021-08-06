import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  fetching: false,
  posting: false,
  settings: null,
  error: null,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.FETCH_ACCOUNT_SETTINGS_ATTEMPT:
      return update(state, 'fetching', true);
    case types.FETCH_ACCOUNT_SETTINGS_SUCCESS:
      return update(state, {
        fetching: false,
        settings: payload,
        error: null,
      });
    case types.FETCH_ACCOUNT_SETTINGS_FAILURE:
      return update(state, {
        fetching: false,
        error: payload,
      });
    case types.POST_ACCOUNT_SETTINGS_ATTEMPT:
      return update(state, 'posting', true);
    case types.POST_ACCOUNT_SETTINGS_SUCCESS:
      return update(state, {
        posting: false,
        settings: payload,
        error: null,
      });
    case types.POST_ACCOUNT_SETTINGS_FAILURE:
      return update(state, {
        posting: false,
        error: payload,
      });
    case types.RESET_ACCOUNT_ERROR:
      return update(state, 'error', null);
    default:
      return state;
  }
};

export default reducer;
