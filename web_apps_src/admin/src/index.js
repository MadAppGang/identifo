import React from 'react';
import ReactDOM from 'react-dom';
import Root from './components/Root';
import services from './services';
import configureStore from './modules';

const store = configureStore(services);

const markup = (
  <Root store={store} services={services} />
);

ReactDOM.render((markup), document.getElementById('root'));
