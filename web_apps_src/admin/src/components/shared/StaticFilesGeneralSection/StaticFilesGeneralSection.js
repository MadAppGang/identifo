import update from '@madappgang/update-by-path';
import React from 'react';
import { useDispatch } from 'react-redux';
import { useVerification } from '~/hooks/useVerification';
import { hideSettingsDialog, handleSettingsDialog } from '~/modules/applications/actions';
import { storageTypes, verificationStatuses } from '~/enums';
import { settingsConfig, dialogActions } from '~/modules/applications/dialogsConfigs';
import { StaticFilesForm } from './Form';


const isVerificationType = (type) => {
  return type === storageTypes.local || type === storageTypes.s3;
};

const serializeVerificationPayload = (values) => {
  return {
    file_storage: {
      type: values.type,
      [values.type]: values[values.type],
    },
  };
};

export const StaticFilesGeneralSection = ({ title, subtitle, settings, onSubmit, onVerify }) => {
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

      if (res === dialogActions.submit) {
        onSubmit(payload);
      }
      if (res === dialogActions.verify) {
        onVerify(serializeVerificationPayload(formValues));
      }
    } else {
      onSubmit(payload);
    }
  };

  const handleSettingsVerification = (values) => {
    onVerify(serializeVerificationPayload(values));
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
