import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import ErrorIcon from '~/components/icons/ErrorIcon';
import SuccessIcon from '~/components/icons/SuccessIcon';

const types = {
  SUCCESS: 'success',
  FAILURE: 'failure',
};

const Notification = ({ title, text, type, onClick }) => {
  const className = classnames({
    'iap-notification': true,
    'iap-notification--failure': type === types.FAILURE,
    'iap-notification--success': type === types.SUCCESS,
  });

  const Icon = type === 'failure' ? ErrorIcon : SuccessIcon;

  return (
    <button className={className} onClick={onClick}>
      <Icon className="iap-notification__icon" />
      <div>
        <p className="iap-notification__title">{title}</p>
        <p className="iap-notification__text">{text}</p>
      </div>
    </button>
  );
};

Notification.propTypes = {
  title: PropTypes.string,
  text: PropTypes.string.isRequired,
  type: PropTypes.oneOf(Object.values(types)).isRequired,
};

Notification.defaultProps = {
  title: 'Notification',
};

export default Notification;
