import React from 'react';
import PropTypes from 'prop-types';
import FetchFailureIcon from '~/components/icons/FetchFailure';
import ReplayIcon from '~/components/icons/ReplayIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import Button from '~/components/shared/Button';

const SettingsPlaceholder = props => (
  <div className="iap-section-placeholder">
    <div className="iap-section-placeholder__title">
      Account
    </div>

    <FetchFailureIcon className="iap-placeholder__fetch-failure-icon" />

    <p className="iap-section-placeholder__msg">
      Could not load account settings
    </p>

    <Button
      error
      Icon={props.fetching ? LoadingIcon : ReplayIcon}
      onClick={props.onTryAgainClick}
      disabled={props.fetching}
    >
      Try again
    </Button>
  </div>
);

SettingsPlaceholder.propTypes = {
  fetching: PropTypes.bool,
  onTryAgainClick: PropTypes.func.isRequired,
};

SettingsPlaceholder.defaultProps = {
  fetching: false,
};

export default SettingsPlaceholder;
