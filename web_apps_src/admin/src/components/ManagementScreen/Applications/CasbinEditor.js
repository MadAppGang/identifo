import React from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import 'codemirror/lib/codemirror.css';

const editorOptions = {
  lineNumbers: true,
  preserveScrollPosition: true,
  viewportMargin: Infinity,
};

const CasbinEditor = (props) => {
  return (
    <CodeMirror
      className="casbin-code-editor"
      value={props.value}
      onBeforeChange={(editor, _, value) => props.onChange(value)}
      options={editorOptions}
    />
  );
};

export default CasbinEditor;
