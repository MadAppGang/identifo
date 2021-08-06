import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import AccountForm from './AdminAccountForm';
import { fetchAccountSettings, postAccountSettings } from '~/modules/account/actions';
import SettingsPlaceholder from './Placeholder';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const AdminAccountSettings = () => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const error = useSelector(s => s.account.error);
  const settings = useSelector(s => s.account.settings);

  const fetchSettings = async () => {
    setProgress(70);
    await dispatch(fetchAccountSettings());
    setProgress(100);
  };

  React.useEffect(() => {
    fetchSettings();
  }, []);

  const handleFormSubmit = async () => {
    setProgress(70);
    try {
      await dispatch(postAccountSettings(settings));
      notifySuccess({
        title: 'Saved',
        text: 'Account settings have been successfully saved',
      });
    } catch (_) {
      notifyFailure({
        title: 'Error',
        text: 'Account settings could not be saved',
      });
    } finally {
      setProgress(100);
    }
  };

  if (error) {
    return (
      <SettingsPlaceholder
        fetching={!!progress}
        onTryAgainClick={fetchSettings}
      />
    );
  }

  return (
    <AccountForm
      error={error}
      loading={!!progress}
      settings={settings}
      onSubmit={handleFormSubmit}
    />
  );
};

export default AdminAccountSettings;
