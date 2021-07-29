import React from 'react';

const UnsavedChangesPrompt = ({ visible, onDiscard, onCancel }) => {
  if (!visible) {
    return null;
  }

  return (
    <div className="editor-changes-warning">
      <div>You have unsaved changes. Discard them?</div>
      <div className="editor-changes-warning__btns">
        <button
          className="editor-changes-warning__btn"
          onClick={onCancel}
        >
          Cancel
        </button>
        <button
          className="editor-changes-warning__btn editor-changes-warning__btn--dimmed"
          onClick={onDiscard}
        >
          Discard
        </button>
      </div>
    </div>
  );
};

export default UnsavedChangesPrompt;
