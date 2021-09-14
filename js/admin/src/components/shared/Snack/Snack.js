import React, { useEffect } from 'react';
import { createPortal } from 'react-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import Button from '~/components/shared/Button';
import './index.css';
import { useDispatch } from 'react-redux';
import { hideNotificationSnack } from '../../../modules/applications/actions';

const SnackComponent = ({ content, buttons, callback, status }) => {
  const dispatch = useDispatch();
  const snackClasses = classnames('iap-snack', {
    'iap-snack--success': status === 'success',
    'iap-snack--error': status === 'error' || status === 'rejected',
  });

  useEffect(() => {
    if (status !== 'changed') {
      setTimeout(() => {
        dispatch(hideNotificationSnack());
      }, 3000);
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
  const Root = document.getElementById('root');
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
  status: PropTypes.string,
};

Snack.defaultProps = {
  buttons: undefined,
  callback: undefined,
  status: '',
};
