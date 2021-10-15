import React, { useState } from 'react';
import { Tab, Tabs } from '~/components/shared/Tabs';
import GeneralTab from './GeneralSection/ServerGeneralTab';
import { JWTSettingsSection } from './JWTSettings/JWTSettingsSection';
import { JWTStorageSection } from './JWTStorage/JWTStorageSection';
import ServerConfigurationForm from './ServerConfiguration/ServerConfigurationForm';

const GeneralSection = () => {
  const [tabIndex, setTabIndex] = useState(0);

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Server Settings
      </p>
      <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
        <Tab title="General" />
        <Tab title="Token Storage" />
        <Tab title="Token Settings" />
        <Tab title="Configuration Storage" />
        <>
          {tabIndex === 0 && <GeneralTab />}
          {tabIndex === 1 && <JWTStorageSection />}
          {tabIndex === 2 && <JWTSettingsSection />}
          {tabIndex === 3 && <ServerConfigurationForm />}
        </>
      </Tabs>
    </section>
  );
};

export default GeneralSection;
