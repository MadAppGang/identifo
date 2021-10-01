import actionCreator from '@madappgang/action-creator';
import { commonAsyncHandler } from '../../utils/commonAsynсHandler';
import { successSnackMessages } from '../applications/constants';
import { showSuccessNotificationSnack } from '../applications/notification-actions';
import types from './types';

const fetchAttempt = actionCreator(types.FETCH_USERS_ATTEMPT);
const fetchSuccess = actionCreator(types.FETCH_USERS_SUCCESS);

const postAttempt = actionCreator(types.POST_USER_ATTEMPT);
const postSuccess = actionCreator(types.POST_USER_SUCCESS);

const alterAttempt = actionCreator(types.ALTER_USER_ATTEMPT);
const alterSuccess = actionCreator(types.ALTER_USER_SUCCESS);

const fetchByIdAttempt = actionCreator(types.FETCH_USER_BY_ID_ATTEMPT);
const fetchByIdSuccess = actionCreator(types.FETCH_USER_BY_ID_SUCCESS);

const deleteAttempt = actionCreator(types.DELETE_USER_BY_ID_ATTEMPT);
const deleteSuccess = actionCreator(types.DELETE_USER_BY_ID_SUCCESS);

const resetUserError = actionCreator(types.RESET_USER_ERROR);
const resetUserById = actionCreator(types.RESET_USER_BY_ID);

const fetchUsers = filters => async (dispatch, _, services) => {
  dispatch(fetchAttempt());
  await commonAsyncHandler(async () => {
    const { users = [], total = 0 } = await services.users.fetchUsers(filters);
    dispatch(fetchSuccess({ users, total }));
  }, dispatch);
};

const postUser = user => async (dispatch, _, services) => {
  dispatch(postAttempt());
  await commonAsyncHandler(async () => {
    const result = await services.users.postUser(user);
    dispatch(postSuccess(result));
    dispatch(showSuccessNotificationSnack(successSnackMessages.postUser));
  }, dispatch);
};

const alterUser = (id, changes) => async (dispatch, _, services) => {
  dispatch(alterAttempt());
  await commonAsyncHandler(async () => {
    const user = await services.users.alterUser(id, changes);
    dispatch(alterSuccess(user));
    dispatch(showSuccessNotificationSnack(successSnackMessages.alterUser));
  }, dispatch);
};

const fetchUserById = id => async (dispatch, _, services) => {
  dispatch(fetchByIdAttempt());
  await commonAsyncHandler(async () => {
    const user = await services.users.fetchUserById(id);
    dispatch(fetchByIdSuccess(user));
  }, dispatch);
};

const deleteUserById = id => async (dispatch, _, services) => {
  dispatch(deleteAttempt());
  await commonAsyncHandler(async () => {
    await services.users.deleteUserById(id);
    dispatch(deleteSuccess(id));
    dispatch(showSuccessNotificationSnack(successSnackMessages.deleteUser));
  }, dispatch);
};

export {
  fetchUsers,
  postUser,
  alterUser,
  fetchUserById,
  deleteUserById,
  resetUserError,
  resetUserById,
};
