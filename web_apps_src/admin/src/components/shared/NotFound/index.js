import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import FetchFailureIcon from '~/components/icons/FetchFailure';
import Button from '~/components/shared/Button';
import './NotFound.css';

const NotFound = props => (
  <div className="iap-404-screen">
    <p className="iap-404-screen__title">404</p>
    <p className="iap-404-screen__subtitle">
      Page not found
    </p>
    <FetchFailureIcon className="iap-404-screen__icon" />
    <Button onClick={() => props.history.push('/management')}>
      ← &nbsp; Go to Home page
    </Button>
  </div>
);

NotFound.propTypes = {
  history: PropTypes.shape({
    push: PropTypes.func,
  }).isRequired,
};

export default withRouter(NotFound);
