import { getError } from '~/utils';
import { setLoadingStatus, showErrorNotificationSnack } from '~/modules/applications/notification-actions';

const withTryCath = (callback, dispatch, withThrowError) => async () => {
  try {
    return await callback();
  } catch (err) {
    dispatch(showErrorNotificationSnack(getError(err).message));
    if (withThrowError) throw new Error(getError(err));
    return Promise.resolve();
  }
};

const withProcess = async (callback, dispatch) => {
  dispatch(setLoadingStatus(true));
  const response = await callback();
  dispatch(setLoadingStatus(false));
  return response;
};

export const commonAsyncHandler = async (callback, dispatch, withThrowError = false) => {
  const wrappedAction = withTryCath(callback, dispatch, withThrowError);
  const response = await withProcess(wrappedAction, dispatch);
  return response;
};
