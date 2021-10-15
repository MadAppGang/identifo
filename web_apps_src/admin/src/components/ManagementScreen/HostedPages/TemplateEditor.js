import React, { useEffect, useState } from 'react';
import { Controlled as CodeMirror } from 'react-codemirror2';
import FileIcon from '~/components/icons/FileIcon.svg';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import LanguageSelector from './LanguageSelector';
import ChangesIndicator from './ChangesIndicator';
import SourceCodeUploader from './SourceCodeUploader';
import UnsavedChangesPrompt from './UnsavedChangesPrompt';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/eclipse.css';
import 'codemirror/mode/htmlmixed/htmlmixed';
import 'codemirror/mode/css/css';
import 'codemirror/mode/javascript/javascript';

const [HTML, CSS, JS] = ['html', 'css', 'js'];

const mode = {
  [HTML]: 'htmlmixed',
  [CSS]: 'css',
  [JS]: 'javascript',
};

let editor = null;

const composeFilename = (name, ext) => {
  if (!name) {
    return '';
  }

  if (!ext) {
    return name;
  }

  return `${name}.${ext}`;
};

const TemplateEditor = (props) => {
  const { name, extension, source, progress, onChange } = props;
  const [code, setCode] = useState(source || '');
  const [hasChanged, setHasChanged] = useState(false);
  const [displayUnsavedChangesWarning, setDisplayUnsavedChangesWarning] = useState(false);

  useEffect(() => {
    if (!source) return;
    setCode(source);
  }, [source]);

  useEffect(() => {
    setHasChanged(code !== source);
  }, [code]);

  const handleEditorClick = () => {
    if (editor) {
      editor.display.input.focus();
    }
  };

  const handleLanguageChange = (value) => {
    if (hasChanged) {
      setDisplayUnsavedChangesWarning(true);
      return;
    }

    props.onExtensionChange(value);
  };


  const handleSubmit = () => {

  };

  return (
    <>
      <header className="template-editor-header">
        <div className="template-editor__filename">
          <FileIcon className="template-editor__file-icon" />
          {composeFilename(name, extension)}
          <ChangesIndicator visible={hasChanged} />
        </div>

        <SourceCodeUploader onCode={setCode} />
      </header>

      {/* eslint-disable-next-line */}
      <div className="template-editor" onClick={handleEditorClick}>
        <CodeMirror
          className="template-editor-inner"
          editorDidMount={v => editor = v}
          value={progress ? '' : code}
          options={{
            lineNumbers: true,
            theme: 'eclipse',
            mode: mode[extension],
          }}
          onBeforeChange={(_, data, value) => setCode(value)}
        />
        <div className="template-editor__numpad-area" />
        <LanguageSelector
          languages={[HTML, CSS, JS]}
          selected={extension}
          onChange={handleLanguageChange}
        />
        <UnsavedChangesPrompt
          visible={displayUnsavedChangesWarning}
          onDiscard={null}
          onCancel={() => setDisplayUnsavedChangesWarning(false)}
        />
      </div>

      <footer className="template-editor-footer">
        <Button
          Icon={progress ? LoadingIcon : SaveIcon}
          onClick={handleSubmit}
          disabled={progress}
        >
          Submit Template
        </Button>
      </footer>
    </>
  );
};

export default TemplateEditor;
