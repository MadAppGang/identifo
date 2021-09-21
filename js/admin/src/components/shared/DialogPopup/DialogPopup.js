import React, { useRef } from 'react';
import { createPortal } from 'react-dom';
import Button from '~/components/shared/Button';
import PropTypes from 'prop-types';
import './index.css';

const Dialog = ({ title, content, buttons, children, callback, onClose }) => {
  const node = useRef();
  const onOverlayClick = (e) => {
    if (node.current && e.target === node.current) {
      onClose();
    }
  };
  return (
    <div className="iap-dialog-popup--overlay" onClick={onOverlayClick} ref={node} role="presentation">
      <div className="iap-dialog-popup">
        <div className="iap-dialog-popup--in">
          {title && <h3 className="iap-dialog-popup--title">{title}</h3>}
          <div className="iap-dialog-popup--content">
            <span>{content}</span>
            <div className="iap-dialog-popup--content-children">{children}</div>
          </div>
          <div className="iap-dialog-popup--controls">
            {buttons.map((btn) => {
              return (
                <Button
                  key={btn.label}
                  onClick={() => callback(btn.data)}
                  error={!!btn.error}
                  outline={!!btn.outline}
                  disabled={!!btn.disabled}
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
