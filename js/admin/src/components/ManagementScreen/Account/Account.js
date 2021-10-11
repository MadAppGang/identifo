import React from 'react';
import { Tab, Tabs } from '~/components/shared/Tabs';
import { tabGroups } from '../../../enums';
import { useQuery } from '../../../hooks/useQuery';
import AccountSettings from './AdminAccountSettings';
import SessionStorageSettings from './SessionStorageSettings';

const renderComponents = (tab) => {
  const components = {
    admin_account: <AccountSettings />,
    session_storage: <SessionStorageSettings />,
    default: null,
  };
  return components[tab] || components.default;
};

const AccountSection = () => {
  const activeTab = useQuery().get(tabGroups.account_group);
  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Account Settings
      </p>

      <Tabs group={tabGroups.account_group}>
        <Tab title="Admin Account" />
        <Tab title="Session Storage" />
        <>
          {renderComponents(activeTab)}
        </>
      </Tabs>
    </section>
  );
};

export default AccountSection;
