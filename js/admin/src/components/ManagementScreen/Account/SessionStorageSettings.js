import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import SessionStorageForm from './SessionStorageForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const SessionStorageSettings = ({ error }) => {
  const dispatch = useDispatch();
  const { notifySuccess } = useNotifications();
  const { progress, setProgress } = useProgressBar();
  const settings = useSelector(state => state.settings.sessionStorage);

  const handleSubmit = async (data) => {
    setProgress(70);
    // TODO: Nikita k update settings
    setProgress(100);

    notifySuccess({
      title: 'Updated',
      text: 'Settings have been updated successfully',
    });
  };

  return (
    <SessionStorageForm
      error={error}
      loading={!!progress}
      settings={settings}
      onSubmit={handleSubmit}
    />
  );
};

export default SessionStorageSettings;
