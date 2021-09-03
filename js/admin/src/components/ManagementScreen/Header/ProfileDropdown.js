import React from 'react';
import { useSelector } from 'react-redux';
import useDropdown from 'use-dropdown';
import DropdownIcon from '~/components/icons/DropdownIcon';
import AccountSection from './AccountSection';
import LogoutSection from './LogoutSection';

const ProfileDropdown = () => {
  const [containerRef, isOpen, open, close] = useDropdown();
  const admin = useSelector(s => s.settings.adminAccount);

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
