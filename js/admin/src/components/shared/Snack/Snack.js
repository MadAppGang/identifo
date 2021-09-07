import React from 'react';
import { createPortal } from 'react-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import Button from '~/components/shared/Button';
import './index.css';

const SnackComponent = ({ content, buttons, callback, success, error }) => {
  const snackClasses = classnames('iap-snack', {
    'iap-snack--success': success,
    'iap-snack--error': error,
  });
  return (
    <div className={snackClasses}>
      <div className="iap-snack--in">
        <div className="iap-snack--content"><span>{content}</span></div>
        <div className="iap-snack--controls">
          {buttons.map((btn, idx) => {
            return (
              <Button
                white
                outline={idx > 0}
                onClick={() => callback(btn.data)}
              >
                {btn.label}
              </Button>
            );
          })}
        </div>
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
  callback: PropTypes.func.isRequired,
  buttons: PropTypes.arrayOf(PropTypes.shape({
    label: PropTypes.string.isRequired,
    data: PropTypes.any.isRequired,
  })).isRequired,
};
