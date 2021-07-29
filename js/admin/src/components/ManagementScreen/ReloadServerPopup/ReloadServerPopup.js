import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { restartServer } from '~/modules/settings/actions';
import LoadingIcon from '~/components/icons/LoadingIcon';
import useProgressBar from '~/hooks/useProgressBar';

const ReloadServerPopup = () => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar({ local: true });
  const settingsChanged = useSelector(s => s.settings.changed);

  if (!settingsChanged) {
    return null;
  }

  const handleRestartClick = () => {
    setProgress(70);
    dispatch(restartServer());
  };

  return (
    <div className="popup">
      {!!progress && (
        <>
          <p className="popup-text">
            Restarting the server
          </p>

          <LoadingIcon className="popup-loading" />
        </>
      )}

      {!progress && (
        <>
          <p className="popup-text">
            In order for the changes to take effect, please restart the server.
          </p>

          <button className="popup-btn" onClick={handleRestartClick}>
            Restart
          </button>
        </>
      )}
    </div>
  );
};

export default ReloadServerPopup;
