import React, { useState } from 'react';
import SessionStorageSettings from './SessionStorageSettings';
import AccountSettings from './AdminAccountSettings';
import { Tabs, Tab } from '~/components/shared/Tabs';

const AccountSection = () => {
  const [tabIndex, setTabIndex] = useState(0);

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Account Settings
      </p>

      <Tabs
        activeTabIndex={tabIndex}
        onChange={setTabIndex}
      >
        <Tab title="Admin Account" />
        <Tab title="Session Storage" />

        <>
          {tabIndex === 0 && <AccountSettings />}
          {tabIndex === 1 && <SessionStorageSettings />}
        </>
      </Tabs>
    </section>
  );
};

export default AccountSection;
