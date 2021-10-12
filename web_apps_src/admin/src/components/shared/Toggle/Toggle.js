import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import LoadingIcon from '~/components/icons/LoadingIcon';

let loadingTimeout;

const Toggle = ({ label, value, onChange }) => {
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (loading) {
      setLoading(false);
    }

    clearTimeout(loadingTimeout);
  }, [value]);

  const handleToggle = () => {
    onChange(!value);
    loadingTimeout = setTimeout(setLoading, 70, true);
  };

  const rootClassName = classnames({
    'iap-default-toggle__body': true,
    'iap-default-toggle__body--on': value,
  });

  return (
    <div className="iap-default-toggle">
      {!!label && (
        <span>
          {label}
        </span>
      )}
      <button
        type="button"
        className={rootClassName}
        onClick={handleToggle}
      >
        <div className="iap-default-toggle__handle">
          {loading && (
            <LoadingIcon className="iap-default-toggle__handle-icon" />
          )}
        </div>
      </button>
    </div>
  );
};

Toggle.propTypes = {
  label: PropTypes.string,
  value: PropTypes.bool,
  onChange: PropTypes.func,
};

Toggle.defaultProps = {
  label: '',
  value: false,
  onChange: () => {},
};

export default Toggle;
