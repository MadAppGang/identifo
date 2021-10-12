import React, { useState } from 'react';
import { Tabs, Tab } from '~/components/shared/Tabs';
import AppSiteAssociationForm from './AppSiteAssociationForm';
import DomainAssociationForm from './DomainAssociationForm';

const AppleIntegration = () => {
  const [tabIndex, setTabIndex] = useState(0);

  return (
    <section className="iap-management-section">
      <header>
        <p className="iap-management-section__title">
          Apple Integration
        </p>

        <p className="iap-management-section__description">
          Sign In with Apple configuration for users to sign in using their Apple ID
        </p>
      </header>

      <main className="iap-settings-section">
        <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
          <Tab title="App Site Association" />
          <Tab title="Developer Domain Association" />

          <>
            {tabIndex === 0 && (
              <AppSiteAssociationForm />
            )}

            {tabIndex === 1 && (
              <DomainAssociationForm />
            )}
          </>
        </Tabs>
      </main>
    </section>
  );
};

export default AppleIntegration;
