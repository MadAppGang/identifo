import actionCreator from '@madappgang/action-creator';
import { getError } from '~/utils';
import types from './types';

const fetchSettingsAttempt = actionCreator(types.FETCH_ACCOUNT_SETTINGS_ATTEMPT);
const fetchSettingsSuccess = actionCreator(types.FETCH_ACCOUNT_SETTINGS_SUCCESS);
const fetchSettingsFailure = actionCreator(types.FETCH_ACCOUNT_SETTINGS_FAILURE);

const postSettingsAttempt = actionCreator(types.POST_ACCOUNT_SETTINGS_ATTEMPT);
const postSettingsSuccess = actionCreator(types.POST_ACCOUNT_SETTINGS_SUCCESS);
const postSettingsFailure = actionCreator(types.POST_ACCOUNT_SETTINGS_FAILURE);

const fetchAccountSettings = () => async (dispatch, _, services) => {
  dispatch(fetchSettingsAttempt());

  try {
    const settings = await services.account.fetchSettings();
    dispatch(fetchSettingsSuccess(settings));
  } catch (err) {
    dispatch(fetchSettingsFailure(getError(err)));
  }
};

const postAccountSettings = settings => async (dispatch, _, services) => {
  dispatch(postSettingsAttempt());

  try {
    await services.account.postSettings(settings);
    dispatch(postSettingsSuccess(settings));
  } catch (err) {
    dispatch(postSettingsFailure(getError(err)));
  }
};

const resetAccountError = actionCreator(types.RESET_ACCOUNT_ERROR);

export {
  fetchAccountSettings,
  postAccountSettings,
  resetAccountError,
};
