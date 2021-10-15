import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import SessionStorageForm from './SessionStorageForm';
import useProgressBar from '~/hooks/useProgressBar';
import { updateServerSettings } from '../../../modules/settings/actions';
import { getSessionStorageSettings } from '~/modules/settings/selectors';

const SessionStorageSettings = ({ error }) => {
  const dispatch = useDispatch();
  const { progress } = useProgressBar();
  const settings = useSelector(getSessionStorageSettings);

  const handleSubmit = async (data) => {
    try {
      await dispatch(updateServerSettings({ sessionStorage: data }));
    } catch (e) {
      console.log(e);
    }
  };

  return (
    <SessionStorageForm
      error={error}
      loading={!!progress}
      settings={settings}
      onSubmit={handleSubmit}
    />
  );
};

export default SessionStorageSettings;
