import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import SessionStorageForm from './SessionStorageForm';
import {
  fetchSessionStorageSettings, updateSessionStorageSettings,
} from '~/modules/settings/actions';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const SessionStorageSettings = ({ error }) => {
  const dispatch = useDispatch();
  const { notifySuccess } = useNotifications();
  const { progress, setProgress } = useProgressBar();
  const settings = useSelector(state => state.settings.sessionStorage);

  const fetchSettings = async () => {
    setProgress(70);
    await dispatch(fetchSessionStorageSettings());
    setProgress(100);
  };

  React.useEffect(() => {
    fetchSettings();
  }, []);

  const handleSubmit = async (data) => {
    setProgress(70);
    await dispatch(updateSessionStorageSettings(data));
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
