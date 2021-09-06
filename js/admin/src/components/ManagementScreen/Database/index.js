import update from '@madappgang/update-by-path';
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Tab, Tabs } from '~/components/shared/Tabs';
import useNotifications from '~/hooks/useNotifications';
import useProgressBar from '~/hooks/useProgressBar';
import { settingsActionsEnum, settingsConfig } from '~/modules/applications/dialogsConfigs';
import { postSettings, verifyConnection } from '~/modules/database/actions';
import { CONNECTION_SUCCEED } from '~/modules/database/connectionReducer';
import { fetchServerSetings } from '~/modules/settings/actions';
import { handleSettingsDialog, hideSettingsDialog } from '../../../modules/applications/actions';
import './index.css';
import DatabasePlaceholder from './Placeholder';
import StorageSettings from './StorageSettings';
import { getStorageSettings } from '~/modules/settings/selectors';

const StoragesSection = () => {
  const dispatch = useDispatch();
  const [tabIndex, setTabIndex] = useState(0);
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess } = useNotifications();
  const settings = useSelector(getStorageSettings);
  const error = useSelector(state => state.database.settings.error);
  const connectionState = useSelector(state => state.database.connection.state);

  const triggerFetchSettings = async () => {
    setProgress(70);

    try {
      await dispatch(fetchServerSetings());
    } finally {
      setProgress(100);
    }
  };

  const saveHandler = async (node, nodeSettings) => {
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

  const handleSettingsVerification = async (nodeSettings) => {
    setProgress(70);
    try {
      await dispatch(verifyConnection(nodeSettings));
    } finally {
      setProgress(100);
    }
  };

  const handleSettingsSubmit = node => async (nodeSettings) => {
    if (connectionState !== CONNECTION_SUCCEED) {
      const config = {
        ...settingsConfig[connectionState],
        onClose: () => dispatch(hideSettingsDialog()),
      };
      const res = await dispatch(handleSettingsDialog(config));
      switch (res) {
        case settingsActionsEnum.save:
          await saveHandler(node, nodeSettings);
          break;
        case settingsActionsEnum.verify:
          await handleSettingsVerification(nodeSettings);
          await saveHandler(node, nodeSettings);
          break;
        case settingsActionsEnum.close:
          dispatch(hideSettingsDialog());
          break;
        default:
          dispatch(hideSettingsDialog());
      }
    } else {
      await saveHandler(node, nodeSettings);
    }
  };

  useEffect(() => {
    setProgress(100);
  }, []);

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
        settings: settings ? settings.appStorage : null,
        postSettings: handleSettingsSubmit('app_storage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'User Storage',
        description: 'Setup a connection to the database all the users are stored in.',
        settings: settings ? settings.userStorage : null,
        postSettings: handleSettingsSubmit('user_storage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Token Storage',
        description: 'Setup a connection to the database all the tokens are stored in.',
        settings: settings ? settings.tokenStorage : null,
        postSettings: handleSettingsSubmit('token_storage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Verification Code Storage',
        description: 'Setup a connection to the database all the verification codes are stored in.',
        settings: settings ? settings.verificationCodeStorage : null,
        postSettings: handleSettingsSubmit('verification_code_storage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Token Blacklist Storage',
        description: 'Setup a connection to the database all the blacklisted tokens are stored in.',
        settings: settings ? settings.tokenBlacklist : null,
        postSettings: handleSettingsSubmit('token_blacklist'),
        verifySettings: handleSettingsVerification,
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
        <StorageSettings
          connectionState={connectionState}
          progress={!!progress}
          {...storageSettingsProps}
        />
      </Tabs>
    </section>
  );
};

export default StoragesSection;
