import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import { useDispatch, useSelector } from 'react-redux';
import StorageSettings from './StorageSettings';
import { fetchSettings, postSettings } from '~/modules/database/actions';
import DatabasePlaceholder from './Placeholder';
import { Tabs, Tab } from '~/components/shared/Tabs';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

import './index.css';

const StoragesSection = () => {
  const dispatch = useDispatch();
  const [tabIndex, setTabIndex] = useState(0);
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess } = useNotifications();
  const settings = useSelector(state => state.database.settings.config);
  const error = useSelector(state => state.database.settings.error);

  const triggerFetchSettings = async () => {
    setProgress(70);

    try {
      await dispatch(fetchSettings());
    } finally {
      setProgress(100);
    }
  };

  useEffect(() => {
    triggerFetchSettings();
  }, []);

  const handleSettingsSubmit = node => async (nodeSettings) => {
    setProgress(70);

    const updatedSettings = update(settings, {
      [node]: nodeSettings,
    });

    try {
      await dispatch(postSettings(updatedSettings));

      notifySuccess({
        title: 'Saved',
        text: 'Storage settings have been successfully saved',
      });
    } finally {
      setProgress(100);
    }
  };

  if (error) {
    return (
      <section className="iap-management-section">
        <DatabasePlaceholder
          fetching={progress}
          onTryAgainClick={triggerFetchSettings}
        />
      </section>
    );
  }

  const getStorageSettingsProps = (index) => {
    return [
      {
        title: 'Application Storage',
        description: 'Setup a connection to the database all the applications are stored in.',
        settings: settings ? settings.app_storage : null,
        postSettings: handleSettingsSubmit('app_storage'),
      },
      {
        title: 'User Storage',
        description: 'Setup a connection to the database all the users are stored in.',
        settings: settings ? settings.user_storage : null,
        postSettings: handleSettingsSubmit('user_storage'),
      },
      {
        title: 'Token Storage',
        description: 'Setup a connection to the database all the tokens are stored in.',
        settings: settings ? settings.token_storage : null,
        postSettings: handleSettingsSubmit('token_storage'),
      },
      {
        title: 'Verification Code Storage',
        description: 'Setup a connection to the database all the verification codes are stored in.',
        settings: settings ? settings.verification_code_storage : null,
        postSettings: handleSettingsSubmit('verification_code_storage'),
      },
      {
        title: 'Token Blacklist Storage',
        description: 'Setup a connection to the database all the blacklisted tokens are stored in.',
        settings: settings ? settings.token_blacklist : null,
        postSettings: handleSettingsSubmit('token_blacklist'),
      },
    ][index];
  };

  const storageSettingsProps = getStorageSettingsProps(tabIndex);

  return (
    <section className="iap-management-section">
      <header className="iap-management-section__header">
        <p className="iap-management-section__title">
          Storages
        </p>
      </header>

      <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
        <Tab title="Applications" />
        <Tab title="Users" />
        <Tab title="Tokens" />
        <Tab title="Verification Codes" />
        <Tab title="Blacklist" />

        <StorageSettings progress={!!progress} {...storageSettingsProps} />
      </Tabs>

    </section>
  );
};

export default StoragesSection;
