import update from '@madappgang/update-by-path';
import types from './types';

export const CONNECTION_ESTABLISHED = 'connection_established';
export const CONNECTION_FAILED = 'connection_failed';
export const CONNECTION_TEST_REQUIRED = 'connection_test_required';

const INITIAL_STATE = {
  checking: false,
  state: CONNECTION_TEST_REQUIRED,
};

const reducer = (state = INITIAL_STATE, { type }) => {
  switch (type) {
    case types.TEST_CONNECTION_ATTEMPT:
      return update(state, 'checking', true);
    case types.TEST_CONNECTION_SUCCESS:
      return update(state, {
        checking: false,
        state: CONNECTION_ESTABLISHED,
      });
    case types.TEST_CONNECTION_FAILURE:
      return update(state, {
        checking: false,
        state: CONNECTION_FAILED,
      });
    case types.FETCH_DB_SETTINGS_ATTEMPT:
    case types.POST_DB_SETTINGS_SUCCESS:
      return update(state, 'state', CONNECTION_TEST_REQUIRED);
    default:
      return state;
  }
};

export default reducer;
