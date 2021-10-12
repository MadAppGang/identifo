import React, { useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { logout } from '~/modules/auth/actions';
import LogoutIcon from '~/components/icons/LogoutIcon.svg';

const LogoutSection = ({ onClick }) => {
  const dispatch = useDispatch();

  const handleClick = () => {
    dispatch(logout());
    onClick();
  };

  return (
    <button
      type="button"
      className="iap-profile-dropdown__section"
      onClick={useCallback(handleClick, [logout, onClick])}
    >
      <LogoutIcon className="iap-profile-dropdown__icon" fill="#6d6d6d" />
      <span>Logout</span>
    </button>
  );
};

LogoutSection.defaultProps = {
  onClick: () => {},
};

export default LogoutSection;
