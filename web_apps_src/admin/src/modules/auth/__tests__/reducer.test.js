import reducer from '../reducer';
import types from '../types';

describe('auth module reducer', () => {
  test('sign in progress indicator is false by default', () => {
    expect(reducer(undefined, {}).inProgress).toBe(false);
  });

  test('authentication state is false by default', () => {
    expect(reducer(undefined, {}).authenticated).toBe(false);
  });

  test('error is absent by default', () => {
    expect(reducer(undefined, {}).error).toBe(null);
  });

  test('sign out progress indicator is false by default', () => {
    expect(reducer(undefined, {}).inProgress).toBe(false);
  });

  test('sets sign in progress indicator to true on login attempt', () => {
    const action = { type: types.LOGIN_ATTEMPT };
    expect(reducer(undefined, action).inProgress).toBe(true);
  });

  test('sets sign in progress indicator to false on login success', () => {
    const action = { type: types.AUTH_STATE_CHANGE, payload: false };
    expect(reducer(undefined, action).inProgress).toBe(false);
  });

  test('sets auth state to true on login success', () => {
    const action = { type: types.AUTH_STATE_CHANGE, payload: true };
    expect(reducer(undefined, action).authenticated).toBe(true);
  });

  test('sets sign in progress indicator to false on login failure', () => {
    const action = { type: types.LOGIN_FAILURE };
    expect(reducer(undefined, action).inProgress).toBe(false);
  });

  test('sets error on login failure', () => {
    const error = new Error();
    const action = { type: types.LOGIN_FAILURE, payload: error };
    expect(reducer(undefined, action).error).toBe(error);
  });

  test('sets sign out progress indicator to true on logout attempt', () => {
    const action = { type: types.LOGOUT_ATTEMPT };
    expect(reducer(undefined, action).inProgress).toBe(true);
  });

  test('sets sign out progress indicator to true on logout attempt', () => {
    const action = { type: types.LOGOUT_ATTEMPT };
    expect(reducer(undefined, action).inProgress).toBe(true);
  });

  test('sets sign out progress indicator to false on logout success', () => {
    const action = { type: types.LOGOUT_SUCCESS };
    expect(reducer(undefined, action).inProgress).toBe(false);
  });

  test('sets auth state to false on logout success', () => {
    const action = { type: types.LOGOUT_SUCCESS };
    expect(reducer(undefined, action).authenticated).toBe(false);
  });
});
