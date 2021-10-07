import React from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { MarkodwnScreen } from '~/components/ManagementScreen/MardownScreen/MardownScreen';
import { ServicesContext } from '../hooks/useServices';
import ensureAuthState, { SIGNED_IN, SIGNED_OUT } from './ensureAuthState';
import LoginScreen from './LoginScreen';
import ManagementScreen from './ManagementScreen';
import NotFoundScreen from './NotFoundScreen';
import './Root.css';


const Root = ({ store, services }) => {
  return (
    <Provider store={store}>
      <ServicesContext.Provider value={services}>
        <BrowserRouter basename={process.env.BASE_URL}>
          <Switch>
            <Route
              exact
              path="/"
              component={ensureAuthState(SIGNED_OUT, LoginScreen, '/management')}
            />
            <Route
              path="/management/:section?"
              component={ensureAuthState(SIGNED_IN, ManagementScreen, '/')}
            />
            <Route
              path="/faq"
              component={ensureAuthState(SIGNED_IN, MarkodwnScreen, '/')}
            />
            <Route component={NotFoundScreen} />
          </Switch>
        </BrowserRouter>
      </ServicesContext.Provider>
    </Provider>
  );
};

export default Root;
