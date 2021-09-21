import React, { useMemo, useState } from 'react';
import { DialogPopup } from '~/components/shared/DialogPopup/DialogPopup';
import { Option, Select } from '~/components/shared/Select';
import { dialogActions } from '~/modules/applications/dialogsConfigs';


export const AlgorithmDialog = ({ onClose, dialogHandler, loading }) => {
  const [algorithm, setAlgorithm] = useState('');

  const algorithmConfig = useMemo(() => ({
    title: 'Generate new keys',
    content: 'If you generate new keys, all previous JWT will be removed, do you want to generate?',
    buttons: [
      { label: 'Generate and Save', data: dialogActions.submit, disabled: !algorithm || loading },
      { label: 'Cancel', data: dialogActions.cancel, outline: true },
    ],
  }), [algorithm, loading]);

  const handler = (act) => {
    dialogHandler(act, algorithm);
  };

  return (
    <DialogPopup {...algorithmConfig} callback={handler} onClose={onClose}>
      <Select
        value={algorithm}
        onChange={setAlgorithm}
        placeholder="Select algorithm"
      >
        <Option value="RS256" title="RS256" />
        <Option value="ES256" title="ES256" />
      </Select>
    </DialogPopup>
  );
};
