import React from 'react';
import { Switch, Route } from 'react-router-dom';
import UsersMainView from './MainView';
import NewUserView from './NewUserView';
import EditUserView from './EditUserView';

const UsersSection = () => {
  return (
    <Switch>
      <Route
        exact
        path="/management/users"
        render={props => <UsersMainView {...props} />}
      />
      <Route
        path="/management/users/new"
        render={props => <NewUserView {...props} />}
      />
      <Route
        path="/management/users/:userid"
        render={props => <EditUserView {...props} />}
      />
    </Switch>
  );
};

export default UsersSection;
