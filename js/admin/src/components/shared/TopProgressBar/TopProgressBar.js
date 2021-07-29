import React, { useState, useEffect, createContext } from 'react';

export const ProgressBarContext = createContext();

const TopProgressBar = ({ children }) => {
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    if (progress >= 100) {
      setProgress(0);
    }
  }, [progress]);

  const setProgressProxy = (value) => {
    if (value >= 100) {
      setProgress(99.9999);

      setTimeout(setProgress, 400, 100);
      return;
    }

    setProgress(value);
  };

  return (
    <ProgressBarContext.Provider value={{ progress, setProgress: setProgressProxy }}>
      {progress < 100 && (
        <div className="top-progress-bar" style={{ width: `${progress}vw` }} />
      )}

      {children}
    </ProgressBarContext.Provider>
  );
};

export default TopProgressBar;
