import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

const SelectOption = ({ title, active, onClick }) => {
  const className = classnames({
    'iap-db-dropdown__option': true,
    'iap-db-dropdown__option--active': active,
  });

  return (
    <button
      type="button"
      className={className}
      onClick={onClick}
    >
      {title}
    </button>
  );
};

SelectOption.propTypes = {
  title: PropTypes.string,
  active: PropTypes.bool,
  onClick: PropTypes.func,
};

SelectOption.defaultProps = {
  title: '',
  active: false,
  onClick: () => {},
};

export default SelectOption;
