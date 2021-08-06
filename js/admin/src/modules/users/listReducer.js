import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  list: [],
  error: null,
  fetching: false,
  total: 0,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.FETCH_USERS_ATTEMPT:
      return update(state, {
        fetching: true,
      });
    case types.FETCH_USERS_SUCCESS:
      return update(state, {
        fetching: false,
        list: payload.users,
        total: payload.total,
      });
    case types.FETCH_USERS_FAILURE:
      return update(state, {
        fetching: false,
        error: payload,
      });
    default:
      return state;
  }
};

export default reducer;
