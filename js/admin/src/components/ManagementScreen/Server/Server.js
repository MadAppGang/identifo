import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Tab, Tabs } from '~/components/shared/Tabs';
import useNotifications from '~/hooks/useNotifications';
import useProgressBar from '~/hooks/useProgressBar';
import { uploadJWTKeys } from '~/modules/settings/actions';
import ConfigurationForm from './ServerConfigurationForm';
import GeneralTab from './ServerGeneralTab';
import JWTForm from './ServerJWTForm';
import { getKeyStorageSettings, getSettingsConfig } from '~/modules/settings/selectors';

const GeneralSection = () => {
  const [tabIndex, setTabIndex] = useState(0);
  const dispatch = useDispatch();
  const settings = useSelector(getKeyStorageSettings);
  const configSettings = useSelector(getSettingsConfig);
  const { notifySuccess } = useNotifications();

  const { progress, setProgress } = useProgressBar();

  const handleSubmit = async (nextSettings) => {
    setProgress(70);
    // TODO: Nikita k update settings

    const { privateKey, publicKey } = nextSettings;
    if (privateKey && publicKey) {
      await dispatch(uploadJWTKeys(publicKey, privateKey));
    }

    setProgress(100);

    notifySuccess({
      title: 'Updated',
      text: 'Server settings have been updated successfully',
    });
  };

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Server Settings
      </p>

      <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
        <Tab title="General" />
        <Tab title="Token Storage" />
        <Tab title="Token Settings" />
        <Tab title="Configuration Storage" />

        <>
          {tabIndex === 0 && <GeneralTab />}

          {tabIndex === 1 && (
            <JWTForm
              loading={!!progress}
              settings={settings}
              onSubmit={handleSubmit}
            />
          )}

          {tabIndex === 3 && (
            <ConfigurationForm
              loading={!!progress}
              settings={configSettings}
              onSubmit={handleSubmit}
            />
          )}
        </>
      </Tabs>
    </section>
  );
};

export default GeneralSection;
