import { combineReducers } from 'redux';
import settingsReducer from './settingsReducer';
import connectionReducer from './connectionReducer';

const databaseReducer = combineReducers({
  settings: settingsReducer,
  connection: connectionReducer,
});

export default databaseReducer;
