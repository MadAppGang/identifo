import React from 'react';
import PropTypes from 'prop-types';
import './Field.css';

const Field = ({ label, subtext, Icon, children }) => (
  <div className="form-field">
    <span className="form-field__label">
      {label}

      {!!Icon && (
        <Icon.type {...Icon.props} className="form-field__label-icon" />
      )}
    </span>
    {children}
    {!!subtext && (
      <p className="form-field__subtext">{subtext}</p>
    )}
  </div>
);

Field.propTypes = {
  label: PropTypes.string,
  children: PropTypes.node,
  Icon: PropTypes.element,
  subtext: PropTypes.string,
};

Field.defaultProps = {
  label: '',
  children: null,
  Icon: null,
  subtext: '',
};

export default Field;
