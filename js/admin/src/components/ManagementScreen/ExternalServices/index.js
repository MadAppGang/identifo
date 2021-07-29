import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import { useDispatch, useSelector } from 'react-redux';
import { Tabs, Tab } from '~/components/shared/Tabs';
import MailServiceSettings from './MailServiceSettings';
import SmsServiceSettings from './SmsServiceSettings';
import {
  fetchExternalServicesSettings, updateExternalServicesSettings,
} from '~/modules/settings/actions';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const ExternalServicesSection = () => {
  const [tabIndex, setTabIndex] = useState(0);
  const dispatch = useDispatch();
  const settings = useSelector(state => state.settings.externalServices);
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess } = useNotifications();

  useEffect(() => {
    if (!settings) {
      setProgress(70);
      dispatch(fetchExternalServicesSettings());
    }
  }, []);

  useEffect(() => {
    if (settings && progress) {
      setProgress(100);
    }
  }, [settings]);

  const handleSubmit = async (service, value) => {
    setProgress(70);

    const nextSettings = update(settings, {
      [service]: value,
    });

    await dispatch(updateExternalServicesSettings(nextSettings));

    notifySuccess({
      title: 'Updated',
      text: 'Settings have been updated successfully',
    });

    setProgress(100);
  };

  return (
    <section className="iap-management-section">
      <header>
        <p className="iap-management-section__title">
          External Services
        </p>

        <p className="iap-management-section__description">
          Configure external Email ans SMS service integrations
        </p>
      </header>

      <main className="iap-settings-section">
        <div className="iap-management-section__tabs">
          <Tabs activeTabIndex={tabIndex} onChange={setTabIndex}>
            <Tab title="Email Service" />
            <Tab title="SMS Service" />

            <>
              {tabIndex === 0 && (
                <MailServiceSettings
                  loading={!!progress}
                  settings={settings ? settings.emailService : null}
                  onSubmit={value => handleSubmit('emailService', value)}
                />
              )}

              {tabIndex === 1 && (
                <SmsServiceSettings
                  loading={!!progress}
                  settings={settings ? settings.smsService : null}
                  onSubmit={value => handleSubmit('smsService', value)}
                />
              )}
            </>
          </Tabs>
        </div>
      </main>
    </section>
  );
};

export default ExternalServicesSection;
