import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useNotifications from '~/hooks/useNotifications';
import useProgressBar from '~/hooks/useProgressBar';
import { updateServerSettings } from '~/modules/settings/actions';
import { getGeneralServerSettings } from '~/modules/settings/selectors';
import ServerGeneralForm from './ServerGeneralForm';

const ServerGeneralTab = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getGeneralServerSettings);
  const { progress } = useProgressBar();
  const { notifySuccess } = useNotifications();

  const handleSubmit = async (nextSettings) => {
    const payload = { general: { ...nextSettings } };
    dispatch(updateServerSettings(payload));
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
