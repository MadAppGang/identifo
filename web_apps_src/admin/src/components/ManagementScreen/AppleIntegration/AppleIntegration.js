import React from 'react';
import { Tab, Tabs } from '~/components/shared/Tabs';
import { tabGroups } from '../../../enums';
import { useQuery } from '../../../hooks/useQuery';
import AppSiteAssociationForm from './AppSiteAssociationForm';
import DomainAssociationForm from './DomainAssociationForm';

const renderComponent = (tab) => {
  const components = {
    app_site_association: <AppSiteAssociationForm />,
    developer_domain_association: <DomainAssociationForm />,
    default: null,
  };
  return components[tab] || components.default;
};

const AppleIntegration = () => {
  const activeTab = useQuery().get(tabGroups.apple_integration_group);
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
        <Tabs group={tabGroups.apple_integration_group}>
          <Tab title="App Site Association" />
          <Tab title="Developer Domain Association" />
          <>
            {renderComponent(activeTab)}
          </>
        </Tabs>
      </main>
    </section>
  );
};

export default AppleIntegration;
