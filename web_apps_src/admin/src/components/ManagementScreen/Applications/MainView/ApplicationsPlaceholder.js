import React from 'react';
import PropTypes from 'prop-types';
import Button from '~/components/shared/Button';
import ApplicationsIcon from '~/components/icons/ApplicationsIcon';
import AddIcon from '~/components/icons/AddIcon';

const ApplicationsPlaceholder = (props) => {
  return (
    <div className="iap-section-placeholder">
      <h2 className="iap-section-placeholder__title">
        Applications
      </h2>

      <ApplicationsIcon className="iap-section-placeholder__icon" />

      <p className="iap-section-placeholder__msg">
        No applications have been added so far.
      </p>

      <Button Icon={AddIcon} onClick={props.onCreateApplicationClick}>
        Create application
      </Button>
    </div>
  );
};

ApplicationsPlaceholder.propTypes = {
  onCreateApplicationClick: PropTypes.func.isRequired,
};

export default ApplicationsPlaceholder;
