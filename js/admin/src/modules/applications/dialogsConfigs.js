import { verificationStatuses } from '~/enums';

export const dialogActions = {
  cancel: 'cancel',
  submit: 'submit',
  verify: 'verify',
};

export const settingsConfig = {
  [verificationStatuses.fail]: {
    content: 'The server was unable to connect with these settings',
    buttons: [
      { label: 'Save', data: dialogActions.submit },
      { label: 'Don`t Save', data: dialogActions.cancel },
    ],
  },
  [verificationStatuses.required]: {
    content: 'It would be a good idea to test the connection before saving.',
    buttons: [
      { label: 'Verify', data: dialogActions.verify },
      { label: 'Save without verification', data: dialogActions.submit },
    ],
  },
};

export const privateKeyChangedConfig = {
  content: 'If you save new key, all previous JWT will be removed, do you want to save?',
  buttons: [
    { label: 'Save', data: dialogActions.submit },
    { label: 'Cancel', data: dialogActions.cancel },
  ],
};

export const showPrivateKeyConfig = {
  content: 'Are u sure you want to see the private key? This actions is not safe.',
  buttons: [
    { label: 'Show the key', data: dialogActions.submit },
    { label: 'Cancel', data: dialogActions.cancel, error: true },
  ],
};
