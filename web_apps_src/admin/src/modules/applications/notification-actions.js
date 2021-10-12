import actionCreator from '@madappgang/action-creator';
import types from './types';

export const showSuccessNotificationSnack = actionCreator(types.SHOW_SUCESS_NOTIFICATION_SNACK);
export const showErrorNotificationSnack = actionCreator(types.SHOW_ERROR_NOTIFICATION_SNACK);
export const hideNotificationSnack = actionCreator(types.HIDE_NOTIFICATION_SNACK);
export const setLoadingStatus = actionCreator(types.SET_LOADING_STATUS);
