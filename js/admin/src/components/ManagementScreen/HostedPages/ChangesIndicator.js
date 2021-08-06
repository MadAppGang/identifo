import React from 'react';

const ChangesIndicator = ({ visible }) => {
  if (!visible) {
    return null;
  }

  return (
    <div className="file-changes-indicator-container">
      <div className="file-changes-indicator" />
      <div className="file-changes-tooltip">
        Unsaved changes
      </div>
    </div>
  );
};

ChangesIndicator.defaultProps = {
  visible: false,
};

export default ChangesIndicator;
