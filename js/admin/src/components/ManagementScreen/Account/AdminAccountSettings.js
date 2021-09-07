import React, { useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import update from '@madappgang/update-by-path';
import AccountForm from './AdminAccountForm';
import { fetchServerSetings, updateServerSettings } from '~/modules/settings/actions';
import SettingsPlaceholder from './Placeholder';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';
import { getAdminAccountSettings, getSessionStorageSettings } from '~/modules/settings/selectors';


const AdminAccountSettings = () => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const error = useSelector(s => s.account.error);
  const accountSettings = useSelector(getAdminAccountSettings);
  const sessionSettings = useSelector(getSessionStorageSettings);
  const settings = useMemo(() => {
    if (accountSettings && sessionSettings) {
      return { ...accountSettings, sessionDuration: sessionSettings.sessionDuration };
    }
    return null;
  }, [accountSettings, sessionSettings]);

  const fetchSettings = async () => {
    setProgress(70);
    await dispatch(fetchServerSetings());
    setProgress(100);
  };

  const handleFormSubmit = async (data) => {
    setProgress(70);
    try {
      const { sessionDuration, ...rest } = data;
      const payload = {
        adminAccount: update(accountSettings, rest),
        sessionStorage: update(sessionSettings, { sessionDuration }),
      };
      await dispatch(updateServerSettings(payload));
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
