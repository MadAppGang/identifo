import React, { useRef } from 'react';
import { createPortal } from 'react-dom';
import Button from '~/components/shared/Button';
import PropTypes from 'prop-types';
import './index.css';
import classNames from 'classnames';
import { dialogTypes } from '~/modules/applications/dialogsConfigs';

const Dialog = ({ title, content, buttons, children, type, callback, onClose }) => {
  const node = useRef();
  const popupClass = classNames('iap-dialog-popup', { 'iap-dialog-popup__danger': type === dialogTypes.danger });
  const onOverlayClick = (e) => {
    if (node.current && e.target === node.current) {
      onClose();
    }
  };

  return (
    <div className="iap-dialog-popup--overlay" onClick={onOverlayClick} ref={node} role="presentation">
      <div className={popupClass}>
        <div className="iap-dialog-popup--in">
          {title && <h3 className="iap-dialog-popup--title">{title}</h3>}
          <div className="iap-dialog-popup--content">
            <span>{content}</span>
            {children && <div className="iap-dialog-popup--content-children">{children}</div>}
          </div>
          <div className="iap-dialog-popup--controls">
            {buttons.map(({ label, data, ...btnTypes }) => {
              return (
                <Button
                  key={label}
                  onClick={() => callback(data)}
                  {...btnTypes}
                >
                  {label}
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
  type: PropTypes.string,
  buttons: PropTypes.arrayOf(PropTypes.shape({
    label: PropTypes.string.isRequired,
    data: PropTypes.any.isRequired,
  })).isRequired,
};

DialogPopup.defaultProps = {
  title: '',
  type: 'default',
};
