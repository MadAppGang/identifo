import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Snack } from '~/components/shared/Snack/Snack';
import { postServerSettings, updateServerSettings } from '~/modules/settings/actions';
import { getOriginalSettings } from '~/modules/settings/selectors';
import { hideSettingsSnack } from '~/modules/applications/actions';


const actions = {
  save: 0,
  disgard: 1,
};
const config = {
  content: 'Your settings have been changed, do you want to save it?',
  buttons: [{ label: 'Save', data: actions.save }, { label: 'Don`t save', data: actions.disgard }],
};
export const SaveChangesSnack = () => {
  const dispatch = useDispatch();
  const { show } = useSelector(s => s.applicationDialogs.settingsSnack);
  const originalSettings = useSelector(getOriginalSettings);


  const handler = (action) => {
    if (action === actions.save) {
      dispatch(postServerSettings());
    }
    if (action === actions.disgard) {
      dispatch(hideSettingsSnack());
      dispatch(updateServerSettings(originalSettings));
    }
  };

  if (!show) return null;

  return (
    <Snack {...config} callback={handler} />
  );
};
