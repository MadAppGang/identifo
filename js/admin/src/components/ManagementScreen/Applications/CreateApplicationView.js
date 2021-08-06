import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import { postApplication, resetApplicationError } from '~/modules/applications/actions';
import ApplicationGeneralSettings from './GeneralSettingsForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const goBackPath = '/management/applications';

const CreateApplicationView = ({ history }) => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const error = useSelector(s => s.selectedApplication.error);
  const application = useSelector(s => s.selectedApplication.application);

  React.useEffect(() => {
    if (application && application.id && progress === 100) {
      history.push(`/management/applications/${application.id}`);
    }
  }, [application, progress]);

  const handleSubmit = async (data) => {
    setProgress(70);

    try {
      await dispatch(postApplication(data));

      notifySuccess({
        title: 'Created',
        text: 'Application has been created successfully',
      });
    } catch (_) {
      notifyFailure({
        title: 'Error',
        text: 'Application could not be created',
      });
    } finally {
      setProgress(100);
    }
  };

  const handleCancel = () => {
    dispatch(resetApplicationError());
    history.push(goBackPath);
  };

  return (
    <section className="iap-management-section">
      <header>
        <div>
          <Link to={goBackPath} className="iap-management-section__back">
            ← &nbsp;Applications
          </Link>
        </div>
        <p className="iap-management-section__title">
          Create Application
        </p>
        <p className="iap-management-section__description">
          Configure allowed callback URLs and Secrets for your application.
        </p>
      </header>
      <main>
        <ApplicationGeneralSettings
          error={error}
          loading={!!progress}
          excludeFields={[
            'secret', 'active', 'tfaStatus', 'redirectUrl', 'allowRegistration', 'debugTfaCode',
          ]}
          onCancel={handleCancel}
          onSubmit={handleSubmit}
        />
      </main>
    </section>
  );
};

export default CreateApplicationView;
