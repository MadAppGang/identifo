import React from 'react';
import PropTypes from 'prop-types';

import './index.css';

const PreviewField = ({ label, value }) => {
  if (!value) {
    return null;
  }

  return (
    <div className="iap-section__field">
      <span>
        {label}
      </span>
      <p className="iap-section__value">
        {value}
      </p>
    </div>
  );
};

PreviewField.propTypes = {
  label: PropTypes.string,
  value: PropTypes.string,
};

PreviewField.defaultProps = {
  label: '',
  value: '',
};

export default PreviewField;
