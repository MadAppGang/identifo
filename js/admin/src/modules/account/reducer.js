import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  fetching: false,
  posting: false,
  settings: null,
  error: null,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type } = action;

  switch (type) {
    case types.RESET_ACCOUNT_ERROR:
      return update(state, 'error', null);
    default:
      return state;
  }
};

export default reducer;
