import React, { useEffect, useRef, useState } from 'react';
import Input from '~/components/shared/Input';
import FileUploadButton from './FileUploadButton';

const FileInput = (props) => {
  const fileInputRef = useRef(null);
  const [file, setFile] = useState(props.file || null);

  useEffect(() => {
    props.onFile(file);
  }, [file]);

  const handleReset = () => {
    setFile(null);
    fileInputRef.current.value = '';
  };

  return (
    <div className="file-input">
      <Input
        placeholder={props.placeholder}
        value={props.path}
        onValue={props.disablePathInput ? undefined : props.onPath}
        style={{ caretColor: props.disablePathInput ? 'transparent' : 'unset' }}
        disabled={props.disabled}
        errorMessage={props.errorMessage}
      />

      <FileUploadButton
        isUploaded={!!file}
        filename={file ? file.name : ''}
        onClick={() => fileInputRef.current.click()}
        onReset={handleReset}
      />

      {/* not visible */}
      <input
        type="file"
        ref={fileInputRef}
        style={{ display: 'none' }}
        onChange={e => setFile(e.target.files[0])}
      />
    </div>
  );
};

FileInput.defaultProps = {
  placeholder: 'Select File',
};

export default FileInput;
