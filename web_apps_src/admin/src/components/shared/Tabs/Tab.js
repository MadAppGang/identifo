import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

const Tab = React.forwardRef(({ isActive, title, onClick }, ref) => {
  const className = classnames('iap-tabs-tab', {
    'iap-tabs-tab--active': isActive,
  });

  return (
    <button className={className} onClick={onClick} ref={ref}>
      {title}
    </button>
  );
});

Tab.propTypes = {
  isActive: PropTypes.bool,
  title: PropTypes.string.isRequired,
  onClick: PropTypes.func,
};

Tab.defaultProps = {
  onClick: () => {},
  isActive: false,
};

export default Tab;
