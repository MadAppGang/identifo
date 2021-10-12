import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  authenticated: false,
  inProgress: false,
  error: null,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.LOGIN_ATTEMPT:
      return update(state, 'inProgress', true);
    case types.AUTH_STATE_CHANGE:
      return update(state, {
        inProgress: false,
        authenticated: payload,
      });
    case types.LOGIN_FAILURE:
      return update(state, {
        inProgress: false,
        error: payload,
      });
    case types.LOGOUT_ATTEMPT:
      return update(state, {
        inProgress: true,
        error: null,
      });
    default:
      return state;
  }
};

export default reducer;
