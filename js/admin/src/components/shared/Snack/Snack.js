import React from 'react';
import { createPortal } from 'react-dom';
import PropTypes from 'prop-types';
import './index.css';

const SnackComponent = ({ content, buttons, callback }) => {
  return (
    <div className="iap-snack">
      <div className="iap-snack--in">
        <div className="iap-snack--content">{content}</div>
        <div className="iap-snack--controls">
          {buttons.map((btn) => {
            return (
              <button className="iap-snack--control-item" key={btn.label} onClick={() => callback(btn.data)}>
                {btn.label}
              </button>
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
