import React from 'react';
import { Tab, Tabs } from '~/components/shared/Tabs';
import { useQuery } from '~/hooks/useQuery';
import { tabGroups } from '../../../enums';
import GeneralTab from './GeneralSection/ServerGeneralTab';
import { JWTSettingsSection } from './JWTSettings/JWTSettingsSection';
import { JWTStorageSection } from './JWTStorage/JWTStorageSection';
import ServerConfigurationForm from './ServerConfiguration/ServerConfigurationForm';

const renderComponent = (tabName) => {
  const components = {
    general: <GeneralTab />,
    token_storage: <JWTStorageSection />,
    token_settings: <JWTSettingsSection />,
    configuration_storage: <ServerConfigurationForm />,
    default: null,
  };
  return components[tabName] || components.default;
};

const GeneralSection = () => {
  const activeTab = useQuery().get(tabGroups.server_group);
  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Server Settings
      </p>
      <Tabs group={tabGroups.server_group}>
        <Tab title="General" />
        <Tab title="Token Storage" />
        <Tab title="Token Settings" />
        <Tab title="Configuration Storage" />
        <>
          {renderComponent(activeTab)}
        </>
      </Tabs>
    </section>
  );
};

export default GeneralSection;
