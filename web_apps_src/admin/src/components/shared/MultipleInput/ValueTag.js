import React from 'react';
import CancelIcon from '~/components/icons/CancelIcon.svg';

const ValueTag = (props) => {
  const { onDeleteClick, children } = props;

  return (
    <li className="value-tag">
      {children}
      <button
        type="button"
        onClick={onDeleteClick}
        className="value-tag__delete"
      >
        <CancelIcon className="value-tag__delete-icon" />
      </button>
    </li>
  );
};

export default ValueTag;
