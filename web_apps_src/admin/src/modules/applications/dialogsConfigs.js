import { verificationStatuses } from '~/enums';

export const dialogActions = {
  cancel: 'cancel',
  submit: 'submit',
  verify: 'verify',
};

export const dialogTypes = {
  danger: 'danger',
};

export const settingsConfig = {
  [verificationStatuses.fail]: {
    content: 'The server was unable to connect with these settings',
    buttons: [
      { label: 'Save', data: dialogActions.submit },
      { label: 'Don`t Save', data: dialogActions.cancel, outline: true },
    ],
  },
  [verificationStatuses.required]: {
    content: 'It would be a good idea to test the connection before saving.',
    buttons: [
      { label: 'Verify', data: dialogActions.verify },
      { label: 'Save without verification', data: dialogActions.submit, outline: true },
    ],
  },
};

export const privateKeyChangedConfig = {
  content: 'If you save new key, all previous JWT will be removed, do you want to save?',
  buttons: [
    { label: 'Save', data: dialogActions.submit },
    { label: 'Cancel', data: dialogActions.cancel, outline: true },
  ],
};

export const showPrivateKeyConfig = {
  title: 'Confirmation alert',
  content: 'Are you sure you want to see the private key? This actions is not safe.',
  type: dialogTypes.danger,
  buttons: [
    { label: 'Show the key', data: dialogActions.submit, error: true },
    { label: 'Cancel', data: dialogActions.cancel, white: true, error: true },
  ],
};

export const disableServeAdmin = {
  title: 'Confirmation alert',
  content: 'If you disable serve admin you will no longer have access to the panel.',
  type: dialogTypes.danger,
  buttons: [
    { label: 'Disable', data: dialogActions.submit, error: true },
    { label: 'Cancel', data: dialogActions.cancel, white: true, error: true },
  ],
};
