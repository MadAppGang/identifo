import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import './ToggleView.css';

const ToggleDefaultView = ({ on, toggle, label }) => {
  const rootClassName = classnames({
    'iap-default-toggle__body': true,
    'iap-default-toggle__body--on': on,
  });

  return (
    <div className="iap-default-toggle">
      <button
        type="button"
        className={rootClassName}
        onClick={toggle}
      >
        <div className="iap-default-toggle__handle" />
      </button>
      <span>
        {label}
      </span>
    </div>
  );
};

ToggleDefaultView.propTypes = {
  toggle: PropTypes.func.isRequired,
  label: PropTypes.string,
  on: PropTypes.bool.isRequired,
};

ToggleDefaultView.defaultProps = {
  label: '',
};

export default ToggleDefaultView;
