import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { verificationStatuses } from '~/enums';
import useProgressBar from '~/hooks/useProgressBar';
import { useVerification } from '~/hooks/useVerification';
import { handleSettingsDialog, hideSettingsDialog } from '~/modules/applications/actions';
import { dialogActions, settingsConfig } from '~/modules/applications/dialogsConfigs';
import { updateServerSettings } from '~/modules/settings/actions';
import { getKeyStorageSettings } from '~/modules/settings/selectors';
import JWTForm from './Form';

export const JWTStorageSection = () => {
  const dispatch = useDispatch();
  const [verificationStatus, verify, setStatus] = useVerification();
  const settings = useSelector(getKeyStorageSettings);

  const { progress } = useProgressBar();

  const handleSettingsVerification = async (nodeSettings) => {
    const payload = { type: 'key_storage', keyStorage: nodeSettings };
    await dispatch(verify(payload));
  };

  const tokenStorageSubmit = async (nextSettings) => {
    if (verificationStatus !== verificationStatuses.success) {
      const config = {
        ...settingsConfig[verificationStatus],
        onClose: () => dispatch(hideSettingsDialog()),
      };
      const action = await dispatch(handleSettingsDialog(config));
      switch (action) {
        case dialogActions.submit:
          dispatch(updateServerSettings({ keyStorage: nextSettings }));
          break;
        case dialogActions.verify:
          await handleSettingsVerification(nextSettings);
          break;
        default:
          dispatch(hideSettingsDialog());
          break;
      }
    } else {
      await dispatch(updateServerSettings({ keyStorage: nextSettings }));
    }
  };
  if (!settings.type) return null;
  return (
    <JWTForm
      loading={!!progress}
      settings={settings}
      onSubmit={tokenStorageSubmit}
      onChange={() => setStatus(verificationStatuses.required)}
      verificationStatus={verificationStatus}
      handleVerify={handleSettingsVerification}
    />
  );
};
