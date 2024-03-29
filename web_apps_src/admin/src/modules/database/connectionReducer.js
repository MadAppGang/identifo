import update from '@madappgang/update-by-path';
import types from './types';

export const CONNECTION_ESTABLISHED = 'connection_established';
export const CONNECTION_FAILED = 'connection_failed';
export const CONNECTION_SUCCEED = 'connection_succeed';
export const CONNECTION_TEST_REQUIRED = 'connection_test_required';

const INITIAL_STATE = {
  checking: false,
  state: CONNECTION_TEST_REQUIRED,
  connectionStatus: CONNECTION_TEST_REQUIRED,
};

const reducer = (state = INITIAL_STATE, { type }) => {
  switch (type) {
    case types.TEST_CONNECTION_ATTEMPT:
      return update(state, 'checking', true);
    case types.TEST_CONNECTION_SUCCESS:
      return update(state, {
        checking: false,
        state: CONNECTION_SUCCEED,
      });
    case types.TEST_CONNECTION_FAILURE:
      return update(state, {
        checking: false,
        state: CONNECTION_FAILED,
      });
    case types.RESET_CONNECTION_STORE:
      return INITIAL_STATE;
    case types.FETCH_DB_SETTINGS_ATTEMPT:
      return update(state, 'state', CONNECTION_TEST_REQUIRED);
    default:
      return state;
  }
};

export default reducer;
