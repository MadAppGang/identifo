import React from 'react';
import { createPortal } from 'react-dom';
import Button from '~/components/shared/Button';
import PropTypes from 'prop-types';
import './index.css';

export const Dialog = ({ title, content, buttons, callback }) => {
  return (
    <div className="iap-dialog-popup--overlay">
      <div className="iap-dialog-popup">
        <div className="iap-dialog-popup--in">
          {title && <h3 className="iap-dialog-popup--title">{title}</h3>}
          <div className="iap-dialog-popup--content">{content}</div>
          <div className="iap-dialog-popup--controls">
            {buttons.map((btn) => {
              return (
                <Button
                  key={btn.label}
                  onClick={() => callback(btn.data)}
                >
                  {btn.label}
                </Button>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export const DialogPopup = (props) => {
  const Root = document.getElementById('root');
  if (Root) {
    return createPortal(<Dialog {...props} />, Root);
  }
  return null;
};

DialogPopup.propTypes = {
  title: PropTypes.string,
  content: PropTypes.string.isRequired,
  callback: PropTypes.func.isRequired,
  buttons: PropTypes.arrayOf(PropTypes.shape({
    label: PropTypes.string.isRequired,
    data: PropTypes.any.isRequired,
  })).isRequired,
};

DialogPopup.defaultProps = {
  title: '',
};
