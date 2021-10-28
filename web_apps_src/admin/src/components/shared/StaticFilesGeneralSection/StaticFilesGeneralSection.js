import update from '@madappgang/update-by-path';
import React from 'react';
import { useDispatch } from 'react-redux';
import { useVerification } from '~/hooks/useVerification';
import { hideSettingsDialog, handleSettingsDialog } from '~/modules/applications/actions';
import { storageTypes, verificationStatuses } from '~/enums';
import { settingsConfig, dialogActions } from '~/modules/applications/dialogsConfigs';
import { StaticFilesForm } from './Form';

// TODO: Add verification after backed will be done
const isVerificationType = (type) => {
  return type === storageTypes.local || type === storageTypes.s3;
};

export const StaticFilesGeneralSection = ({ title, subtitle, settings, onSubmit }) => {
  const dispatch = useDispatch();
  const [verificationStatus, , setStatus] = useVerification();


  const onChange = () => {
    setStatus(verificationStatuses.required);
  };

  const submitHandler = async (formValues) => {
    const payload = update(settings, formValues);

    const config = {
      ...settingsConfig[verificationStatus],
      onClose: () => dispatch(hideSettingsDialog()),
    };

    if (isVerificationType(formValues.type)
        && verificationStatus !== verificationStatuses.success) {
      const res = await dispatch(handleSettingsDialog(config));
      switch (res) {
        case dialogActions.submit: {
          onSubmit(payload);
          break;
        }
        case dialogActions.verify: {
          setStatus(verificationStatuses.success);
          break;
        }
        default:
          break;
      }
    } else {
      dispatch(onSubmit(payload));
    }
  };

  const handleSettingsVerification = () => {
    setStatus(verificationStatuses.success);
  };

  if (!settings) return null;

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        {title}
      </p>
      <p className="iap-management-section__description">
        {subtitle}
      </p>
      <StaticFilesForm
        settings={settings}
        verificationStatus={verificationStatus}
        onSubmit={submitHandler}
        onVerify={handleSettingsVerification}
        onChange={onChange}
      />
    </section>
  );
};
