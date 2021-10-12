import React from 'react';
import { DatagridRow } from '~/components/shared/Datagrid';

const ApplicationListPreloader = () => (
  <>
    <DatagridRow>
      <div className="iap-apps-row__fake-icon" />
      <div className="iap-apps-row__fake-type" />
      <div className="iap-apps-row__fake-client-id" />
    </DatagridRow>
    <DatagridRow>
      <div className="iap-apps-row__fake-icon" />
      <div className="iap-apps-row__fake-type" />
      <div className="iap-apps-row__fake-client-id" />
    </DatagridRow>
  </>
);

export default ApplicationListPreloader;
