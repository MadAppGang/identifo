import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link, Redirect } from 'react-router-dom';
import ActionsButton from '~/components/shared/ActionsButton';
import { Tab, Tabs } from '~/components/shared/Tabs';
import useProgressBar from '~/hooks/useProgressBar';
import {
  alterApplication,
  deleteApplicationById,
  fetchApplicationById, fetchFederatedProviders, resetApplicationError,
} from '~/modules/applications/actions';
import { errorSnackMessages } from '~/modules/applications/constants';
import { showErrorNotificationSnack } from '~/modules/applications/notification-actions';
import { tabGroups } from '../../../enums';
import { useQuery } from '../../../hooks/useQuery';
import ApplicationAuthSettings from './AuthSettingsForm';
import ApplicationFederatedLoginSettings from './FederatedLoginSettingsForm';
import ApplicationGeneralSettings from './GeneralSettingsForm';
import { RegistrationSettingsForm } from './RegistrationSettingsForm';
import ApplicationTokenSettings from './TokenSettingsForm';

const goBackPath = '/management/applications';

const tabsTitles = {
  general: 'General',
  registration: 'Registration',
  authorization: 'Authorization',
  tokens: 'Tokens',
  federated_login: 'Federated Login',
};

const tabsMatcher = Object.values(tabsTitles).reduce((p, n) => {
  // eslint-disable-next-line no-param-reassign
  p[n] = n.toLowerCase().replaceAll(' ', '_');
  return p;
}, {});

const EditApplicationView = ({ match, history }) => {
  const activeTab = useQuery().get(tabGroups.edit_app_group);
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const id = match.params.appid;

  const application = useSelector(s => s.selectedApplication.application);
  const federatedProviders = useSelector(s => s.applications.federatedProviders);
  const error = useSelector(s => s.selectedApplication.error);

  const fetchData = async () => {
    setProgress(70);
    await dispatch(fetchApplicationById(id));
    await dispatch(fetchFederatedProviders());
    setProgress(100);
  };

  React.useEffect(() => {
    fetchData();
  }, []);

  const handleSubmit = async (data) => {
    setProgress(70);
    try {
      await dispatch(alterApplication(id, data));
    } finally {
      setProgress(100);
    }
  };

  const handleCancel = () => {
    dispatch(resetApplicationError());
    history.push(goBackPath);
  };

  const handleDeleteClick = async () => {
    setProgress(70);
    try {
      await dispatch(deleteApplicationById(id));
      history.push(goBackPath);
    } finally {
      setProgress(100);
    }
  };

  const availableActions = React.useMemo(() => [{
    title: 'Delete Application',
    onClick: handleDeleteClick,
  }], [id, handleDeleteClick]);

  if (progress > 70 && !application) {
    dispatch(showErrorNotificationSnack(`${errorSnackMessages.appNotFound} ${id}`));
    dispatch(resetApplicationError());

    return <Redirect to={goBackPath} />;
  }

  return (
    <section className="iap-management-section">
      <header>
        <div>
          <Link to={goBackPath} className="iap-management-section__back">
            ← &nbsp;Applications
          </Link>
        </div>
        <div className="iap-management-section__title">
          Application Details

          <ActionsButton loading={!!progress} actions={availableActions} />
        </div>
        <p className="iap-management-section__description">
          <span className="iap-section-description__id">
            Client ID:&nbsp;
            {id}
          </span>
        </p>
      </header>
      <main>
        <div className="iap-management-section__tabs">
          <Tabs group={tabGroups.edit_app_group}>
            <Tab title={tabsTitles.general} />
            <Tab title={tabsTitles.registration} />
            <Tab title={tabsTitles.authorization} />
            <Tab title={tabsTitles.tokens} />
            <Tab title={tabsTitles.federated_login} />

            <>
              {activeTab === tabsMatcher[tabsTitles.general] && (
                <ApplicationGeneralSettings
                  error={error}
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                  excludeFields={['newUserDefaultRole', 'newUserDefaultScopes', 'allowRegistration', 'allowAnonymousRegistration']}
                />
              )}
              {activeTab === tabsMatcher[tabsTitles.registration] && (
              <RegistrationSettingsForm
                error={error}
                loading={!!progress}
                application={application}
                onCancel={handleCancel}
                onSubmit={handleSubmit}
              />
              )}
              {activeTab === tabsMatcher[tabsTitles.authorization] && (
                <ApplicationAuthSettings
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}

              {activeTab === tabsMatcher[tabsTitles.tokens] && (
                <ApplicationTokenSettings
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}

              {activeTab === tabsMatcher[tabsTitles.federated_login] && (
                <ApplicationFederatedLoginSettings
                  federatedProviders={federatedProviders}
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}
            </>
          </Tabs>
        </div>
      </main>
    </section>
  );
};

export default EditApplicationView;
