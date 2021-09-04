import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { updateGeneralSettings } from '~/modules/settings/actions';
import ServerGeneralForm from './ServerGeneralForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const ServerGeneralTab = () => {
  const dispatch = useDispatch();
  const settings = useSelector(s => s.settings.general);
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess } = useNotifications();

  const handleSubmit = async (nextSettings) => {
    setProgress(70);
    await dispatch(updateGeneralSettings(nextSettings));
    setProgress(100);

    notifySuccess({
      title: 'Updated',
      text: 'Server settings have been updated successfully',
    });
  };

  return (
    <ServerGeneralForm
      loading={!!progress}
      settings={settings}
      onSubmit={handleSubmit}
    />
  );
};

export default ServerGeneralTab;
