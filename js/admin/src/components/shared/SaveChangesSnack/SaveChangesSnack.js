import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Snack } from '~/components/shared/Snack/Snack';
import { postServerSettings, updateServerSettings } from '~/modules/settings/actions';
import { getOriginalSettings } from '~/modules/settings/selectors';
import { showErrorNotificationSnack } from '~/modules/applications/notification-actions';
import { notificationStatuses } from '~/enums';
import useProgressBar from '~/hooks/useProgressBar';
import { errorSnackMessages } from '~/modules/applications/constants';

const actions = {
  save: 0,
  disgard: 1,
};
const config = {
  content: 'You have changed server settings, do you want to apply and restart the server?',
  buttons: [{ label: 'Apply and restart', data: actions.save }, { label: 'Discard all changes', data: actions.disgard }],
  status: notificationStatuses.changed,
};
export const SaveChangesSnack = () => {
  const dispatch = useDispatch();
  const { show } = useSelector(s => s.notifications.settingsSnack);
  const originalSettings = useSelector(getOriginalSettings);
  const { setProgress } = useProgressBar();

  const handler = async (action) => {
    if (action === actions.save) {
      setProgress(50);
      await dispatch(postServerSettings());
      setProgress(100);
    }
    if (action === actions.disgard) {
      dispatch(showErrorNotificationSnack(errorSnackMessages.settingsRejected));
      dispatch(updateServerSettings(originalSettings));
    }
  };

  if (!show) return null;

  return (
    <Snack {...config} callback={handler} />
  );
};
