import { useContext } from 'react';
import { NotificationContext } from '~/components/shared/Notifications';

const useNotifications = () => {
  return useContext(NotificationContext);
};

export default useNotifications;
