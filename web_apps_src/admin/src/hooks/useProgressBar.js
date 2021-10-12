import { useState, useContext } from 'react';
import { ProgressBarContext } from '~/components/shared/TopProgressBar';

const useProgressBar = (options = {}) => {
  const progressBar = useContext(ProgressBarContext);
  const [localProgress, setLocalProgress] = useState(0);

  if (options.local) {
    return {
      progress: localProgress,
      setProgress: setLocalProgress,
    };
  }

  return progressBar;
};

export default useProgressBar;
