import { verificationStatuses } from '~/enums';


export const settingsActionsEnum = {
  save: 'save',
  verify: 'verify',
  close: 'close',
};

export const settingsConfig = {
  [verificationStatuses.fail]: {
    content: 'The server was unable to connect with these settings',
    buttons: [
      { label: 'Save', data: settingsActionsEnum.save },
      { label: 'Don`t Save', data: settingsActionsEnum.close },
    ],
  },
  [verificationStatuses.required]: {
    content: 'It would be a good idea to test the connection before saving.',
    buttons: [
      { label: 'Verify', data: settingsActionsEnum.verify },
      { label: 'Save without verification', data: settingsActionsEnum.save },
    ],
  },
};

export const generateKeyActions = {
  generate: 'generate',
  cancel: 'cancel',
};

export const generateKeyConfig = {
  content: 'If you generate new keys, all previous JWT will be removed, do you want to generate?',
  buttons: [
    { label: 'Generate and Save', data: generateKeyActions.generate },
    { label: 'Cancel', data: generateKeyActions.cancel },
  ],
};
export const privateKeyChangedActions = {
  save: 'save',
  cancel: 'cancel',
};

export const privateKeyChangedConfig = {
  content: 'If you save new key, all previous JWT will be removed, do you want to save?',
  buttons: [
    { label: 'Save', data: privateKeyChangedActions.save },
    { label: 'Cancel', data: privateKeyChangedActions.cancel },
  ],
};
