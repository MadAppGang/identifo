import React from 'react';
import { Link } from 'react-router-dom';
import Container from '~/components/shared/Container';
import ProfileDropdown from './ProfileDropdown';
import './Header.css';

const ManagementScreenHeader = () => (
  <header className="iap-management-header">
    <Container>
      <div className="iap-management-header__inner">
        <Link to="/management" className="rrdl">
          <h1 className="iap-management-header__logo">
            <span>identifo</span>
            <span>Admin Panel</span>
          </h1>
        </Link>
        <ProfileDropdown />
      </div>
    </Container>
  </header>
);

export default ManagementScreenHeader;
