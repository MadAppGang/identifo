import actionCreator from '@madappgang/action-creator';
import { getError } from '~/utils';
import types from './types';

const fetchAttempt = actionCreator(types.FETCH_USERS_ATTEMPT);
const fetchSuccess = actionCreator(types.FETCH_USERS_SUCCESS);
const fetchFailure = actionCreator(types.FETCH_USERS_FAILURE);

const postAttempt = actionCreator(types.POST_USER_ATTEMPT);
const postSuccess = actionCreator(types.POST_USER_SUCCESS);
const postFailure = actionCreator(types.POST_USER_FAILURE);

const alterAttempt = actionCreator(types.ALTER_USER_ATTEMPT);
const alterSuccess = actionCreator(types.ALTER_USER_SUCCESS);
const alterFailure = actionCreator(types.ALTER_USER_FAILURE);

const fetchByIdAttempt = actionCreator(types.FETCH_USER_BY_ID_ATTEMPT);
const fetchByIdSuccess = actionCreator(types.FETCH_USER_BY_ID_SUCCESS);
const fetchByIdFailure = actionCreator(types.FETCH_USER_BY_ID_FAILURE);

const deleteAttempt = actionCreator(types.DELETE_USER_BY_ID_ATTEMPT);
const deleteSuccess = actionCreator(types.DELETE_USER_BY_ID_SUCCESS);
const deleteFailure = actionCreator(types.DELETE_USER_BY_ID_FAILURE);

const fetchUsers = filters => async (dispatch, _, services) => {
  dispatch(fetchAttempt());

  try {
    const { users = [], total = 0 } = await services.users.fetchUsers(filters);
    dispatch(fetchSuccess({ users, total }));
  } catch (err) {
    dispatch(fetchFailure(getError(err)));
  }
};

const postUser = user => async (dispatch, _, services) => {
  dispatch(postAttempt());

  try {
    const result = await services.users.postUser(user);
    dispatch(postSuccess(result));
  } catch (err) {
    dispatch(postFailure(getError(err)));
  }
};

const alterUser = (id, changes) => async (dispatch, _, services) => {
  dispatch(alterAttempt());

  try {
    const user = await services.users.alterUser(id, changes);
    dispatch(alterSuccess(user));
  } catch (err) {
    dispatch(alterFailure(getError(err)));
  }
};

const fetchUserById = id => async (dispatch, _, services) => {
  dispatch(fetchByIdAttempt());

  try {
    const user = await services.users.fetchUserById(id);
    dispatch(fetchByIdSuccess(user));
  } catch (err) {
    dispatch(fetchByIdFailure(getError(err)));
  }
};

const deleteUserById = id => async (dispatch, _, services) => {
  dispatch(deleteAttempt());

  try {
    await services.users.deleteUserById(id);
    dispatch(deleteSuccess(id));
  } catch (err) {
    dispatch(deleteFailure(getError(err)));
  }
};

const resetUserError = actionCreator(types.RESET_USER_ERROR);

export {
  fetchUsers,
  postUser,
  alterUser,
  fetchUserById,
  deleteUserById,
  resetUserError,
};
