import actionCreator from '@madappgang/action-creator';
import { getError } from '~/utils';
import types from './types';

const postSettingsAttempt = actionCreator(types.POST_ACCOUNT_SETTINGS_ATTEMPT);
const postSettingsSuccess = actionCreator(types.POST_ACCOUNT_SETTINGS_SUCCESS);
const postSettingsFailure = actionCreator(types.POST_ACCOUNT_SETTINGS_FAILURE);

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
  postAccountSettings,
  resetAccountError,
};
