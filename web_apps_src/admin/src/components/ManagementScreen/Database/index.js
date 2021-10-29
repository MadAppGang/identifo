import update from '@madappgang/update-by-path';
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Tab, Tabs } from '~/components/shared/Tabs';
import { tabGroups, verificationStatuses } from '~/enums';
import useProgressBar from '~/hooks/useProgressBar';
import { useVerification } from '~/hooks/useVerification';
import { dialogActions, settingsConfig } from '~/modules/applications/dialogsConfigs';
import { fetchServerSetings, updateServerSettings } from '~/modules/settings/actions';
import { getStorageSettings } from '~/modules/settings/selectors';
import { useQuery } from '../../../hooks/useQuery';
import { handleSettingsDialog, hideSettingsDialog } from '../../../modules/applications/actions';
import './index.css';
import DatabasePlaceholder from './Placeholder';
import StorageSettings from './StorageSettings';

const tabsTitles = {
  applications: 'Applications',
  users: 'Users',
  tokens: 'Tokens',
  verifications_codes: 'Verification Codes',
  blacklist: 'Blacklist',
};

const getTabIndex = (urlTab) => {
  const tabsUrls = Object.keys(tabsTitles);
  const idx = tabsUrls.indexOf(urlTab);
  return idx === -1 ? 0 : idx;
};

const StoragesSection = () => {
  const dispatch = useDispatch();
  const activeTab = useQuery().get(tabGroups.storages_group);
  const [verificationStatus, verify, setStatus] = useVerification();
  const { progress, setProgress } = useProgressBar();
  const settings = useSelector(getStorageSettings);
  const error = useSelector(state => state.database.settings.error);

  const triggerFetchSettings = async () => {
    setProgress(70);
    try {
      await dispatch(fetchServerSetings());
    } finally {
      setProgress(100);
    }
  };

  const saveHandler = async (node, nodeSettings) => {
    const updatedSettings = { storage: update(settings, {
      [node]: nodeSettings,
    }) };
    dispatch(updateServerSettings(updatedSettings));
  };

  const handleSettingsVerification = async (nodeSettings) => {
    setProgress(70);
    try {
      const paylaod = { type: 'database', database: nodeSettings };
      await dispatch(verify(paylaod));
    } finally {
      setProgress(100);
    }
  };

  const handleSettingsSubmit = node => async (nodeSettings) => {
    if (verificationStatus !== verificationStatuses.success) {
      const config = {
        ...settingsConfig[verificationStatus],
        onClose: () => dispatch(hideSettingsDialog()),
      };
      const res = await dispatch(handleSettingsDialog(config));
      switch (res) {
        case dialogActions.submit:
          await saveHandler(node, nodeSettings);
          break;
        case dialogActions.verify:
          await handleSettingsVerification(nodeSettings);
          await saveHandler(node, nodeSettings);
          break;
        default:
          dispatch(hideSettingsDialog());
      }
    } else {
      await saveHandler(node, nodeSettings);
    }
  };

  useEffect(() => {
    setStatus(verificationStatuses.required);
  }, [getTabIndex(activeTab)]);

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
        postSettings: handleSettingsSubmit('appStorage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'User Storage',
        description: 'Setup a connection to the database all the users are stored in.',
        settings: settings ? settings.userStorage : null,
        postSettings: handleSettingsSubmit('userStorage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Token Storage',
        description: 'Setup a connection to the database all the tokens are stored in.',
        settings: settings ? settings.tokenStorage : null,
        postSettings: handleSettingsSubmit('tokenStorage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Verification Code Storage',
        description: 'Setup a connection to the database all the verification codes are stored in.',
        settings: settings ? settings.verificationCodeStorage : null,
        postSettings: handleSettingsSubmit('verificationCodeStorage'),
        verifySettings: handleSettingsVerification,
      },
      {
        title: 'Token Blacklist Storage',
        description: 'Setup a connection to the database all the blacklisted tokens are stored in.',
        settings: settings ? settings.tokenBlacklist : null,
        postSettings: handleSettingsSubmit('tokenBlacklist'),
        verifySettings: handleSettingsVerification,
      },
    ][index];
  };

  const storageSettingsProps = getStorageSettingsProps(getTabIndex(activeTab));
  return (
    <section className="iap-management-section">
      <header className="iap-management-section__header">
        <p className="iap-management-section__title">
          Storages
        </p>
      </header>
      <Tabs group={tabGroups.storages_group}>
        <Tab title={tabsTitles.applications} />
        <Tab title={tabsTitles.users} />
        <Tab title={tabsTitles.tokens} />
        <Tab title={tabsTitles.verifications_codes} />
        <Tab title={tabsTitles.blacklist} />
        <StorageSettings
          activeTabIndex={getTabIndex(activeTab)}
          connectionState={verificationStatus}
          onChange={() => setStatus(verificationStatuses.required)}
          progress={!!progress}
          {...storageSettingsProps}
        />
      </Tabs>
    </section>
  );
};

export default StoragesSection;
