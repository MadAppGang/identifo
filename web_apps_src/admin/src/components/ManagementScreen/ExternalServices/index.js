import update from '@madappgang/update-by-path';
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Tab, Tabs } from '~/components/shared/Tabs';
import useProgressBar from '~/hooks/useProgressBar';
import { getExternalServicesSettings } from '~/modules/settings/selectors';
import { tabGroups } from '../../../enums';
import { useQuery } from '../../../hooks/useQuery';
import { updateServerSettings } from '../../../modules/settings/actions';
import MailServiceSettings from './MailServiceSettings';
import SmsServiceSettings from './SmsServiceSettings';

const tabsTitles = {
  email_service: 'Email Service',
  sms_service: 'SMS Service',
};

const tabsMatcher = Object.values(tabsTitles).reduce((p, n) => {
  // eslint-disable-next-line no-param-reassign
  p[n] = n.toLowerCase().replaceAll(' ', '_');
  return p;
}, {});

const ExternalServicesSection = () => {
  const activeTab = useQuery().get(tabGroups.external_services_group);
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
          <Tabs group={tabGroups.external_services_group}>
            <Tab title={tabsTitles.email_service} />
            <Tab title={tabsTitles.sms_service} />
            <>
              {activeTab === tabsMatcher[tabsTitles.email_service] && (
                <MailServiceSettings
                  loading={!!progress}
                  settings={settings ? settings.emailService : null}
                  onSubmit={value => handleSubmit('emailService', value)}
                />
              )}

              {activeTab === tabsMatcher[tabsTitles.sms_service] && (
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
