import actionCreator from '@madappgang/action-creator';
import { getError } from '~/utils';
import types from './types';
import { successSnackMessages } from '~/modules/applications/constants';
import { commonAsyncHandler } from '~/utils/commonAsynсHandler';
import { showSuccessNotificationSnack } from './notification-actions';

const fetchAttempt = actionCreator(types.FETCH_APPLICATIONS_ATTEMPT);
const fetchSuccess = actionCreator(types.FETCH_APPLICATIONS_SUCCESS);

const postAttempt = actionCreator(types.POST_APPLICATION_ATTEMPT);
const postSuccess = actionCreator(types.POST_APPLICATION_SUCCESS);

const deleteAttempt = actionCreator(types.DELETE_APPLICATION_ATTEMPT);
const deleteSuccess = actionCreator(types.DELETE_APPLICATION_SUCCESS);

const alterAttempt = actionCreator(types.ALTER_APPLICATION_ATTEMPT);
const alterSuccess = actionCreator(types.ALTER_APPLICATION_SUCCESS);

const fetchByIdAttempt = actionCreator(types.FETCH_APPLICATION_BY_ID_ATTEMPT);
const fetchByIdSuccess = actionCreator(types.FETCH_APPLICATION_BY_ID_SUCCESS);

const fetchFederatedProvidersAttempt = actionCreator(types.FETCH_FEDERATED_PROVIDERS_ATTEMTP);
const fetchFederatedProviderssSuccess = actionCreator(types.FETCH_FEDERATED_PROVIDERS_SUCCESS);
const fetchFederatedProviderssFailure = actionCreator(types.FETCH_FEDERATED_PROVIDERS_FAILURE);

const showSettingsDialog = actionCreator(types.SHOW_SETTINGS_DIALOG);
export const hideSettingsDialog = actionCreator(types.HIDE_SETTINGS_DIALOG);

export const showSettingsSnack = actionCreator(types.SHOW_SETTINGS_SNACK);
export const hideSettingsSnack = actionCreator(types.HIDE_SETTINGS_SNACK);

export const handleSettingsDialog = config => async (dispatch) => {
  return new Promise((resolve) => {
    const callback = (d) => {
      dispatch(hideSettingsDialog());
      resolve(d);
    };
    dispatch(showSettingsDialog({ ...config, callback }));
  });
};

const fetchApplications = () => async (dispatch, _, services) => {
  dispatch(fetchAttempt());
  await commonAsyncHandler(async () => {
    const { apps = [], total = 0 } = await services.applications.fetchApplications();
    dispatch(fetchSuccess({ apps, total }));
  }, dispatch);
};

const fetchApplicationById = id => async (dispatch, _, services) => {
  dispatch(fetchByIdAttempt());
  await commonAsyncHandler(async () => {
    const application = await services.applications.fetchApplicationById(id);
    dispatch(fetchByIdSuccess(application));
  }, dispatch);
};

const postApplication = application => async (dispatch, _, services) => {
  dispatch(postAttempt());
  await commonAsyncHandler(async () => {
    const result = await services.applications.postApplication(application);
    dispatch(postSuccess(result));
    dispatch(showSuccessNotificationSnack(successSnackMessages.postApp));
  }, dispatch);
};

const deleteApplicationById = id => async (dispatch, _, services) => {
  dispatch(deleteAttempt());
  await commonAsyncHandler(async () => {
    await services.applications.deleteApplicationById(id);
    dispatch(deleteSuccess(id));
    dispatch(showSuccessNotificationSnack(successSnackMessages.deleteApp));
  }, dispatch);
};

const alterApplication = (id, changes) => async (dispatch, _, services) => {
  dispatch(alterAttempt());
  await commonAsyncHandler(async () => {
    const result = await services.applications.alterApplication(id, changes);
    dispatch(alterSuccess(result));
    dispatch(showSuccessNotificationSnack(successSnackMessages.alterApp));
  }, dispatch);
};

const fetchFederatedProviders = () => async (dispatch, _, services) => {
  dispatch(fetchFederatedProvidersAttempt());

  try {
    const result = await services.applications.fetchFederatedLoginProviders();
    dispatch(fetchFederatedProviderssSuccess(result));
  } catch (err) {
    dispatch(fetchFederatedProviderssFailure(getError(err)));
  }
};

const resetApplicationError = actionCreator(types.RESET_APPLICATION_ERROR);

export {
  fetchApplications,
  postApplication,
  deleteApplicationById,
  alterApplication,
  fetchApplicationById,
  resetApplicationError,
  fetchFederatedProviders,
};
