import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useProgressBar from '~/hooks/useProgressBar';
import { handleSettingsDialog, hideSettingsDialog } from '~/modules/applications/actions';
import {
  generateKeyActions, generateKeyConfig,
  privateKeyChangedActions, privateKeyChangedConfig,
} from '~/modules/applications/dialogsConfigs';
import { generateKeys, getJWTKeys, uploadJWTKeys, setJWTKeys } from '~/modules/settings/actions';
import { selectJWTKeys } from '~/modules/settings/selectors';


import update from '@madappgang/update-by-path';
import { JWTSettingsForm } from './Form';

export const JWTSettingsSection = () => {
  const dispatch = useDispatch();
  const keys = useSelector(selectJWTKeys);

  const { progress } = useProgressBar();

  const getPrivateKey = async () => {
    if (!keys.private) {
      await dispatch(getJWTKeys(true));
    }
  };

  const onGenerateKey = async (alg) => {
    const config = {
      ...generateKeyConfig,
      onClose: () => dispatch(hideSettingsDialog()),
    };
    const res = await dispatch(handleSettingsDialog(config));
    if (res === generateKeyActions.generate) {
      await dispatch(generateKeys(alg));
    } else {
      dispatch(hideSettingsDialog());
    }
  };

  const tokenSettingsSubmit = async (nextSettings) => {
    if (nextSettings.private !== keys.private || nextSettings.alg !== keys.alg) {
      const config = {
        ...privateKeyChangedConfig,
        onClose: () => dispatch(hideSettingsDialog()),
      };
      const res = await dispatch(handleSettingsDialog(config));
      switch (res) {
        case privateKeyChangedActions.save:
          await dispatch(uploadJWTKeys(nextSettings));
          break;
        case privateKeyChangedActions.cancel:
          dispatch(hideSettingsDialog());
          break;
        default:
          break;
      }
    }
  };

  useEffect(() => {
    return () => {
      dispatch(setJWTKeys(update(keys, { private: '' })));
    };
  }, []);

  return (
    <JWTSettingsForm
      settings={keys}
      onShowPassword={getPrivateKey}
      onGenerateKey={onGenerateKey}
      onSubmit={tokenSettingsSubmit}
      loading={!!progress}
    />

  );
};
