import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useProgressBar from '~/hooks/useProgressBar';
import { handleSettingsDialog, hideSettingsDialog } from '~/modules/applications/actions';
import { dialogActions, disableServeAdmin } from '~/modules/applications/dialogsConfigs';
import { updateServerSettings } from '~/modules/settings/actions';
import { getAdminPanelSettings, getGeneralServerSettings } from '~/modules/settings/selectors';
import ServerGeneralForm from './ServerGeneralForm';

const ServerGeneralTab = () => {
  const settings = useSelector(getGeneralServerSettings);
  const adminPanelSettigns = useSelector(getAdminPanelSettings);
  const [serveAdmin, setServeAdmin] = useState(false);
  const dispatch = useDispatch();

  const { progress } = useProgressBar();

  const handleSubmit = async (nextSettings) => {
    const payload = { general: { ...nextSettings }, adminPanel: { enabled: serveAdmin } };
    dispatch(updateServerSettings(payload));
  };

  const onServeAdminChange = async (status) => {
    const config = {
      ...disableServeAdmin,
      onClose: () => dispatch(hideSettingsDialog()),
    };
    if (!status) {
      const res = await dispatch(handleSettingsDialog(config));
      if (res === dialogActions.submit) {
        setServeAdmin(status);
      }
    } else {
      setServeAdmin(status);
    }
  };

  useEffect(() => {
    if (adminPanelSettigns) {
      setServeAdmin(adminPanelSettigns.enabled);
    }
  }, [adminPanelSettigns]);

  return (
    <ServerGeneralForm
      loading={!!progress}
      settings={settings}
      serveAdmin={serveAdmin}
      onServeAdminChange={onServeAdminChange}
      onSubmit={handleSubmit}
    />
  );
};

export default ServerGeneralTab;
