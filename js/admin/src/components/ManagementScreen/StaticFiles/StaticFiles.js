import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import StaticFilesGeneralForm from './StaticFilesGeneralForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';
import { getStaticFilesSettings } from '~/modules/settings/selectors';
import { updateServerSettings } from '../../../modules/settings/actions';

const StaticFilesSection = () => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess } = useNotifications();
  const settings = useSelector(getStaticFilesSettings);

  const handleSubmit = async (nextSettings) => {
    setProgress(70);
    const payload = { staticFilesStorage: nextSettings };
    await dispatch(updateServerSettings(payload));
    setProgress(100);

    notifySuccess({
      title: 'Updated',
      text: 'Settings have been updated successfully',
    });
  };

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Static Files
      </p>

      <p className="iap-management-section__description">
        These settings allow to specify paths to various static files directories.
      </p>

      <StaticFilesGeneralForm
        settings={settings || {}}
        loading={!!progress}
        onSubmit={handleSubmit}
      />
    </section>
  );
};

export default StaticFilesSection;
