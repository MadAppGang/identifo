import { login, logout } from '../actions';
import types from '../types';

describe('auth module actions', () => {
  let auth;
  let dispatch;

  beforeEach(() => {
    dispatch = jest.fn();
    auth = {
      login: jest.fn(() => Promise.resolve()),
      logout: jest.fn(() => Promise.resolve()),
    };
  });

  test('dispatches login attempt on login', () => {
    const expectedAction = {
      type: types.LOGIN_ATTEMPT,
    };

    login('email', 'password')(dispatch, null, { auth });
    expect(dispatch).toHaveBeenCalledWith(expectedAction);
  });

  test('dispatches login success with credentials on successful login', async () => {
    const expectedAction = {
      type: types.AUTH_STATE_CHANGE,
      payload: true,
    };

    await login('email', 'password')(dispatch, null, { auth });
    expect(dispatch).toHaveBeenLastCalledWith(expectedAction);
  });

  test('dispatches login failure with error on failed login', async () => {
    const error = new Error();
    const expectedAction = {
      type: types.LOGIN_FAILURE,
      payload: error,
    };
    auth.login = jest.fn(() => Promise.reject(error));

    await login('email', 'password')(dispatch, null, { auth });
    expect(dispatch).toHaveBeenLastCalledWith(expectedAction);
  });

  test('invokes auth service on login', () => {
    login('email', 'password')(dispatch, null, { auth });
    expect(auth.login).toHaveBeenCalledWith('email', 'password');
  });

  test('dispatches logout attempt on logout', () => {
    const expectedAction = {
      type: types.LOGOUT_ATTEMPT,
    };

    logout()(dispatch, null, { auth });
    expect(dispatch).toHaveBeenCalledWith(expectedAction);
  });

  test('dispatches logout success on logout', async () => {
    const expectedAction = {
      type: types.AUTH_STATE_CHANGE,
      payload: false,
    };

    await logout()(dispatch, null, { auth });
    expect(dispatch).toHaveBeenLastCalledWith(expectedAction);
  });

  test('invokes auth service on logout', () => {
    logout()(dispatch, null, { auth });
    expect(auth.logout).toHaveBeenCalled();
  });
});
