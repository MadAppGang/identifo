import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useProgressBar from '~/hooks/useProgressBar';
import { handleSettingsDialog, hideSettingsDialog } from '~/modules/applications/actions';
import {
  dialogActions, showPrivateKeyConfig, privateKeyChangedConfig,
} from '~/modules/applications/dialogsConfigs';
import { generateKeys, getJWTKeys, uploadJWTKeys, setJWTKeys } from '~/modules/settings/actions';
import { selectJWTKeys } from '~/modules/settings/selectors';


import update from '@madappgang/update-by-path';
import { JWTSettingsForm } from './Form';
import { AlgorithmDialog } from '~/components/ManagementScreen/Server/JWTSettings/AlgorithmDialog';

export const JWTSettingsSection = () => {
  const dispatch = useDispatch();
  const keys = useSelector(selectJWTKeys);
  const [generateKeyDialogShown, setGenerateKeyDialogShown] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const { progress } = useProgressBar();

  const getPrivateKey = async () => {
    const config = {
      ...showPrivateKeyConfig,
      onClose: () => dispatch(hideSettingsDialog()),
    };
    const res = await dispatch(handleSettingsDialog(config));
    if (!keys.private && res === dialogActions.submit) {
      await dispatch(getJWTKeys(true));
    } else {
      dispatch(hideSettingsDialog());
    }
  };

  const onGenerateKeyHandler = async (act, alg) => {
    if (act === dialogActions.submit) {
      setIsLoading(true);
      try {
        await dispatch(generateKeys(alg));
        setGenerateKeyDialogShown(false);
        setIsLoading(false);
      } catch (error) {
        setIsLoading(false);
      }
    } else {
      setGenerateKeyDialogShown(false);
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
        case dialogActions.submit:
          await dispatch(uploadJWTKeys(nextSettings));
          break;
        case dialogActions.cancel:
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
    <>
      <JWTSettingsForm
        settings={keys}
        onShowPassword={getPrivateKey}
        onGenerateKey={() => setGenerateKeyDialogShown(true)}
        onSubmit={tokenSettingsSubmit}
        loading={!!progress}
      />
      {generateKeyDialogShown && (
        <AlgorithmDialog
          dialogHandler={onGenerateKeyHandler}
          onClose={() => setGenerateKeyDialogShown(false)}
          loading={isLoading}
        />
      )}
    </>
  );
};
