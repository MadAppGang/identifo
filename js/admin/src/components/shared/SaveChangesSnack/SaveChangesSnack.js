import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Snack } from '~/components/shared/Snack/Snack';
import { postServerSettings, updateServerSettings } from '~/modules/settings/actions';
import { getOriginalSettings } from '~/modules/settings/selectors';
import { showNotificationSnack } from '~/modules/applications/actions';
import { notificationStates } from '~/modules/applications/notificationsStates';


const actions = {
  save: 0,
  disgard: 1,
};
const config = {
  content: 'You have changed server settings, do you want to apply and restart the server?',
  buttons: [{ label: 'Apply and restart', data: actions.save }, { label: 'Discard all changes', data: actions.disgard }],
  status: 'changed',
};
export const SaveChangesSnack = () => {
  const dispatch = useDispatch();
  const { show } = useSelector(s => s.notifications.settingsSnack);
  const originalSettings = useSelector(getOriginalSettings);


  const handler = async (action) => {
    if (action === actions.save) {
      dispatch(postServerSettings());
    }
    if (action === actions.disgard) {
      dispatch(showNotificationSnack(notificationStates.rejected.status));
      dispatch(updateServerSettings(originalSettings));
    }
  };

  if (!show) return null;

  return (
    <Snack {...config} callback={handler} />
  );
};
