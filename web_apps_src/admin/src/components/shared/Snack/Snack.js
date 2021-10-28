import classnames from 'classnames';
import PropTypes from 'prop-types';
import React, { useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useDispatch } from 'react-redux';
import Button from '~/components/shared/Button';
import { hideNotificationSnack } from '~/modules/applications/notification-actions';
import './index.css';
import { notificationStatuses } from '~/enums';

const SnackComponent = ({ content, buttons, callback, status }) => {
  const dispatch = useDispatch();
  const snackClasses = classnames('iap-snack', {
    'iap-snack--success': status === notificationStatuses.success,
    'iap-snack--error': status === notificationStatuses.error,
    'iap-snack--changes': status === notificationStatuses.changed,
  });

  useEffect(() => {
    if (status !== notificationStatuses.changed) {
      setTimeout(() => {
        dispatch(hideNotificationSnack());
      }, 5000);
    }
  }, []);

  return (
    <div className={snackClasses}>
      <div className="iap-snack--in">
        <div className="iap-snack--content"><span>{content}</span></div>
        {buttons
        && (
        <div className="iap-snack--controls">
          {buttons.map((btn, idx) => {
            return (
              <Button
                white
                outline={idx > 0}
                onClick={() => callback(btn.data)}
                key={btn.label}
              >
                {btn.label}
              </Button>
            );
          })}
        </div>
        )}
      </div>
    </div>
  );
};

export const Snack = (props) => {
  const Root = document.getElementById('iap-notifications');
  if (Root) {
    return createPortal(<SnackComponent {...props} />, Root);
  }
  return null;
};

Snack.propTypes = {
  content: PropTypes.string.isRequired,
  callback: PropTypes.func,
  buttons: PropTypes.arrayOf(PropTypes.shape({
    label: PropTypes.string.isRequired,
    data: PropTypes.any.isRequired,
  })),
  status: PropTypes.number.isRequired,
};

Snack.defaultProps = {
  buttons: undefined,
  callback: undefined,
};
