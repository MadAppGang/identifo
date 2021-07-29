import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link, Redirect } from 'react-router-dom';
import {
  alterApplication,
  deleteApplicationById,
  fetchApplicationById,
  resetApplicationError,
  fetchFederatedProviders,
} from '~/modules/applications/actions';
import ActionsButton from '~/components/shared/ActionsButton';
import { Tabs, Tab } from '~/components/shared/Tabs';
import ApplicationGeneralSettings from './GeneralSettingsForm';
import ApplicationAuthSettings from './AuthSettingsForm';
import ApplicationTokenSettings from './TokenSettingsForm';
import ApplicationFederatedLoginSettings from './FederatedLoginSettingsForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const goBackPath = '/management/applications';

const EditApplicationView = ({ match, history }) => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const id = match.params.appid;

  const application = useSelector(s => s.selectedApplication.application);
  const federatedProviders = useSelector(s => s.applications.federatedProviders);
  const error = useSelector(s => s.selectedApplication.error);
  const [tabIndex, setTabIndex] = React.useState(0);

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
      notifySuccess({
        title: 'Updated',
        text: 'Application has been updated successfully',
      });
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'Application could not be updated',
      });
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
      notifySuccess({
        title: 'Deleted',
        text: 'Application has been deleted successfully',
      });
      history.push(goBackPath);
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'Application could not be deleted',
      });
    } finally {
      setProgress(100);
    }
  };

  const availableActions = React.useMemo(() => [{
    title: 'Delete Application',
    onClick: handleDeleteClick,
  }], [id, handleDeleteClick]);

  if (progress > 70 && !application) {
    notifyFailure({
      title: 'Something went wrong',
      text: `Can't find application with id ${id}`,
    });

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
          <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
            <Tab title="General" />
            <Tab title="Authorization" />
            <Tab title="Tokens" />
            <Tab title="Federated Login" />

            <>
              {tabIndex === 0 && (
                <ApplicationGeneralSettings
                  error={error}
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}

              {tabIndex === 1 && (
                <ApplicationAuthSettings
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}

              {tabIndex === 2 && (
                <ApplicationTokenSettings
                  loading={!!progress}
                  application={application}
                  onCancel={handleCancel}
                  onSubmit={handleSubmit}
                />
              )}

              {tabIndex === 3 && (
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
