import React from 'react';
import update from '@madappgang/update-by-path';
import Notification from './Notification';
import usePrevious from '~/hooks/usePrevious';

const NOTIFICATION_LIFETIME = 3000;

export const NotificationContext = React.createContext();

const NotificationContainer = ({ children }) => {
  const [notifications, setNotifications] = React.useState([]);
  const previousNotifications = usePrevious(notifications);

  const isNewNotification = (notification) => {
    return !previousNotifications.map(n => n.id).includes(notification.id);
  };

  const excludeExistingNotifications = (list) => {
    return list.filter(isNewNotification);
  };

  const removeNotification = (notification) => {
    setNotifications(notifications.filter(n => n.id !== notification.id));
  };

  const removeAfterDelay = (notification) => {
    setTimeout(removeNotification, NOTIFICATION_LIFETIME, notification);
  };

  React.useEffect(() => {
    excludeExistingNotifications(notifications).forEach(removeAfterDelay);
  }, [notifications]);

  const createNotificationOfType = (type, notification) => {
    return update(notification, { type, id: Date.now() });
  };

  const context = {
    notifySuccess(notification) {
      setNotifications([
        ...notifications,
        createNotificationOfType('success', notification),
      ]);
    },
    notifyFailure(notification) {
      setNotifications([
        ...notifications,
        createNotificationOfType('failure', notification),
      ]);
    },
  };

  return (
    <NotificationContext.Provider value={context}>
      <div className="iap-notification-container">
        {notifications.map(notification => (
          <Notification
            key={notification.id}
            {...notification}
            onClick={() => removeNotification(notification)}
          />
        ))}
      </div>

      {children}
    </NotificationContext.Provider>
  );
};

export default NotificationContainer;
