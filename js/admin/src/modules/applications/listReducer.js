import update from '@madappgang/update-by-path';
import types from './types';

const INITIAL_STATE = {
  fetching: false,
  list: [],
  federatedProviders: {},
  error: null,
  total: 0,
};

const reducer = (state = INITIAL_STATE, action) => {
  const { type, payload } = action;

  switch (type) {
    case types.FETCH_APPLICATIONS_ATTEMPT:
      return update(state, {
        fetching: true,
      });
    case types.FETCH_APPLICATIONS_SUCCESS:
      return update(state, {
        fetching: false,
        list: payload.apps,
        total: payload.total,
      });
    case types.FETCH_APPLICATIONS_FAILURE:
      return update(state, {
        fetching: false,
        error: payload,
      });
    case types.DELETE_APPLICATION_SUCCESS:
      return update(state, {
        list: list => list.filter(app => app.id !== payload),
      });
    case types.FETCH_FEDERATED_PROVIDERS_ATTEMTP:
      return update(state, {
        fetching: true,
      });
    case types.FETCH_FEDERATED_PROVIDERS_SUCCESS:
      return update(state, {
        fetching: false,
        federatedProviders: payload,
      });
    case types.FETCH_FEDERATED_PROVIDERS_FAILURE:
      return update(state, {
        fetching: false,
        error: payload,
      });
    default:
      return state;
  }
};

export default reducer;
