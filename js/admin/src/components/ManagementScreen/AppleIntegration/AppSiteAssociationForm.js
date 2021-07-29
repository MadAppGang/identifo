import React, { useEffect, useRef, useState } from 'react';
import { Controlled as CodeMirror } from 'react-codemirror2';
import Button from '~/components/shared/Button';
import FileIcon from '~/components/icons/FileIcon.svg';
import UploadIcon from '~/components/icons/UploadIcon.svg';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import useServices from '~/hooks/useServices';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';
import 'codemirror/mode/javascript/javascript';

let editor = null;

const AppSiteAssociationForm = () => {
  const services = useServices();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();
  const [content, setContent] = useState('{\n\t\n}');
  const fileInputRef = useRef(null);

  const fetchFileContents = async () => {
    setProgress(70);
    try {
      const result = await services.apple.fetchAppSiteAssociationFileContents();
      setContent(result);
    } finally {
      setProgress(100);
    }
  };

  useEffect(() => {
    fetchFileContents();
  }, []);

  const handleEditorClick = () => {
    if (editor) {
      editor.display.input.focus();
    }
  };

  const handleUpload = ({ target }) => {
    const file = target.files[0];
    const reader = new FileReader();

    reader.onload = () => setContent(reader.result);
    reader.readAsText(file);
  };

  const handleSubmit = async () => {
    setProgress(70);

    try {
      await services.apple.uploadAppSiteAssociationFileContents(content);

      notifySuccess({
        title: 'Uploaded',
        text: 'File has been uploaded.',
      });
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'File could not be uploaded.',
      });
    } finally {
      setProgress(100);
    }
  };

  return (
    <div className="app-site-association-form">
      <header className="template-editor-header">
        <p className="template-editor__filename">
          <FileIcon className="template-editor__file-icon" />
          apple-app-site-association
        </p>

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
      </header>

      {/* eslint-disable-next-line */}
      <div className="template-editor" onClick={handleEditorClick}>
        <CodeMirror
          editorDidMount={v => editor = v}
          value={content}
          options={{
            lineNumbers: true,
            theme: 'eclipse',
            mode: 'javascript',
            tabSize: 4,
          }}
          onBeforeChange={(_, data, value) => setContent(value)}
          className="template-editor-inner"
        />
        <div className="template-editor__numpad-area" />
      </div>

      <footer className="template-editor-footer">
        <Button
          Icon={progress ? LoadingIcon : SaveIcon}
          onClick={handleSubmit}
          disabled={!!progress}
        >
          Upload
        </Button>
      </footer>
    </div>
  );
};

export default AppSiteAssociationForm;
