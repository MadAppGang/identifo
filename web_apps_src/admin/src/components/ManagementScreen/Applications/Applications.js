import React from 'react';
import { Switch, Route } from 'react-router-dom';
import ApplicationsMainView from './MainView';
import CreateApplicationView from './CreateApplicationView';
import EditApplicationView from './EditApplicationView';

const ApplicationsSection = () => {
  return (
    <Switch>
      <Route
        exact
        path="/management/applications/"
        render={props => <ApplicationsMainView {...props} />}
      />
      <Route
        path="/management/applications/new"
        render={props => <CreateApplicationView {...props} />}
      />
      <Route
        path="/management/applications/:appid"
        render={props => <EditApplicationView {...props} />}
      />
    </Switch>
  );
};

export default ApplicationsSection;
