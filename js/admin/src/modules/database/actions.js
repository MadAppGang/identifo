import actionCreator from '@madappgang/action-creator';
import { getError } from '~/utils';
import types from './types';

const testConnectionAttempt = actionCreator(types.TEST_CONNECTION_ATTEMPT);
const testConnectionSuccess = actionCreator(types.TEST_CONNECTION_SUCCESS);
const testConnectionFailure = actionCreator(types.TEST_CONNECTION_FAILURE);

export const resetConnectionState = actionCreator(types.RESET_CONNECTION_STORE);

const testConnection = () => async (dispatch, getState, services) => {
  dispatch(testConnectionAttempt());

  try {
    await services.database.testConnection(getState().database.settings.config);
    dispatch(testConnectionSuccess());
  } catch (err) {
    dispatch(testConnectionFailure(getError(err)));
  }
};

const verifyConnection = settings => async (dispatch, _, services) => {
  dispatch(testConnectionAttempt());
  try {
    await services.database.verifySettings(settings);
    dispatch(testConnectionSuccess());
  } catch (err) {
    dispatch(testConnectionFailure(new Error(getError(err))));
    throw new Error(getError(err));
  }
};

const resetError = actionCreator(types.RESET_DB_SETTINGS_ERROR);

export {
  testConnection,
  verifyConnection,
  resetError,
};
