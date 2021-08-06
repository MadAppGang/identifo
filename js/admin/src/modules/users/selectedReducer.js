import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  user: null,
  error: null,
  fetching: false,
  saving: false,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.FETCH_USER_BY_ID_ATTEMPT:
      return update(state, {
        fetching: true,
        user: null,
        error: null,
      });
    case types.FETCH_USER_BY_ID_SUCCESS:
      return update(state, {
        fetching: false,
        user: payload,
      });
    case types.FETCH_USER_BY_ID_FAILURE:
      return update(state, {
        fetching: false,
        error: payload,
      });
    case types.POST_USER_ATTEMPT:
      return update(state, {
        saving: true,
      });
    case types.POST_USER_SUCCESS:
      return update(state, {
        error: null,
        saving: false,
        user: payload,
      });
    case types.POST_USER_FAILURE:
      return update(state, {
        saving: false,
        error: payload,
      });
    case types.ALTER_USER_ATTEMPT:
      return update(state, {
        saving: true,
      });
    case types.ALTER_USER_SUCCESS:
      return update(state, {
        saving: false,
        user: payload,
        error: null,
      });
    case types.ALTER_USER_FAILURE:
      return update(state, {
        saving: false,
        error: payload,
      });
    case types.DELETE_USER_BY_ID_ATTEMPT:
      return update(state, 'saving', true);
    case types.DELETE_USER_BY_ID_SUCCESS:
      return update(state, {
        saving: false,
        user: null,
      });
    case types.DELETE_USER_BY_ID_FAILURE:
      return update(state, {
        saving: false,
        error: payload,
      });
    case types.RESET_USER_ERROR:
      return update(state, 'error', null);
    default:
      return state;
  }
};


export default reducer;
