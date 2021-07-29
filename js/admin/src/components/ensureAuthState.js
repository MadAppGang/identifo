import React from 'react';
import PropTypes from 'prop-types';
import { useSelector  } from 'react-redux';
import { Redirect, withRouter } from 'react-router-dom';

const SIGNED_IN = true;
const SIGNED_OUT = false;

const ensureAuthState = (expectedAuthState, Component, redirectPath) => {
  const ConnectedComponent = ({ location, ...props }) => {
    const actualAuthState = useSelector(state => state.auth.authenticated);

    if (expectedAuthState !== actualAuthState) {
      if (expectedAuthState === SIGNED_IN) {
        const to = {
          pathname: redirectPath,
          state: {
            path: location.pathname,
          },
        };

        return <Redirect to={to} />;
      }

      const previousAttemptPath = (location.state || {}).path;

      return <Redirect to={previousAttemptPath || redirectPath} />;
    }

    return <Component {...props} />;
  };

  ConnectedComponent.propTypes = {
    location: PropTypes.shape({
      pathname: PropTypes.string,
      state: PropTypes.object,
    }).isRequired,
  };

  return withRouter(ConnectedComponent);
};

export {
  SIGNED_IN, SIGNED_OUT,
};

export default ensureAuthState;
