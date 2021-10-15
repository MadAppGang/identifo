import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import useProgressBar from '~/hooks/useProgressBar';
import { updateServerSettings } from '~/modules/settings/actions';
import { getGeneralServerSettings } from '~/modules/settings/selectors';
import ServerGeneralForm from './ServerGeneralForm';

const ServerGeneralTab = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getGeneralServerSettings);
  const { progress } = useProgressBar();

  const handleSubmit = async (nextSettings) => {
    const payload = { general: { ...nextSettings } };
    dispatch(updateServerSettings(payload));
  };

  return (
    <ServerGeneralForm
      loading={!!progress}
      settings={settings}
      onSubmit={handleSubmit}
    />
  );
};

export default ServerGeneralTab;
