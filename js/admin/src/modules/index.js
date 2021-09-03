import { combineReducers } from 'redux';
import authReducer from './auth/reducer';
import databaseReducer from './database/reducer';
import accountReducer from './account/reducer';
import userListReducer from './users/listReducer';
import selectedUserReducer from './users/selectedReducer';
import applicationListReducer from './applications/listReducer';
import selectedApplicationReducer from './applications/selectedReducer';
import dialogsApplicationReducer from './applications/dialogsReducer';
import settingsReducer from './settings/reducer';

import configureStore from './store';

import { curry } from '../utils/fn';

const rootReducer = combineReducers({
  auth: authReducer,
  database: databaseReducer,
  account: accountReducer,
  users: userListReducer,
  selectedUser: selectedUserReducer,
  applications: applicationListReducer,
  selectedApplication: selectedApplicationReducer,
  applicationDialogs: dialogsApplicationReducer,
  settings: settingsReducer,
});

export default curry(configureStore, rootReducer);
