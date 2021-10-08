import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import { useDispatch, useSelector } from 'react-redux';
import { Tabs, Tab } from '~/components/shared/Tabs';
import MailServiceSettings from './MailServiceSettings';
import SmsServiceSettings from './SmsServiceSettings';
import useProgressBar from '~/hooks/useProgressBar';
import { getExternalServicesSettings } from '~/modules/settings/selectors';
import { updateServerSettings } from '../../../modules/settings/actions';

const ExternalServicesSection = () => {
  const [tabIndex, setTabIndex] = useState(0);
  const dispatch = useDispatch();
  const settings = useSelector(getExternalServicesSettings);
  const { progress, setProgress } = useProgressBar();

  useEffect(() => {
    if (settings && progress) {
      setProgress(100);
    }
  }, [settings]);

  const handleSubmit = async (service, value) => {
    const nextSettings = { externalServices: update(settings, {
      [service]: value,
    }) };

    dispatch(updateServerSettings(nextSettings));
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
