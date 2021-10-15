import React from 'react';
import PropTypes from 'prop-types';
import ErrorIcon from '~/components/icons/ErrorIcon';
import './FormErrorMessage.css';

const FormErrorMessage = ({ error }) => (
  <div className="iap-management-section-error">
    <ErrorIcon className="iap-management-section-error__icon" />
    <div className="iap-management-section-error__msg">
      <p className="iap-management-section-error__title">
        Server error:
      </p>
      <p>{error.message}</p>
    </div>
  </div>
);

FormErrorMessage.propTypes = {
  error: PropTypes.instanceOf(Error).isRequired,
};

export default FormErrorMessage;
