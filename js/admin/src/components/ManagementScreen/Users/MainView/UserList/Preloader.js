import React from 'react';
import { DatagridRow } from '~/components/shared/Datagrid';

const Preloader = () => (
  <>
    <DatagridRow>
      <div className="iap-users-row__fake-icon" />
      <div className="iap-users-row__fake-name" />
      <div className="iap-users-row__fake-email" />
      <div className="iap-users-row__fake-latest-login" />
      <div className="iap-users-row__fake-logins" />
    </DatagridRow>
    <DatagridRow>
      <div className="iap-users-row__fake-icon" />
      <div className="iap-users-row__fake-name" />
      <div className="iap-users-row__fake-email" />
      <div className="iap-users-row__fake-latest-login" />
      <div className="iap-users-row__fake-logins" />
    </DatagridRow>
  </>
);

export default Preloader;
