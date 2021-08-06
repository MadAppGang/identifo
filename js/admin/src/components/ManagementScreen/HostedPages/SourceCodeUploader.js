import React from 'react';
import UploadIcon from '~/components/icons/UploadIcon.svg';

const SourceCodeUploader = ({ onCode }) => {
  const fileInputRef = React.useRef(null);

  const handleUpload = ({ target }) => {
    const file = target.files[0];
    const reader = new FileReader();

    reader.onload = () => onCode(reader.result);
    reader.readAsText(file);
  };

  return (
    <>
      <button
        className="template-editor__upload-code"
        onClick={() => fileInputRef.current.click()}
      >
        <UploadIcon className="template-editor__upload-icon" />
        Upload source code
      </button>

      <input
        type="file"
        onChange={handleUpload}
        ref={fileInputRef}
        style={{ display: 'none' }}
      />
    </>
  );
};

export default SourceCodeUploader;
