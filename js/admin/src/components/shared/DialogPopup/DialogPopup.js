import React from 'react';
import { createPortal } from 'react-dom';
import Button from '~/components/shared/Button';
import PropTypes from 'prop-types';
import './index.css';

export const Dialog = ({ title, content, onSubmit, onCancel }) => {
  return (
    <div className="iap-dialog-popup--overlay">
      <div className="iap-dialog-popup">
        <div className="iap-dialog-popup--in">
          {title && <h3 className="iap-dialog-popup--title">{title}</h3>}
          <div className="iap-dialog-popup--content">{content}</div>
          <div className="iap-dialog-popup--controls">
            <Button onClick={onSubmit}>Confirm</Button>
            <Button onClick={onCancel} transparent>Cancel</Button>
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
  onSubmit: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
};

DialogPopup.defaultProps = {
  title: '',
};
