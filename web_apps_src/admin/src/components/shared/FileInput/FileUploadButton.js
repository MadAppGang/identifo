import React from 'react';
import classnames from 'classnames';
import UploadIcon from '~/components/icons/UploadIcon.svg';
import CheckIcon from '~/components/icons/CheckIcon.svg';

const FileUploadButton = ({ filename, isUploaded, onClick, onReset }) => {
  const [showPrompt, setShowPrompt] = React.useState(false);

  const classname = classnames('file-input__upload', {
    'file-input__upload--success': isUploaded,
  });

  const Icon = isUploaded ? CheckIcon : UploadIcon;

  const handleResetClick = () => {
    onReset();
    setShowPrompt(false);
  };

  return (
    <div className="file-input__upload-wrapper">
      <button
        type="button"
        className={classname}
        onClick={isUploaded ? () => setShowPrompt(!showPrompt) : onClick}
      >
        <Icon className="file-input__upload-icon" />
      </button>

      {showPrompt && (
        <div className="upload-prompt">
          <p className="upload-prompt__filename">
            {filename}
          </p>
          <button
            type="button"
            className="upload-prompt__another"
            onClick={handleResetClick}
          >
            Reset
          </button>
        </div>
      )}
    </div>
  );
};

export default FileUploadButton;
