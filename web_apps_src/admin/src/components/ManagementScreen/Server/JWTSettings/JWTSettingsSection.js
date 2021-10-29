import update from '@madappgang/update-by-path';
import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AlgorithmDialog } from '~/components/ManagementScreen/Server/JWTSettings/AlgorithmDialog';
import useProgressBar from '~/hooks/useProgressBar';
import { handleSettingsDialog, hideSettingsDialog } from '~/modules/applications/actions';
import {
  dialogActions, privateKeyChangedConfig, showPrivateKeyConfig,
} from '~/modules/applications/dialogsConfigs';
import { generateKeys, getJWTKeys, setJWTKeys, uploadJWTKeys } from '~/modules/settings/actions';
import { selectJWTKeys } from '~/modules/settings/selectors';
import { JWTSettingsForm } from './Form';


export const JWTSettingsSection = () => {
  const dispatch = useDispatch();
  const keys = useSelector(selectJWTKeys);
  const [generateKeyDialogShown, setGenerateKeyDialogShown] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const { progress, setProgress } = useProgressBar();

  const getPrivateKey = async () => {
    const config = {
      ...showPrivateKeyConfig,
      onClose: () => dispatch(hideSettingsDialog()),
    };
    const res = await dispatch(handleSettingsDialog(config));
    if (!keys.private && res === dialogActions.submit) {
      setProgress(50);
      await dispatch(getJWTKeys(true));
      setProgress(100);
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
          setProgress(100);
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
