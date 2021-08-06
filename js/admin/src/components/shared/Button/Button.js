import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

const Button = (props) => {
  const {
    stretch, Icon, children, error, transparent, iconClassName, extraClassName,
    ...domProps } = props;

  const className = classnames({
    'iap-btn': true,
    'iap-btn--stretch': stretch,
    'iap-btn--iconized': !!Icon,
    'iap-btn--transparent': transparent,
    'iap-btn--error': error,
    [extraClassName]: !!extraClassName,
  });

  return (
    <button
      className={className}
      {...domProps}
    >
      {Icon && (
        <Icon className={'iap-btn__icon '.concat(iconClassName).trim()} />
      )}
      <span>
        {children}
      </span>
    </button>
  );
};

Button.propTypes = {
  type: PropTypes.string,
  children: PropTypes.node,
  onClick: PropTypes.func,
  disabled: PropTypes.bool,
  stretch: PropTypes.bool,
  transparent: PropTypes.bool,
  iconClassName: PropTypes.string,
  Icon: PropTypes.func,
  error: PropTypes.bool,
  extraClassName: PropTypes.string,
};

Button.defaultProps = {
  type: 'button',
  onClick: null,
  children: null,
  disabled: false,
  stretch: false,
  transparent: false,
  iconClassName: '',
  Icon: null,
  error: false,
  extraClassName: '',
};

export default Button;
