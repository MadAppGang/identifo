import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import SessionStorageForm from './SessionStorageForm';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';
import { updateServerSettings } from '../../../modules/settings/actions';
import { getSessionStorageSettings } from '~/modules/settings/selectors';

const SessionStorageSettings = ({ error }) => {
  const dispatch = useDispatch();
  const { notifySuccess, notifyFailure } = useNotifications();
  const { progress } = useProgressBar();
  const settings = useSelector(getSessionStorageSettings);

  const handleSubmit = async (data) => {
    try {
      await dispatch(updateServerSettings({ sessionStorage: data }));
      notifySuccess({
        title: 'Updated',
        text: 'Settings have been updated successfully',
      });
    } catch (e) {
      notifyFailure({
        title: 'Erro',
        text: 'some error',
      });
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
