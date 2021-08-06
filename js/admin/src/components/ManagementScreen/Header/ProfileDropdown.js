import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useDropdown from 'use-dropdown';
import DropdownIcon from '~/components/icons/DropdownIcon';
import LogoutSection from './LogoutSection';
import AccountSection from './AccountSection';
import { fetchAccountSettings } from '~/modules/account/actions';

const ProfileDropdown = () => {
  const [containerRef, isOpen, open, close] = useDropdown();
  const dispatch = useDispatch();
  const admin = useSelector(s => s.account.settings);

  useEffect(() => {
    dispatch(fetchAccountSettings());
  }, []);

  return (
    <div className="iap-header-profile" ref={containerRef}>
      <button
        type="button"
        className="iap-header-profile__trigger"
        onClick={open}
      >
        <span>{admin ? admin.email : 'Admin Panel'}</span>
        <DropdownIcon className="iap-dropdown-icon" />
      </button>
      {isOpen && (
        <div className="iap-profile-dropdown">
          <AccountSection onClick={close} />
          <LogoutSection onClick={close} />
        </div>
      )}
    </div>
  );
};

export default ProfileDropdown;
