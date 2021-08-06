import React from 'react';
import { Switch, Route } from 'react-router-dom';
import Header from './Header';
import Sidebar from './Sidebar';
import UsersSection from './Users';
import DatabaseSection from './Database';
import ApplicationsSection from './Applications';
import AccountSection from './Account';
import ServerSection from './Server';
import ExternalServicesSection from './ExternalServices';
import LoginTypesSection from './LoginTypes';
import MultiFactorAuthSection from './MultiFactorAuth';
import HostedPagesSection from './HostedPages';
import NotFoundSection from './NotFoundSection';
import StaticFilesSection from './StaticFiles';
import AppleIntegrationSection from './AppleIntegration';
import ReloadServerPopup from './ReloadServerPopup';
import Container from '~/components/shared/Container';
import { NotificationContainer } from '~/components/shared/Notifications';
import ProgressBar from '~/components/shared/TopProgressBar';
import './Management.css';

const ManagementScreen = () => (
  <div className="iap-management-layout">
    <ProgressBar>
      <NotificationContainer>
        <Header />
        <div className="iap-management-content">
          <Container>
            <Sidebar />
            <Switch>
              <Route exact path="/management" component={ServerSection} />
              <Route path="/management/users" component={UsersSection} />
              <Route path="/management/database" component={DatabaseSection} />
              <Route path="/management/applications" component={ApplicationsSection} />
              <Route path="/management/email_integration" component={ExternalServicesSection} />
              <Route path="/management/account" component={AccountSection} />
              <Route path="/management/settings" component={LoginTypesSection} />
              <Route path="/management/multi-factor_auth" component={MultiFactorAuthSection} />
              <Route path="/management/static" component={StaticFilesSection} />
              <Route path="/management/hosted_pages" component={HostedPagesSection} />
              <Route path="/management/apple" component={AppleIntegrationSection} />
              <Route component={NotFoundSection} />
            </Switch>
          </Container>
        </div>
        <ReloadServerPopup />
      </NotificationContainer>
    </ProgressBar>
  </div>
);

export default ManagementScreen;
