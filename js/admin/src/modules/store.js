import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';

const configureStore = (rootReducer, services) => {
  const middleware = applyMiddleware(
    thunk.withExtraArgument(services),
  );

  return createStore(rootReducer, middleware);
};

export default configureStore;
