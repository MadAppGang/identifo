import React from 'react';
import { useSelector } from 'react-redux';
import { Route, Switch } from 'react-router-dom';
import Container from '~/components/shared/Container';
import { DialogPopup } from '~/components/shared/DialogPopup/DialogPopup';
import { NotificationContainer } from '~/components/shared/Notifications';
import ProgressBar from '~/components/shared/TopProgressBar';
import { useSettings } from '../../hooks/useSettings';
import AccountSection from './Account';
import AppleIntegrationSection from './AppleIntegration';
import ApplicationsSection from './Applications';
import DatabaseSection from './Database';
import ExternalServicesSection from './ExternalServices';
import Header from './Header';
import HostedPagesSection from './HostedPages';
import LoginTypesSection from './LoginTypes';
import './Management.css';
import MultiFactorAuthSection from './MultiFactorAuth';
import NotFoundSection from './NotFoundSection';
import ReloadServerPopup from './ReloadServerPopup';
import ServerSection from './Server';
import Sidebar from './Sidebar';
import StaticFilesSection from './StaticFiles';
import UsersSection from './Users';
import { SaveChangesSnack } from '~/components/shared/SaveChangesSnack/SaveChangesSnack';
import { Snack } from '~/components/shared/Snack/Snack';
import { notificationStatuses } from '~/enums';

const ManagementScreen = () => {
  const dialogConfig = useSelector(s => s.notifications.settingsDialog);
  const notificationSnack = useSelector(s => s.notifications.notificationSnack);
  useSettings();

  return (
    <div className="iap-management-layout">
      {/* iap-notifications is needed to render portal snack */}
      <div id="iap-notifications" className="iap-notifications" />
      {dialogConfig.show && <DialogPopup {...dialogConfig.config} />}
      {notificationSnack.status !== notificationStatuses.idle
        && <Snack content={notificationSnack.message} status={notificationSnack.status} />}
      <ProgressBar>
        <SaveChangesSnack />
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
};

export default ManagementScreen;
